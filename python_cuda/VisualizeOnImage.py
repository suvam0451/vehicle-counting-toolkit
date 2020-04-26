import cv2
import json, math
from os import listdir
from os.path import isfile, join

def pixel_to_screen(coords, size):
    return (coords[0]/size[0], coords[1]/size[1])

# Converts screen-space position to pixel-space position
def screen_to_pixel(coords, size):
    return (math.floor(coords[0]*size[0]), math.floor(coords[1]*size[1]))

img = cv2.imread('images/image_02_02.jpg', 0)

radius = 5
thickness = 2
shape = (img.shape[1], img.shape[0])

# get all intermediate json files
inputFileList = [f for f in listdir("./intermediate") if isfile(join("./intermediate", f))]

for inputfile in inputFileList:
    with open("repo/out_02_04.json", "r") as f:
        distros_dict = json.load(f)
        for frame in distros_dict:
            for uobject in frame["objects"]:
                normalized_counter = min(max(uobject["tagcounter"] / 1024, 0), 1)
                if(normalized_counter == 0):
                    continue
                # print("pass for: ", uobject["tagcounter"], uobject["relative_coordinates"]["center_x"], uobject["relative_coordinates"]["center_y"])
                effectpow = normalized_counter * 255 # (0,255) clamped
                screen_coordinates = (uobject["relative_coordinates"]["center_x"], uobject["relative_coordinates"]["center_y"])
                pixel_coordinates = screen_to_pixel(screen_coordinates, shape)
                color = (effectpow, 0, 0)
                img = cv2.circle(img, pixel_coordinates, radius, color, thickness)

    cv2.imwrite('images/image_02_04_out.jpg',img)

