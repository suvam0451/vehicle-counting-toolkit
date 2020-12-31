/**
Outputs the following :-

1. filename.json --
2. filename_veh.json -- Tracked frames for each vehicle tag. Entries with low score are eliminated.
*/

package traffic

import (
	"encoding/json"
	"fmt"
	"gitlab.com/suvam0451/trafficdetection/utility"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type mainArchive []archiveRecord

type trackArchive []VehicleTracks

// VehicleTracks :  List of positions w/confidence through which vehicle has passed
type VehicleTracks struct {
	VehicleID   int             `json:"vehicle_id"`  // ID given to the vehicle
	FrameCount  int             `json:"frame_count"` // Number of frames for which this object was detected
	ClassID     int             `json:"class_id"`    // ClassID for this vehicle type
	TrackPoints []CompactCoords `json:"objects"`     // List of co-ordinates
}

// Holds record for all the previous frames
type archiveRecord struct {
	FrameID     int         `json:"frame_id"`
	FrameRecord frameRecord `json:"objects"`
}

type frameRecord []objectHistory

type objectHistory struct {
	VehicleID           int                 `json:"id"` // ID assigned to every unique vehicle path
	ClassID             int                 `json:"class_id"`
	Name                Name                `json:"name"`
	RelativeCoordinates relativeCoordinates `json:"relative_coordinates"`
	TagCounter          int                 `json:"tagcounter"`
	tagged              bool
}

type relativeCoordinates struct {
	CenterX float64 `json:"center_x"`
	CenterY float64 `json:"center_y"`
}

// Name : List of tag names in darknet
type Name string

// ModelParameters : Parameters for our model
type ModelParameters struct {
	Rewards            int
	Penalty            int
	XThreshold         float64
	YThreshold         float64
	EliminateThreshold int
}

// Filter : Filters struct array for minimum threshold value
func Filter02(vs []PreviousFrameObject, threshold int) ([]PreviousFrameObject, []PreviousFrameObject) {
	accepted := make([]PreviousFrameObject, 0)
	rejected := make([]PreviousFrameObject, 0)
	for _, v := range vs {
		if v.TagCounter > threshold {
			accepted = append(accepted, v)
		} else {
			rejected = append(rejected, v)
		}
	}
	return accepted, rejected
}

// Filter : Filters struct array for minimum threshold value
func Filter(vs []objectHistory, threshold int) ([]objectHistory, []objectHistory) {
	accepted := make([]objectHistory, 0)
	rejected := make([]objectHistory, 0)
	for _, v := range vs {
		if v.TagCounter > threshold {
			accepted = append(accepted, v)
		} else {
			rejected = append(rejected, v)
		}
	}
	return accepted, rejected
}

// PruneFalsePositives : Any vehicle trail with < minThreshold number of data will be pruned
func PruneFalsePositives(archive []VehicleTracks, minThreshold int) (accepted []VehicleTracks, rejected []VehicleTracks) {
	for _, v := range archive {
		if v.FrameCount > minThreshold {
			accepted = append(accepted, v)
		} else {
			rejected = append(rejected, v)
		}
	}
	return
}

func CreateMissingDirectories() {
	if r := recover(); r != nil {
		utility.EnsureFile("./config.json")
		fmt.Println("Restored required missing file: ", r)
	}
}

// DetectTrail detect trails for all files in given path
func DetectTrail(inputPath string, config TrailDetectAltConfig) {
	defer CreateMissingDirectories()
	var inputFilesPath []string
	var inputFilesName []string
	var wg sync.WaitGroup

	if eval, _ := utility.PathExists(inputPath); eval == false {
		panic("input filepath missing !")
	}

	//// Error already handled above
	filepath.Walk(inputPath,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() == false {
				inputFilesPath = append(inputFilesPath, path)
				inputFilesName = append(inputFilesName, info.Name())
			}
			return nil
		})

	var ParsedStruct TrailData_Source
	for i, file := range inputFilesPath {
		if openFile, err := os.Open(file); err == nil {
			byteValue, _ := ioutil.ReadAll(openFile)
			if err := json.Unmarshal(byteValue, &ParsedStruct); err == nil {
				wg.Add(1)

				// Send structs to separate threads
				go func(filename string) {
					detectIndividualTrail(ParsedStruct, config, filename)
					wg.Done()
				}(inputFilesName[i])
			}
		}
	}
	wg.Wait()
}

