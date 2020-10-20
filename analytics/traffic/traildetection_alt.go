package traffic

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sync"

	utility "gitlab.com/suvam0451/trafficdetection/utility"
)

// FrameTaggedArchive : Holds history of all past frames
type FrameTaggedArchive struct {
	FrameID     int                   `json:"frame_id"`
	FrameRecord []PreviousFrameObject `json:"frames"`
}

// CustomFrameData : Main container for our combined data segments
type CustomFrameData []CompactFrame

// CompactFrame : More compacted frame datas
type CompactFrame struct {
	FrameID int             `json:"frame_id"`
	Objects []CompactObject `json:"objects"`
}

// TaggedObject : Represents a tagged vehicle and its trajectory
type TaggedObject struct {
	ObjectID int             `json:"id"`
	ClassID  int             `json:"class_id"`
	Data     []CompactCoords `json:"data"`
}

// CompactObject : More compact single object data
type CompactObject struct {
	ClassID    int     `json:"class_id"`
	CenterX    float64 `json:"center_x"`
	CenterY    float64 `json:"center_y"`
	confidence float64 `json:"confidence"`
}

// PreviousFrameObject : Struct to hold
type PreviousFrameObject struct {
	ObjectID   int     `json:"id"`
	ClassID    int     `json:"class_id"`
	CenterX    float64 `json:"center_x"`
	CenterY    float64 `json:"center_y"`
	confidence float64 `json:"confidence"`
	TagCounter int     `json:"tagcounter"`
	tagged     bool
}

// CompactCoords : More compact ccord data
type CompactCoords struct {
	CenterX    float64 `json:"center_x"`
	CenterY    float64 `json:"center_y"`
	Confidence float64 `json:"confidence"`
}

// DetectTrailCustom : I am modifying the codebase for more compact data format
func DetectTrailCustom(inputpath string, params ModelParameters) {

	var wg sync.WaitGroup

	// regex groups
	reA := regexp.MustCompile(`^A_2_[0-9]_02\.json$`)
	reB := regexp.MustCompile(`^B_2_[0-9]_02\.json$`)
	reD := regexp.MustCompile(`^D_2_[0-9]_02\.json$`)
	reF := regexp.MustCompile(`^F_2_[0-9]_02\.json$`)
	reG := regexp.MustCompile(`^G_2_[0-9]_02\.json$`)
	var regexArray = []*regexp.Regexp{reA, reB, reD, reF, reG}
	var outputPath = []string{"A", "B", "D", "F", "G"}

	for i, regex := range regexArray {
		// ounter := 1
		var inputfilespath []string

		// Error already handled above
		filepath.Walk(inputpath,
			func(path string, info os.FileInfo, err error) error {
				if info.IsDir() == false && regex.MatchString(info.Name()) {
					inputfilespath = append(inputfilespath, path)
				}
				return nil
			})

		// Unify the data files and genearte a single file
		var SourceObject CustomFrameData
		for _, file := range inputfilespath {
			var tmpData CustomFrameData
			if openfile, err := os.Open(file); err == nil {
				byteValue, _ := ioutil.ReadAll(openfile)
				if err := json.Unmarshal(byteValue, &tmpData); err == nil {
					SourceObject = append(SourceObject, tmpData...)
				}
			}
		}

		outpath := "out_traildetect_alt"
		utility.MakeDirectory(outpath)

		// Every thread handles one group of files
		wg.Add(1)
		go func(data CustomFrameData, param ModelParameters, outpath, outfile string) {
			runAnalysis(SourceObject, params, outpath, outfile)
			wg.Done()
		}(SourceObject, params, outpath, outputPath[i])
	}

	wg.Wait()
}

// FindIntInSlice takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func FindIntInSlice(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func itemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Array {
		panic("Invalid data-type")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}
	return false
}

