import cv2 as cv


def pixel_to_unit_coord(video_size, dims):
    return dims[0] / video_size[0], dims[1] / video_size[1]