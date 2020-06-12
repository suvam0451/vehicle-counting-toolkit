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

	inputfiles := []string{"outputnew/veh_A.json",
		"outputnew/veh_B.json",
		"outputnew/veh_D.json",
		"outputnew/veh_F.json",
		"outputnew/veh_G.json"}

	// Iterate over input files
	for _, inputfile := range inputfiles {
		if jsonFile, err := os.Open(inputfile); err == nil {
			// These need to be unique per file
			var data []FrameTaggedData
			var outdata []StackplotStruct

			_stat, _ := jsonFile.Stat()
			byteValue, _ := ioutil.ReadAll(jsonFile)
			json.Unmarshal(byteValue, &data)

			fmt.Println("Everything OK. Wait for program to exit...")
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
			// filepath.
			file, _ := json.MarshalIndent(outdata, "", " ")
			_ = ioutil.WriteFile("stackplot/"+str+"_stackplot.json", file, 0644)
		}
	}

}
