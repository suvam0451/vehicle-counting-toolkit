/**
Outputs the following :-

1. filename.json --
2. filename_veh.json -- Tracked frames for each vehicle tag. Entries with low score are eliminated.
*/

package traffic

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type trailData []trailDatum

type mainArchive []archiveRecord

type trackArchive []VehicleTracks

// VehicleTracks :  List of positions w/confidence through which vehicle has passed
type VehicleTracks struct {
	VehicleID   int             `json:"vehicle_id"`  // ID given to the vehicle
	FrameCount  int             `json:"frame_count"` // Number of frames for which this object was detected
	ClassID     int             `json:"class_id"`    // ClassID for this vehicle tyype
	TrackPoints []CompactCoords `json:"objects"`     // List of co-ordinates
	// TrackPoints []object `json:"objects"`     // List of co-ordinates
}

// Holds record for all the previous frames
type archiveRecord struct {
	FrameID     int         `json:"frame_id"`
	FrameRecord frameRecord `json:"objects"`
}

type frameRecord []objectHistory

type trailDatum struct {
	FrameID  int64    `json:"frame_id"`
	Filename string   `json:"filename"`
	Objects  []object `json:"objects"`
}

type object struct {
	ClassID             int                 `json:"class_id"`
	Name                Name                `json:"name"`
	RelativeCoordinates relativeCoordinates `json:"relative_coordinates"`
	Confidence          float64             `json:"confidence"`
}

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

const (
	Bicycle      Name = "bicycle"
	Bus          Name = "bus"
	Car          Name = "car"
	Motorbike    Name = "motorbike"
	Person       Name = "person"
	TrafficLight Name = "traffic light"
	Truck        Name = "truck"
)

// ModelParameters : Parameters for our model
type ModelParameters struct {
	Upvote             int
	Downvote           int
	XThreshold         float64
	YThreshold         float64
	EliminateThreshold int
}

// Filters struct array for minimum threshold value
func filter(vs []objectHistory, threshold int) ([]objectHistory, []objectHistory) {
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
	// accepted := make([]vehicleTracks, 0)
	// rejected := make([]vehicleTracks, 0)
	for _, v := range archive {
		if v.FrameCount > minThreshold {
			accepted = append(accepted, v)
		} else {
			rejected = append(rejected, v)
		}
	}
	return
}

// DetectTrail detect trails for all files in given path
func DetectTrail(inputpath string, params ModelParameters) {
	// ounter := 1
	var inputfilespath []string
	var inputfilesname []string
	var wg sync.WaitGroup

	// Error already handled above
	filepath.Walk(inputpath,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() == false {
				inputfilespath = append(inputfilespath, path)
				inputfilesname = append(inputfilesname, info.Name())
			}
			return nil
		})

	var ParsedStruct trailData
	for i, file := range inputfilespath {
		if openfile, err := os.Open(file); err == nil {
			byteValue, _ := ioutil.ReadAll(openfile)
			if err := json.Unmarshal(byteValue, &ParsedStruct); err == nil {
				wg.Add(1)

				// Multi-threaded processing of input files
				go func(filename string) {
					detectIndividualTrail(ParsedStruct, params, filename)
					wg.Done()
				}(inputfilesname[i])
			}
		}
	}

	wg.Wait()
}

