### Intermediate folder

`_veh` files contain id tracked path trajectories of vehicles.

`_veh_x.xs` contain the vehicles which were continously detected for
that much duration.

### outputnew folder

Contains all the actual output files from the latest iteration.

`veh_X_c.json` -- ID tagged vehicle objects
`veh_X.json` -- Frame tagged objects

### input_traildetect

This folder has the following formatted data

```json
{
 "frame_id":1, 
 "filename":"/home/suvam/Documents/MTP/data/out_tinyyolo/L2_V1/00001.jpg", 
 "objects": [ 
  {"class_id":0, "name":"person", "relative_coordinates":{"center_x":0.235643, "center_y":0.605867, "width":0.036381, "height":0.084796}, "confidence":0.397950}, 
  {"class_id":3, "name":"motorbike", "relative_coordinates":{"center_x":0.235643, "center_y":0.605867, "width":0.036381, "height":0.084796}, "confidence":0.304368}, 
  {"class_id":2, "name":"car", "relative_coordinates":{"center_x":0.360789, "center_y":0.561696, "width":0.040099, "height":0.085848}, "confidence":0.950816}, 
  {"class_id":2, "name":"car", "relative_coordinates":{"center_x":0.291200, "center_y":0.572183, "width":0.058488, "height":0.061687}, "confidence":0.855379}, 
  {"class_id":2, "name":"car", "relative_coordinates":{"center_x":0.631011, "center_y":0.550167, "width":0.041236, "height":0.046462}, "confidence":0.350822}, 
  {"class_id":2, "name":"car", "relative_coordinates":{"center_x":0.450887, "center_y":0.516193, "width":0.035826, "height":0.066644}, "confidence":0.308274}, 
  {"class_id":2, "name":"car", "relative_coordinates":{"center_x":0.475797, "center_y":0.514897, "width":0.056791, "height":0.056837}, "confidence":0.254386}
 ] 
}
```