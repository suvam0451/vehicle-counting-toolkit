package traffic

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

// type trailData []trailDatum

// type mainArchive []archiveRecord

// type archiveRecord struct {
// 	FrameID     int         `json:"frame_id"`
// 	FrameRecord frameRecord `json:"objects"`
// }

// type frameRecord []objectHistory

// type trailDatum struct {
// 	FrameID  int64    `json:"frame_id"`
// 	Filename string   `json:"filename"`
// 	Objects  []object `json:"objects"`
// }

// type object struct {
// 	ClassID             int                 `json:"class_id"`
// 	Name                Name                `json:"name"`
// 	RelativeCoordinates relativeCoordinates `json:"relative_coordinates"`
// 	Confidence          float64             `json:"confidence"`
// }

// type objectHistory struct {
// 	EntryID             int                 `json:"id"`
// 	ClassID             int                 `json:"class_id"`
// 	Name                Name                `json:"name"`
// 	RelativeCoordinates relativeCoordinates `json:"relative_coordinates"`
// 	TagCounter          int                 `json:"tagcounter"`
// 	tagged              bool
// }

// Holds history of all past frames
type customArchive struct {
	frameID int
}

// CustomFrameData : Main container for our combined data segments
type CustomFrameData []CompactFrame

// CompactFrame : More compacted frame datas
type CompactFrame struct {
	frameID int
	objects []CompactObject
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

// type relativeCoordinates struct {
// 	CenterX float64 `json:"center_x"`
// 	CenterY float64 `json:"center_y"`
// 	Width   float64 `json:"width"`
// 	Height  float64 `json:"height"`
// }

// // Name : List of tag names in darknet
// type Name string

// const (
// 	Bicycle      Name = "bicycle"
// 	Bus          Name = "bus"
// 	Car          Name = "car"
// 	Motorbike    Name = "motorbike"
// 	Person       Name = "person"
// 	TrafficLight Name = "traffic light"
// 	Truck        Name = "truck"
// )

// // ModelParameters : Parameters for our model
// type ModelParameters struct {
// 	Upvote             int
// 	Downvote           int
// 	XThreshold         float64
// 	YThreshold         float64
// 	EliminateThreshold int
// }

// // DetectTrail detect trails for all files in given path
// func DetectTrail(inputpath string, params ModelParameters) {

// }

// DetectTrailCustom : I am modifying the codebase for more compact data format
func DetectTrailCustom(inputpath string, params ModelParameters) {

	reA := regexp.MustCompile(`^A_2_[0-9]_02\.json$`)
	// reB := regexp.MustCompile(`^B_2_[0-9]_02\.json$`)
	// reC := regexp.MustCompile(`^C_2_[0-9]_02\.json$`)
	// reF := regexp.MustCompile(`^F_2_[0-9]_02\.json$`)
	// reG := regexp.MustCompile(`^G_2_[0-9]_02\.json$`)

	// ounter := 1
	var inputfilespath []string
	var inputfilesname []string
	var wg sync.WaitGroup

	// Error already handled above
	filepath.Walk(inputpath,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() == false && reA.MatchString(info.Name()) {
				inputfilespath = append(inputfilespath, path)
				inputfilesname = append(inputfilesname, info.Name())
			}
			return nil
		})

	// Gnerate a single object from segmented files
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

	runAnalysis(SourceObject, params, "outputnew")
	wg.Wait()
}

func runAnalysis(source CustomFrameData, params ModelParameters, outpath string) {
	var previousFrameData []PreviousFrameObject
	var theArchive mainArchive
	var perVehicleTrack trackArchive
	var vehicleIDIndex int = 0

	for i, frame := range source {
		for _, currentobj := range frame.objects {
			tagged := false

			for idx, prevobj := range previousFrameData {
				if prevobj.ClassID == currentobj.ClassID && !prevobj.tagged {
					// TAG_SUCCESS case : a close enough co-ordinate was detected
					// for a previously existing entry
					previousFrameData[idx].CenterX = currentobj.CenterX
					previousFrameData[idx].CenterY = currentobj.CenterY
					previousFrameData[idx].confidence = currentobj.confidence
					previousFrameData[idx].tagged = true
					tagged = true

					foundID := previousFrameData[idx].ObjectID
					perVehicleTrack[foundID].TrackPoints = append()

				}
			}
			// for idx, prevobj
		}
	}
}