func detectIndividualTrail(data trailData, params ModelParameters, filepath string) {
	var previousFrameData frameRecord
	var theArchive mainArchive
	var perVehicleTrack trackArchive
	var vehicleIDIndex int = 0

	for i, frame := range data {
		for _, currentobj := range frame.Objects {
			tagged := false // will be set to true if object gets assigned to one of the previous frame objects

			for idx, prevobj := range previousFrameData {
				// ID match with untagged object
				if prevobj.ClassID == currentobj.ClassID && !prevobj.tagged {
					// Distance calculations
					if math.Abs(currentobj.RelativeCoordinates.CenterY-prevobj.RelativeCoordinates.CenterY) < params.YThreshold {
						// TAG_SUCCESS case : a close enough co-ordinate was detected for a previously existing entry
						previousFrameData[idx].RelativeCoordinates = currentobj.RelativeCoordinates
						previousFrameData[idx].tagged = true
						tagged = true

						// TAG_SUCCESS case : increment the co-ordinates to the list
						perVehicleTrack[previousFrameData[idx].VehicleID].TrackPoints = append(perVehicleTrack[previousFrameData[idx].VehicleID].TrackPoints,
							CompactCoords{
								CenterX:    currentobj.RelativeCoordinates.CenterX,
								CenterY:    currentobj.RelativeCoordinates.CenterY,
								Confidence: currentobj.Confidence,
							})
						// TAG_SUCCESS case : increment the #frames for which object was tracked
						perVehicleTrack[previousFrameData[idx].VehicleID].FrameCount++
					}
				}
			}

			// Handle if object was untagged (new object detected)
			if !tagged {
				// SKIP : "traffic_light": 9, "person" : 0
				if currentobj.ClassID == 9 || currentobj.ClassID == 0 {
					continue
				}

				// TAG_FAILURE case : Add entry for new vehicleID in list of vehicle tracks
				perVehicleTrack = append(perVehicleTrack, VehicleTracks{
					VehicleID:  vehicleIDIndex,
					FrameCount: 1,
					ClassID:    currentobj.ClassID,
				})

				// TAG_FAILURE case : the vehicleID must exist in the perVehicleTrack arrays
				perVehicleTrack[vehicleIDIndex].TrackPoints = append(perVehicleTrack[vehicleIDIndex].TrackPoints, CompactCoords{
					CenterX:    currentobj.RelativeCoordinates.CenterX,
					CenterY:    currentobj.RelativeCoordinates.CenterY,
					Confidence: currentobj.Confidence,
				})

				tmpStruct := objectHistory{
					VehicleID:           vehicleIDIndex,
					ClassID:             currentobj.ClassID,
					Name:                currentobj.Name,
					RelativeCoordinates: currentobj.RelativeCoordinates,
					TagCounter:          0,
					tagged:              true,
				}

				previousFrameData = append(previousFrameData, tmpStruct)

				// In the end, Increment index for next vehicle ID
				vehicleIDIndex++
			}
		}

		// Increment if tagged and reset tag status
		for idx, prevobj := range previousFrameData {
			if prevobj.tagged {
				previousFrameData[idx].TagCounter += params.Upvote
			} else {
				previousFrameData[idx].TagCounter += params.Downvote
			}
			previousFrameData[idx].tagged = false
		}

		// Eliminate any entry which was not tagged recently
		previousFrameData, _ = filter(previousFrameData, params.EliminateThreshold)

		// Add frame record to archive
		theArchive = append(theArchive, archiveRecord{
			FrameID:     i,
			FrameRecord: previousFrameData,
		})
	}

	// Ensure all output paths exist...
	if _, err := os.Stat("./output"); os.IsNotExist(err) {
		os.Mkdir("./output", os.ModeDir)
	}
	if _, err := os.Stat("./intermediate"); os.IsNotExist(err) {
		os.Mkdir("./intermediate", os.ModeDir)
	}

	// Write data to file
	if jsonString, err := json.MarshalIndent(theArchive, "", " "); err == nil {
		ioutil.WriteFile("./intermediate/"+filepath, jsonString, 0644)
	}

	// Write data to file (vehicle track data)
	if jsonString, err := json.MarshalIndent(perVehicleTrack, "", " "); err == nil {
		// Filename modified to xyz_vehicles.json
		vechicleDataPath := strings.TrimSuffix(filepath, path.Ext(filepath)) + "_veh.json"
		ioutil.WriteFile("./intermediate/"+vechicleDataPath, jsonString, 0644)
	}

	// Test (pruned data - at least 10 frames) --> Noise
	accepted, _ := PruneFalsePositives(perVehicleTrack, 3)
	if jsonString, err := json.MarshalIndent(accepted, "", " "); err == nil {
		vechicleDataPath := strings.TrimSuffix(filepath, path.Ext(filepath)) + "_veh_0.5s.json"
		ioutil.WriteFile("./intermediate/"+vechicleDataPath, jsonString, 0644)
	}

	// Test (pruned data - at least 60 frames) --> Half-second
	accepted, _ = PruneFalsePositives(perVehicleTrack, 6)
	if jsonString, err := json.MarshalIndent(accepted, "", " "); err == nil {
		vechicleDataPath := strings.TrimSuffix(filepath, path.Ext(filepath)) + "_veh_1.0s.json"
		ioutil.WriteFile("./intermediate/"+vechicleDataPath, jsonString, 0644)
	}

	// Test (pruned data - at least 60 frames) --> Half-second
	accepted, _ = PruneFalsePositives(perVehicleTrack, 12)
	if jsonString, err := json.MarshalIndent(accepted, "", " "); err == nil {
		vechicleDataPath := strings.TrimSuffix(filepath, path.Ext(filepath)) + "_veh_2.0s.json"
		ioutil.WriteFile("./intermediate/"+vechicleDataPath, jsonString, 0644)
	}
}
