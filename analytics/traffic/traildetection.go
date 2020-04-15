package traffic

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sync"
)

type trailData []trailDatum

type mainArchive []archiveRecord

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
	ClassID             int                 `json:"class_id"`
	Name                Name                `json:"name"`
	RelativeCoordinates relativeCoordinates `json:"relative_coordinates"`
	TagCounter          int                 `json:"tagcounter"`
	tagged              bool
}

type relativeCoordinates struct {
	CenterX float64 `json:"center_x"`
	CenterY float64 `json:"center_y"`
	Width   float64 `json:"width"`
	Height  float64 `json:"height"`
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

// func filter(vs []objectHistory, f func(objectHistory) bool) []objectHistory {
// 	vsf := make([]objectHistory, 0)
// 	for _, v := range vs {
// 		if f(v) {
// 			vsf = append(vsf, v)
// 		}
// 	}
// 	return vsf
// }

func filter(vs []objectHistory, threshold int) []objectHistory {
	vsf := make([]objectHistory, 0)
	for _, v := range vs {
		if v.TagCounter > threshold {
			vsf = append(vsf, v)
		}
	}
	return vsf
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

				// Capture name of json file and run async
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
	for i, frame := range data {
		for _, currentobj := range frame.Objects {
			tagged := false // will be set to true if object gets assigned to one of the previous frame objects

			for idx, prevobj := range previousFrameData {
				// ID match with untagged object
				if prevobj.ClassID == currentobj.ClassID && !prevobj.tagged {
					// Distance calculations
					if math.Abs(currentobj.RelativeCoordinates.CenterY-prevobj.RelativeCoordinates.CenterY) < params.YThreshold {
						previousFrameData[idx].RelativeCoordinates = currentobj.RelativeCoordinates
						previousFrameData[idx].tagged = true
						tagged = true
					}
				}
			}

			// Handle if object was untagged
			if !tagged {
				previousFrameData = append(previousFrameData, objectHistory{
					ClassID:             currentobj.ClassID,
					Name:                currentobj.Name,
					RelativeCoordinates: currentobj.RelativeCoordinates,
					TagCounter:          0,
					tagged:              true,
				})
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
		previousFrameData = filter(previousFrameData, params.EliminateThreshold)

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
	// Ensure all output paths exist...
	if _, err := os.Stat("./intermediate"); os.IsNotExist(err) {
		os.Mkdir("./intermediate", os.ModeDir)
	}

	// Write data to file
	if jsonString, err := json.MarshalIndent(theArchive, "", " "); err == nil {
		ioutil.WriteFile("./intermediate/"+filepath, jsonString, 0644)
	}
}
