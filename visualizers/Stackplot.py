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
import os
import matplotlib.pyplot as plt
from matplotlib import style
import numpy as np
from sklearn.cluster import KMeans
from os import listdir
from os.path import isfile, join
from numba import jit, float32, cuda


VIDEO_FPS = 25.0  # float
SAMPLE_INTERVAL = 5

# fileArray = [
# "veh_A_stackplot.json",
# "veh_B_stackplot.json",
# "veh_D_stackplot.json",
# "veh_F_stackplot.json",
# "veh_G_stackplot.json"
# ]


inputdir = "./input_stackplot"
outputdir = "./out_stackplot"

# Make directory if not exist
if not os.path.exists(outputdir):
    os.makedirs(outputdir)

# get all json files in target folder
fileList = [f for f in listdir(inputdir) if isfile(join(inputdir, f))]

# x = [1, 2, 3, 4, 5]


INCREMENT_VALUE = SAMPLE_INTERVAL / VIDEO_FPS

# Main loop
for inputfile in fileList:
    # local variables
    x = np.zeros(0, dtype=np.int)
    y1 = np.zeros(0, dtype=np.int)  # Cars
    y2 = np.zeros(0, dtype=np.int)  # Motorbike
    y3 = np.zeros(0, dtype=np.int)  # Truck
    counter = 0.0

    with open(join(inputdir, inputfile), "r") as f:
        distros_dict = json.load(f)

        for entry in distros_dict:
            y1 = np.append(y1, [entry["car"]])
            y2 = np.append(y2, [entry["motorbike"]])
            y3 = np.append(y3, [entry["truck"]])
            x = np.append(x, [counter])
            counter = counter + INCREMENT_VALUE  # Increment counter after use

    # Extarct filename
    filename = inputfile[:len(inputfile)-5]

    # Subplot 1 --> Stackplot of {Cars, Motorbikes, Trucks}
    fig, ax1 = plt.subplots()
    ax1.title.set_text("Combined total active vehicles vs Time")
    ax1.set_xlabel("time(seconds)", fontsize=14)
    ax1.set_ylabel("Number", fontsize=14)
    y = np.vstack([y1, y2, y3])
    labels = ["Cars ", "Motorbikes", "Truck"]
    ax1.stackplot(x, y1, y2, y3, labels=labels)
    ax1.legend(loc='upper left')
    plt.savefig(join("./stackplots", filename + "_01.png"))
    ax1.cla()

    # Subplot 2 --> Stackplot of {Cars, Motorbikes, Trucks}
    fig, ax2 = plt.subplots()
    ax2.title.set_text("Active vehicles vs Time")
    ptr1 = ax2.plot(x, y1, label="Cars")
    ptr2 = ax2.plot(x, y2, label="Motorbikes")
    ptr3 = ax2.plot(x, y3, label="Truck")
    ax2.set_xlabel("time(seconds)", fontsize=20)
    ax2.set_ylabel("Number", fontsize=20)
    ax1.legend(loc='upper left')
    # Size and save
    fig.set_size_inches(16, 8, forward=True)  # play with size
    plt.savefig(join(outputdir, filename + "_02.png"))
    ax2.cla()

    # Clear figure for next round
    plt.clf()
