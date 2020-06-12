"""Entry point."""

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
import cuda_lib as culib
import matplotlib.lines as mlines

style.use("ggplot")

# Parameter list
SEGMENTS = 8  # Number of segments we are dividing to


# COORD_LIST, TAG_DATA, REGION_DATA

@cuda.jit
def pick_segment(inArray, numSlice, SIZEY, outPartition):
    tx = cuda.threadIdx.x
    ty = cuda.blockIdx.x
    bw = cuda.blockDim.x

    # Compute flattened index inside the array
    pos = tx + ty * bw
    # Variable Initialization
    segment_size_x = float(SIZEY / numSlice)
    y_coord = inArray[pos][1]

    # Exit if Out Of Bound
    if pos > outPartition.size - 1:
        return

    for idx in range(numSlice):
        if(y_coord > (idx * segment_size_x) and y_coord < ((1 + idx) * segment_size_x)):
            outPartition[pos] = idx
            return


@cuda.jit
def greedy_pick(coord_list, tag_data, region_data, num_slice, intensity_param):
    """
    Assigns the input array to one of the sum_slice segments

    @intensity_param (float): Sets the greyscale mask value between (0,255) based on it's value
        and sets it to coord_list[idx][3]
    The assigned index is written to coord_list[idx][4]
    """

    # Calculate cuda index
    tx = cuda.threadIdx.x
    ty = cuda.blockIdx.x
    bw = cuda.blockDim.x

    pos = tx + ty * bw
    if pos > coord_list.size - 1:
        return

    # Initialization
    thread_object = coord_list[pos]
    tab_object = tag_data[pos]
    # region_data[pos].effect_power, region_data[pos].region
    region_object = region_data[pos]

    normalized_counter = min(max(tab_object[0] / intensity_param, 0.0), 1.0)

    if(normalized_counter == 0.0):
        region_object[0] = 255
        region_object[1] = 0

    effectpow = int(normalized_counter * intensity_param *
                    255)  # (0, 255) clamped
    region_object[1] = effectpow

    segment_size_y = 1.0 / num_slice
    y_coord = thread_object[1]        # Get Y coordinates

    for idx in range(num_slice):
        if(y_coord > (idx * segment_size_y) and y_coord < ((1 + idx) * segment_size_y)):
            region_object[0] = idx
            return


IMG = cv2.imread('images/image_02_01.jpg', 0)
SHAPE = (IMG.shape[1], IMG.shape[0])

# Converts screen-space position to pixel-space position


def screen_to_pixel(coords, size):
    return (math.floor(coords[0]*size[0]), math.floor(coords[1]*size[1]))

# Struct --> XCoord, YCoord, TagNum, ColorStrength, Segment
# MEGALIST = np.zeros((0, 1), dtype='float64, float64, int16, int16, uint8')


# Stores "list of co-ordinates" from json file
COORD_LIST = np.zeros((0, 2), dtype=np.float)
# Stores "tag counter" from json file
TAG_DATA = np.zeros((0, 1), dtype=np.int)
# FIXME : Not using this
MEGALIST = np.zeros((0, 5), dtype=np.float)

# zeroos = np.empty((0, 1), dtype='float32, float32, int16, uint8')
# zeroos = np.append(zeroos, np.array([[(1.0, 1.0, 1, 2)]], dtype=zeroos.dtype), axis=0)
# zeroos = np.append(zeroos, np.array([(1.0, 1.0, 1, 2)], dtype=zeroos.dtype), axis=0)

afk = np.array([[0, 0]])


# get all json files in "intermediate" folder
inputFileList = [f for f in listdir(
    "./intermediate") if isfile(join("./intermediate", f))]

# Invert Y axis exactly once for correct representation
plt.gca().invert_yaxis()


#
for inputfile in inputFileList:
    print(inputfile)
    with open(join("./intermediate", inputfile), "r") as f:
        distros_dict = json.load(f)
        for frame in distros_dict:
            for uobject in frame["objects"]:
                coords = uobject["relative_coordinates"]
                tag = uobject["tagcounter"]
                COORD_LIST = np.append(COORD_LIST,
                                       [[coords["center_x"],
                                         coords["center_y"]]], axis=0)
                TAG_DATA = np.append(TAG_DATA, [[tag]], axis=0)

                normalized_counter = min(
                    max(uobject["tagcounter"] / 512.0, 0.0), 1.0)
                if(normalized_counter == 0):
                    continue
                effectpow = normalized_counter * 255  # (0,255) clamped
                # screen_coordinates = (uobject["relative_coordinates"]["center_x"],
                #                     uobject["relative_coordinates"]["center_y"])
                # afk = np.append(afk, [screen_to_pixel(screen_coordinates, SHAPE)], axis=0)

    # All arrays have same size
    TPB = 32
    BPG = (TAG_DATA.size + (TPB - 1))

    # index #1 stores strength of trail, index #2 stores the segment to which the point belongs
    REGION_DATA = np.zeros((COORD_LIST.shape[0], 2), dtype=np.uint8)

    # CUDA version (fills REGION_DATA)
    greedy_pick[BPG, TPB](COORD_LIST, TAG_DATA, REGION_DATA,
                          SEGMENTS, 256.0)  # Call the CUDA function

    # print(REGION_DATA[95:100])

    colors = 1000 * ['g', 'r', 'c', 'b', 'k', 'm', 'y']

    #  Non CUDA version
    # lazy_pick(COORD_LIST, TAG_DATA, REGION_DATA,
    #           8, 256.0)  # Call the non-CUDA function

    for idx, point in enumerate(COORD_LIST, start=0):
        coloridx = REGION_DATA[idx][0]
        if coloridx == np.iinfo(np.uint8).max:
            continue
        plt.scatter(point[0], point[1],
                    marker="o", s=0.1, c=colors[REGION_DATA[idx][0]], linewidths=5)

    X_COORD_LIST = COORD_LIST[:, 0:1]
    Y_COORD_LIST = COORD_LIST[:, 1:2]

    #     print(condition)
    #     choices = [X_COORD_LIST]
    #     subarray = np.select(choices, condition)
    #     print(subarray)
    # int(REGION_DATA[idx][0])
    # np.save("output/a.npy", MEGALIST)
    # with open('output/your_file.txt', 'w') as f:
    #     for item in MEGALIST:
    #         f.write("%s\n" %item)

    # clf = KMeans()
    # clf.fit(afk)

    # Declaartions
    # centroids = clf.cluster_centers_

    # stream = cuda.stream()
    # d_partition = cuda.to_device(partition, stream)
    # pick_segment[BPG, TPB](afk, 8, SHAPEY, partition)  # CUDA test

    # d_partition.to_host(stream)
    # stream.synchronize()

    # Get the clusters for the current segment. Then draw on top of the image.
    for i in range(SEGMENTS):
        mask = REGION_DATA[:, 0] == i  # MASK_INFO --> falls under segment
        sampled_coords = COORD_LIST[mask, :]

        # Classifier on simplified ROI ------------------------> (TODO: Needs adjustments)
        clf = KMeans()
        clf.fit(sampled_coords)
        centroids = clf.cluster_centers_
        # ---------------------------------------------------------------

        # Connecting layers based on closeness (TODO: Complete the cpde)
        for j in range(centroids.shape[0]):
            plt.scatter(centroids[j][0], centroids[j][1],
                        marker="*", s=36, c="w", linewidths=5)
        # mlines.

    # Erase the graphs for the next round
    plt.clf()

    # Save the image to output folder by same name
    filename = inputfile[:len(inputfile)-5]
    plt.savefig(join("./output", filename + ".png"))
