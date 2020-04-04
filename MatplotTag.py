import cv2
import json, math
import matplotlib.pyplot as plt
from matplotlib import style
import numpy as np
from sklearn.cluster import KMeans
from os import listdir
from os.path import isfile, join

style.use("ggplot")

img = cv2.imread('images/image_02_02.jpg', 0)
shape = (img.shape[1], img.shape[0])

# Converts screen-space position to pixel-space position
def screen_to_pixel(coords, size):
    return (math.floor(coords[0]*size[0]), math.floor(coords[1]*size[1]))

mainaray = np.array([])
afk = np.array([[0,0]])


# get all intermediate json files
inputFileList = [f for f in listdir("./intermediate") if isfile(join("./intermediate", f))]

# 
for inputfile in inputFileList:
    with open(join("./intermediate", inputfile), "r") as f:
        distros_dict = json.load(f)
        for frame in distros_dict:
            for uobject in frame["objects"]:
                normalized_counter = min(max(uobject["tagcounter"] / 1024, 0), 1)
                if(normalized_counter == 0):
                    continue
                effectpow = normalized_counter * 255 # (0,255) clamped
                screen_coordinates = (uobject["relative_coordinates"]["center_x"], uobject["relative_coordinates"]["center_y"])
                # afk = np.append(afk, [[uobject["relative_coordinates"]["center_x"], uobject["relative_coordinates"]["center_y"]]], axis=0)
                afk = np.append(afk, [screen_to_pixel(screen_coordinates, shape)], axis=0)
    clf = KMeans()
    clf.fit(afk)

    centroids = clf.cluster_centers_

    colors =  10 * ["g.","r.", "c.","b.","k.","m.", "y."]
    # print(len(afk))

    # Number of horizontal slices to segment the image space to...
    n_slices = 16

    for coord_set in afk:
        # Check which segment the point falls to, and color accordingly, coord_set[1] corresponds to y coordinate
        for i in range(0, n_slices - 1):
            if(coord_set[1] > (i * (img.shape[0] / n_slices)) and coord_set[1] < ((1 + i) * (img.shape[0] / n_slices))):
                plt.plot(coord_set[0], coord_set[1], colors[i], markersize=5) # labels --> 0/1/2/3

    plt.gca().invert_yaxis()
    plt.scatter(centroids[:,0], centroids[:,1], marker="x", s=150, linewidths=5)
    # Save the image to output folder by same name
    filename = inputfile[:len(inputfile)-5]
    plt.savefig(join("./output", filename + ".png"))
    # plt.show()
    print("processed:", inputfile)