/* For every alternate frame, compares to the previous frame and classifies objects that may be idle/moving with configurable accuracy */
func runAnalysis(source CustomFrameData, params ModelParameters, outpath, outfile string) {
	var previousFrameData []PreviousFrameObject
	var theArchive []FrameTaggedArchive
	var perVehicleTrack trackArchive
	var vehicleIDIndex int = 0
	var acceptedIDs = []int{2, 3, 5, 7} // car, motorbike, bus, truck

	for i, frame := range source {
		for _, currentobj := range frame.Objects {
			tagged := false

			// Skip for mismatch with following classes : "car", "motorbike", "bus", "truck"
			if !FindIntInSlice(acceptedIDs, currentobj.ClassID) {
				continue
			}

			for idx, prevobj := range previousFrameData {
				if prevobj.ClassID == currentobj.ClassID && !prevobj.tagged && !previousFrameData[idx].tagged {
					// TAG_SUCCESS case : a close enough co-ordinate was detected
					// for a previously existing entry
					previousFrameData[idx].CenterX = currentobj.CenterX
					previousFrameData[idx].CenterY = currentobj.CenterY
					previousFrameData[idx].confidence = currentobj.confidence
					previousFrameData[idx].tagged = true

					// TAG_SUCCESS case : increment the co-ordinates to the list
					_ID := previousFrameData[idx].ObjectID
					perVehicleTrack[_ID].TrackPoints = append(perVehicleTrack[_ID].TrackPoints,
						CompactCoords{
							CenterX:    currentobj.CenterX,
							CenterY:    currentobj.CenterY,
							Confidence: currentobj.confidence,
						})
					// TAG_SUCCESS case : increment the #frames for which object was tracked
					perVehicleTrack[_ID].FrameCount++

					tagged = true
					break
				}
			}

			// Handle if object was untagged (new object detected)
			if !tagged {
				// TAG_FAILURE case : Add entry for new vehicleID in list of vehicle tracks
				perVehicleTrack = append(perVehicleTrack, VehicleTracks{
					VehicleID:  vehicleIDIndex,
					FrameCount: 1,
					ClassID:    currentobj.ClassID,
				})

				// TAG_FAILURE case : the vehicleID must exist in the perVehicleTrack arrays
				perVehicleTrack[vehicleIDIndex].TrackPoints = append(perVehicleTrack[vehicleIDIndex].TrackPoints, CompactCoords{
					CenterX:    currentobj.CenterX,
					CenterY:    currentobj.CenterY,
					Confidence: currentobj.confidence,
				})

				tmpStruct := PreviousFrameObject{
					ObjectID:   vehicleIDIndex,
					ClassID:    currentobj.ClassID,
					CenterX:    currentobj.CenterX,
					CenterY:    currentobj.CenterY,
					confidence: currentobj.confidence,
					TagCounter: 0,
					tagged:     true,
				}

				previousFrameData = append(previousFrameData, tmpStruct)

				// In the end, Increment index for next vehicle ID
				vehicleIDIndex++
			}
		}
		// Increment if tagged and reset tag status (SUCCESS_CASE handled already)
		for idx, prevobj := range previousFrameData {
			if prevobj.tagged {
				previousFrameData[idx].TagCounter += params.Upvote
			} else {
				previousFrameData[idx].TagCounter += params.Downvote
			}
			previousFrameData[idx].tagged = false
		}

		// Eliminate any entry which was not tagged recently
		previousFrameData, _ = Filter02(previousFrameData, params.EliminateThreshold)

		// Add frame record to archive
		theArchive = append(theArchive, FrameTaggedArchive{
			FrameID:     i,
			FrameRecord: previousFrameData,
		})
	}

	// Ensure all output paths exist...
	if _, err := os.Stat("./out_traildetection_alt"); os.IsNotExist(err) {
		os.Mkdir("./out_traildetection_alt", os.ModeDir)
	}

	// Write data to file
	if jsonString, err := json.MarshalIndent(theArchive, "", " "); err == nil {
		ioutil.WriteFile("./out_traildetection_alt/veh_"+outfile+".json", jsonString, 0644)
	}

	// Test (pruned data - at least 10 frames) --> Noise
	accepted, _ := PruneFalsePositives(perVehicleTrack, 5)
	if jsonString, err := json.MarshalIndent(accepted, "", " "); err == nil {
		ioutil.WriteFile("./out_traildetection_alt/veh_"+outfile+"_c.json", jsonString, 0644)
	}
}
