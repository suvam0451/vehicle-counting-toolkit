package traffic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// StackplotStruct : Holds history of all past frames
type StackplotStruct struct {
	Car       int `json:"car"`
	Bicycle   int `json:"bicycle"`
	Motorbike int `json:"motorbike"`
	Truck     int `json:"truck"`
}

// FrameTaggedData : Newest iteration
type FrameTaggedData struct {
	FrameID int              `json:"frame_id"`
	Objects []FrameObjectNew `json:"frames"` // FIXME: This shoudl be objetcs
}

// FrameObjectNew : Newest iteration
type FrameObjectNew struct {
	ID         int     `json:"id"`
	ClassID    int     `json:"class_id"`
	CenterX    float64 `json:"center_x"`
	CenterY    float64 `json:"center_y"`
	TagCounter int     `json:"tagcounter"`
}

// GenerateStackplot : Generates stackplot
func GenerateStackplot() {
	// folder with output files from "traildetection_alt" command
	basedir := "./out_traildetection_alt"
	outputDir := "./out_stackplot"

	// Create folder if not exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModeDir)
	}

	inputfiles := []string{basedir + "/veh_A.json", basedir + "/veh_B.json",
	basedir + "/veh_D.json", basedir + "/veh_F.json", basedir + "/veh_G.json"}

	numFiles := len(inputfiles)
	// Iterate over input files
	for i, inputfile := range inputfiles {
		if jsonFile, err := os.Open(inputfile); err == nil {
			// These need to be unique per file
			var data []FrameTaggedData
			var outdata []StackplotStruct

			_stat, _ := jsonFile.Stat()
			byteValue, _ := ioutil.ReadAll(jsonFile)
			json.Unmarshal(byteValue, &data)

			// print status to terminal
			outStr := fmt.Sprintf("Processing... %d/%d files", i + 1, numFiles)
			fmt.Println(outStr)

			// Iterate frames
			for _, framedata := range data {
				var tmp StackplotStruct
				// Iterate frame objects
				for _, objdata := range framedata.Objects {
					if objdata.ClassID == 7 {
						tmp.Truck++
					} else if objdata.ClassID == 2 {
						tmp.Car++
					} else if objdata.ClassID == 3 {
						tmp.Motorbike++
					}
				}
				outdata = append(outdata, tmp)
			}

			str := strings.TrimSuffix(_stat.Name(), filepath.Ext(_stat.Name()))
			// output filepath
			file, _ := json.MarshalIndent(outdata, "", " ")
			_ = ioutil.WriteFile("out_stackplot/"+str+"_stackplot.json", file, 0644)
		}
	}

}
