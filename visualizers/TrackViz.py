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

with open(join("./input","test_veh_05.json"), "r") as f:
    
    data = json.load(f)
    arr = np.asarray(data)
    print (arr[0]["objects"])
    # for vehicle_object in data:
        # print(vehicle_object)