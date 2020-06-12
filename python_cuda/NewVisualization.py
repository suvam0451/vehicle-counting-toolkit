import cv2
import json
import math
from os import listdir
from os.path import isfile, join


def pixel_to_screen(coords, size):
    return (coords[0]/size[0], coords[1]/size[1])

# Converts screen-space position to pixel-space position


def screen_to_pixel(coords, size):
    return (math.floor(coords[0]*size[0]), math.floor(coords[1]*size[1]))

    img = cv2.imread('images/image_02_02.jpg', 0)