func detectIndividualTrail(data TrailData_Source, params TrailDetectAltConfig, filepath string) {
	outputDir := "./out_traildetect"
	// utility.MakeDirectory(outputDir)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModeDir)
	}

	previousFrameData := frameRecord{}
	theArchive := mainArchive{}
	perVehicleTrack := trackArchive{}
	vehicleIDIndex := 0
	/*
		LOGIC
		----------
		For each frame, loop over the objects, compare the elements to the previously stored elements.
		If classes match and the object has not been tagged, check the displacement and see if it's under the threshold.
		Add the item to list.

		If the element is still untagged, then check the classID and insert it as a new key

		Depending on if the element was tagged/untagged, add/remove scores
		Filter out elements with very low scores
	*/
	for i, frame := range data {
		for _, currentObj := range frame.Objects {
			tagged := false // will be set to true if object gets assigned to one of the previous frame objects

			for idx, prevObj := range previousFrameData {
				// ID match with untagged object
				if prevObj.ClassID == currentObj.ClassID && !prevObj.tagged {
					// Distance calculations
					if math.Abs(currentObj.RelativeCoordinates.CenterY-prevObj.RelativeCoordinates.CenterY) < params.YThreshold {
						// TAG_SUCCESS case : a close enough co-ordinate was detected for a previously existing entry
						previousFrameData[idx].RelativeCoordinates = currentObj.RelativeCoordinates
						previousFrameData[idx].tagged = true
						tagged = true

						// TAG_SUCCESS case : increment the co-ordinates to the list
						perVehicleTrack[previousFrameData[idx].VehicleID].TrackPoints = append(perVehicleTrack[previousFrameData[idx].VehicleID].TrackPoints,
							CompactCoords{
								CenterX:    currentObj.RelativeCoordinates.CenterX,
								CenterY:    currentObj.RelativeCoordinates.CenterY,
								Confidence: currentObj.Confidence,
							})
						// TAG_SUCCESS case : increment the #frames for which object was tracked
						perVehicleTrack[previousFrameData[idx].VehicleID].FrameCount++
					}
				}
			}

			// Handle if object was untagged (new object detected)
			if !tagged {
				// SKIP : "traffic_light": 9, "person" : 0
				if currentObj.ClassID == 9 || currentObj.ClassID == 0 {
					continue
				}

				// TAG_FAILURE case : Add entry for new vehicleID in list of vehicle tracks
				perVehicleTrack = append(perVehicleTrack, VehicleTracks{
					VehicleID:  vehicleIDIndex,
					FrameCount: 1,
					ClassID:    currentObj.ClassID,
				})

				// TAG_FAILURE case : the vehicleID must exist in the perVehicleTrack arrays
				perVehicleTrack[vehicleIDIndex].TrackPoints = append(perVehicleTrack[vehicleIDIndex].TrackPoints, CompactCoords{
					CenterX:    currentObj.RelativeCoordinates.CenterX,
					CenterY:    currentObj.RelativeCoordinates.CenterY,
					Confidence: currentObj.Confidence,
				})

				tmpStruct := objectHistory{
					VehicleID:           vehicleIDIndex,
					ClassID:             currentObj.ClassID,
					Name:                currentObj.Name,
					RelativeCoordinates: currentObj.RelativeCoordinates,
					TagCounter:          0,
					tagged:              true,
				}

				previousFrameData = append(previousFrameData, tmpStruct)

				// In the end, Increment index for next vehicle ID
				vehicleIDIndex++
			}
		}

		// Increment if tagged and reset tag status
		for idx, previousObj := range previousFrameData {
			if previousObj.tagged {
				previousFrameData[idx].TagCounter += params.Rewards
			} else {
				previousFrameData[idx].TagCounter += params.Penalty
			}
			previousFrameData[idx].tagged = false
		}

		// Eliminate any entry which was not tagged recently
		previousFrameData, _ = Filter(previousFrameData, params.EliminateThreshold)

		// Add frame record to archive
		theArchive = append(theArchive, archiveRecord{
			FrameID:     i,
			FrameRecord: previousFrameData,
		})
	}

	// Write data to file
	if jsonString, err := json.MarshalIndent(theArchive, "", " "); err == nil {
		fmt.Println(outputDir + "/" + filepath)
		_ = ioutil.WriteFile(outputDir+"/"+filepath, jsonString, 0644)
	}

	// The following comments assume, that the sampling was done at 1 frame interval, for a 24 fps stream

	// Write data to file (vehicle track data)
	pruneAndWrite(outputDir, filepath, "_veh.json", perVehicleTrack, 0)

	// Test (pruned data - at least 10 frames) --> Noise
	pruneAndWrite(outputDir, filepath, "_veh_0.5s.json", perVehicleTrack, 3)

	// Test (pruned data - at least 60 frames) --> Half-second
	pruneAndWrite(outputDir, filepath, "_veh_1.0s.json", perVehicleTrack, 6)

	// Test (pruned data - at least 60 frames) --> One-second
	pruneAndWrite(outputDir, filepath, "_veh_2.0s.json", perVehicleTrack, 12)
}

func pruneAndWrite(outputPath, filepath, suffix string, data trackArchive, pruneParam int) {
	accepted, _ := PruneFalsePositives(data, pruneParam)
	if jsonString, err := json.MarshalIndent(accepted, "", " "); err == nil {
		vehicleDataPath := strings.TrimSuffix(filepath, path.Ext(filepath)) + suffix
		_ = ioutil.WriteFile(outputPath+"/"+vehicleDataPath, jsonString, 0644)
	}
}
