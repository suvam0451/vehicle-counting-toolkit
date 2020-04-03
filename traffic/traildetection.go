package traffic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
)

type trailData []trailDatum

type mainArchive []archiveRecord

type archiveRecord struct {
	FrmaeID     int
	FrameRecord frameRecord
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
	ClassID             int
	Name                Name
	RelativeCoordinates relativeCoordinates
	TagCounter          int
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

// DetectTrail detect trails for all files in given path
func DetectTrail(inputpath string, params ModelParameters) {
	// ounter := 1
	var inputfiles []string

	// Error already handled above
	filepath.Walk(inputpath,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() == false {
				inputfiles = append(inputfiles, path)
			}
			return nil
		})
	fmt.Println(inputfiles)

	var ParsedStruct trailData
	for _, file := range inputfiles {
		if openfile, err := os.Open(file); err == nil {
			byteValue, _ := ioutil.ReadAll(openfile)
			if err := json.Unmarshal(byteValue, &ParsedStruct); err == nil {
				detectIndividualTrail(ParsedStruct, params)
			}
		}
	}
}

func detectIndividualTrail(data trailData, params ModelParameters) {
	var previousFrames frameRecord
	var theArchive mainArchive
	for i, frame := range data {
		for _, currentobj := range frame.Objects {
			tagged := false // will be set to true if object gets assigned to one of the previous frame objects

			for _, prevobj := range previousFrames {
				// ID match with untagged object
				if prevobj.ClassID == currentobj.ClassID && prevobj.tagged == false {
					// Distance calculations
					if math.Abs(currentobj.RelativeCoordinates.CenterX-prevobj.RelativeCoordinates.CenterX) < params.XThreshold {
						prevobj.RelativeCoordinates = currentobj.RelativeCoordinates
						prevobj.tagged = true
						tagged = true
					}
				}
			}

			// Handle if object was untagged
			if !tagged {
				previousFrames = append(previousFrames, objectHistory{
					ClassID:             currentobj.ClassID,
					Name:                currentobj.Name,
					RelativeCoordinates: currentobj.RelativeCoordinates,
					TagCounter:          0,
					tagged:              true,
				})
			}
		}

		var temporary frameRecord
		// Update counters
		for _, prevobj := range previousFrames {
			temporary = append(temporary, objectHistory{
				ClassID:             prevobj.ClassID,
				Name:                prevobj.Name,
				RelativeCoordinates: prevobj.RelativeCoordinates,
				tagged:              false,
				TagCounter:          prevobj.TagCounter,
			})
		}

		// Add frame record to archive
		theArchive = append(theArchive, archiveRecord{
			FrmaeID:     i,
			FrameRecord: temporary,
		})
	}
}
