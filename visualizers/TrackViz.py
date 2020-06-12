"""Test file to visualize detected trail lines from videos"""

# Usage --> python Trackviz.py 3
#     0 --> Amtala
#     1 --> Bamoner
#     2 --> Diamond
#     3 --> Fotepore
#     4 --> Gangasagar

import cv2
import json
import math
import time
import sys
import matplotlib.pyplot as plt
from matplotlib import style
import numpy as np
from sklearn.cluster import KMeans
from os import listdir
from os.path import isfile, join
from numba import jit, float32, cuda

style.use("ggplot")
plt.title("Tracking distances")
plt.xlabel("Plot Number")
plt.ylabel("Plot points")
plt.xlim(0, 1)
plt.ylim(1, 0)
# plt.gca().invert_yaxis()

# matplotlib axis/bg settings
images = np.asarray(["../images/Sample_Amtala.jpg",
                     "../images/Sample_Bamoner.jpg",
                     "../images/Sample_Diamond.jpg",
                     "../images/Sample_Fotepore.jpg",
                     "../images/Sample_Gangasagar.jpg"])
json_files_track = np.asarray([
    "./inputnew/veh_A_c.json",
    "./inputnew/veh_B_c.json",
    "./inputnew/veh_D_c.json",
    "./inputnew/veh_F_c.json",
    "./inputnew/veh_G_c.json"
])
json_files_frames = np.asarray([
    "./inputnew/veh_A.json",
    "./inputnew/veh_B.json",
    "./inputnew/veh_D.json",
    "./inputnew/veh_F.json",
    "./inputnew/veh_G.json"
])

## Modify 
targetindex = 2
## Primary variables
image_to_open = images[int(sys.argv[1])]
file_to_open = json_files_track[int(sys.argv[1])]
## ----------------------------------------------


img = cv2.imread(image_to_open)
bins = np.fromiter((i*10 for i in range(100)), dtype="float32")

# Setup sub-plot
fig, ax = plt.subplots()
plt.imshow(img, extent=[0, 1, 1, 0])

FRAME_COUNTERS = np.zeros((0, 1), dtype=np.float)

with open(file_to_open, "r") as f:
    data = json.load(f)

    for tracked_vehicle in data:
        # Stores "list of co-ordinates" from json file
        COORD_LIST = np.zeros((0, 2), dtype=np.float)
        FRAME_COUNTERS = np.append(
            FRAME_COUNTERS, [[tracked_vehicle["frame_count"]]], axis=0)
        for coordinates in tracked_vehicle["objects"]:
            COORD_LIST = np.append(COORD_LIST,
                                   [[coordinates["center_x"],
                                     coordinates["center_y"]]], axis=0)
        # print(FRAME_COUNTERS)
        ax.scatter(COORD_LIST[:, 0:1], COORD_LIST[:, 1:2])
        # plt.scatter(COORD_LIST[:,0:1], COORD_LIST[:,1:2])

# plt.savefig(join("./output", "output.png"))
plt.savefig(join("./output", "output.png"))

plt.show()
plt.clf()
plt.hist(FRAME_COUNTERS, bins, histtype="bar", rwidth=0.75)
plt.savefig(join("./output", "track_length.png"))
# print(COORD_LIST[:,0:1])

# secarr = np.asarray(arr[0]["objects"][1])

# lookup = json.JSONDecoder().decode(secarr)
# print (lookup)
# for vehicle_object in data:
# print(vehicle_object)
