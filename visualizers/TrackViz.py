"""Test file to visualize detected trail lines from videos"""

import cv2
import json
import math
import time
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
plt.xlim(0,1)
plt.ylim(1,0)
# matplotlib axis/bg settings
img = cv2.imread("./images/A.jpg")

bins = np.fromiter((i*10 for i in range(100)), dtype="float32")


FRAME_COUNTERS = np.zeros((0, 1), dtype=np.float)

with open(join("./input","test_veh_05.json"), "r") as f:
    data = json.load(f)

    for tracked_vehicle in data:
        # Stores "list of co-ordinates" from json file
        COORD_LIST = np.zeros((0, 2), dtype=np.float)
        FRAME_COUNTERS = np.append(FRAME_COUNTERS, [[tracked_vehicle["frame_count"]]], axis=0)
        for coordinates in tracked_vehicle["objects"]:
            COORD_LIST = np.append(COORD_LIST,
                        [[coordinates["center_x"],
                            coordinates["center_y"]]], axis=0)
        # print(FRAME_COUNTERS)
        plt.scatter(COORD_LIST[:,0:1], COORD_LIST[:,1:2])

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