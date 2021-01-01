import shapes

GREEN = (0, 255, 0)
RED = (0, 0, 255)
BLUE = (255, 0, 0)
YELLOW = (255, 255, 0)
GREY = (50, 50, 50)
import cv2 as cv


def hud_checkpoints(blank, data):
    for checks in data["checkpoints"]:
        start = checks["start"]
        end = checks["end"]
        shapes.draw_line(blank, start["x"], start["y"], end["x"], end["y"], GREY, 5)
