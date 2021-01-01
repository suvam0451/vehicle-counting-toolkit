import cv2 as cv
import json
from collections import namedtuple
import shapes
import util
import overlay
import numpy as np

proj_id ="L2_V1"

# ksreddy1961@gmail.com

filepath = "/home/suvam/Downloads/"
filename = "TheGrind.webm"
datapath = "/home/suvam/Documents/MTP/"
data_1 = "data1/" + proj_id + ".json"

data_2 = "data2/" + proj_id + ".json"
# data_2 = "data2/" + proj_id + "_veh_0.5s" + ".json"
# data_2 = "data2/" + proj_id + "_veh_1.0s" + ".json"
# data_2 = "data2/" + proj_id + "_veh_2.0s" + ".json"

data_3 = "data3/" + proj_id + ".json"
video_path = "videos/" + proj_id + ".mp4"
image_path = "screens/L1_Original.png"

VehicleCountGroup = namedtuple("VehicleCountGroup", ["motorcycle", "car", "truck", "bus"])

GREEN = (0, 255, 0)
RED = (0, 0, 255)
BLUE = (255, 0, 0)
YELLOW = (255, 255, 0)
GREY = (50, 50, 50)

# Global variables (HUD references)
note_mouse_start = 0.0, 0.0
note_mouse_end = 0.0, 0.0
IS_MOUSE_DOWN = False
video_size = 0, 0
is_paused = False
# Global variables (object counter for frame)
frame_trucks = 0
frame_cars = 0
frame_buses = 0
frame_bikes = 0
# Global variables (object counter for frame)
global_cars = 0
global_bikes = 0
global_buses = 0
global_trucks = 0


def get_frame_data_as_string(idx):
    if idx == 0: return str(frame_cars)
    if idx == 1: return str(frame_bikes)
    if idx == 2: return str(frame_buses)
    if idx == 3: return str(frame_trucks)


def get_count_data_as_string(idx):
    if idx == 0: return str(global_cars)
    if idx == 1: return str(global_bikes)
    if idx == 2: return str(global_buses)
    if idx == 3: return str(global_trucks)


def get_dims(blank):
    return blank.shape[1], blank.shape[0]


def midpoint(blank):
    return blank.shape[1] / 2, blank.shape[0] / 2


def show_hud(blank, values):
    global is_paused
    DEFAULT_FONT_SCALE = 0.75
    DEFAULT_FONT_THICKNESS = 2
    hud_start = 0.025, 0.05
    def offset(x): return hud_start[0], hud_start[1] + x * 0.04
    def offset_section(x, y): return hud_start[0], hud_start[1] + y * 0.00 + x * 0.04

    def to_tally_text(x, y): return "Number of " + y + ": " + str(values[x])
    shapes.draw_text(blank, "Number of cars: " + get_frame_data_as_string(0), offset(0), YELLOW, DEFAULT_FONT_SCALE, thickness=DEFAULT_FONT_THICKNESS)
    shapes.draw_text(blank, "Number of bikes: " + get_frame_data_as_string(1), offset(1), YELLOW, DEFAULT_FONT_SCALE, thickness=DEFAULT_FONT_THICKNESS)
    shapes.draw_text(blank, "Number of trucks: " + get_frame_data_as_string(2), offset(2), YELLOW, DEFAULT_FONT_SCALE, thickness=DEFAULT_FONT_THICKNESS)
    shapes.draw_text(blank, "Number of buses: " + get_frame_data_as_string(3), offset(3), YELLOW, DEFAULT_FONT_SCALE, thickness=DEFAULT_FONT_THICKNESS)

    hud_start = 0.25, 0.05
    def to_frame_text(x, y): return y + " in frame: " + get_count_data_as_string(x)
    shapes.draw_text(blank, to_frame_text(0, "Cars"), offset(0), YELLOW, DEFAULT_FONT_SCALE, thickness=DEFAULT_FONT_THICKNESS)
    shapes.draw_text(blank, to_frame_text(1, "Bikes"), offset(1), YELLOW, DEFAULT_FONT_SCALE, thickness=DEFAULT_FONT_THICKNESS)
    shapes.draw_text(blank, to_frame_text(2, "Trucks"), offset(2), YELLOW, DEFAULT_FONT_SCALE, thickness=DEFAULT_FONT_THICKNESS)
    shapes.draw_text(blank, to_frame_text(3, "Buses"), offset(3), YELLOW, DEFAULT_FONT_SCALE, thickness=DEFAULT_FONT_THICKNESS)

    # Menu system (pausing/playing, drawling )
    menu_start = 0.70, 0.80
    def offset(x): return menu_start[0], menu_start[1] + x * 0.04
    def offset_section(x, y): return hud_start[0], hud_start[1] + y * 0.00 + x * 0.04

    if is_paused:
        shapes.draw_text(blank, "Press [P] to continue", offset(0), YELLOW, thickness=DEFAULT_FONT_THICKNESS)
    else:
        shapes.draw_text(blank, "Press [P] to pause", offset(0), YELLOW, thickness=DEFAULT_FONT_THICKNESS)
    shapes.draw_text(blank, "Press [Q] to quit", offset(1), YELLOW, thickness=DEFAULT_FONT_THICKNESS)
    shapes.draw_text(blank, "Press [R] to register offsets", offset(2), YELLOW, thickness=DEFAULT_FONT_THICKNESS)
    # shapes.draw_text(blank, "Press [ESC] to exit", offset(1), YELLOW, thickness=DEFAULT_FONT_THICKNESS)

    util_start = 0.85, 0.10
    UTIL_FONT_SCALE = 0.75
    UTIL_FONT_THICKNESS = 2
    def offset(x): return util_start[0], util_start[1] + x * 0.04
    shapes.draw_text(blank, "Start X: " + str(note_mouse_start[0]), offset(0), BLUE, UTIL_FONT_SCALE, thickness=UTIL_FONT_THICKNESS)
    shapes.draw_text(blank, "Start Y: " + str(note_mouse_start[1]), offset(1), BLUE, UTIL_FONT_SCALE, thickness=UTIL_FONT_THICKNESS)
    shapes.draw_text(blank, "End   X: " + str(note_mouse_end[0]), offset(2), BLUE, UTIL_FONT_SCALE, thickness=UTIL_FONT_THICKNESS)
    shapes.draw_text(blank, "End   Y: " + str(note_mouse_end[1]), offset(3), BLUE, UTIL_FONT_SCALE, thickness=UTIL_FONT_THICKNESS)


def paused_hud(blank):
    hud_start = 0.45, 0.45
    shapes.draw_text(blank, "Paused. Press [P] to continue...", hud_start, GREEN)


def frame_wise_overlay(blank, data):
    DEFAULT_FONT_SCALE = 0.75
    DEFAULT_FONT_THICKNESS = 2
    print(len(data["objects"]))
    for obj in data["objects"]:
        coords = obj["relative_coordinates"]
        center = coords["center_x"], coords["center_y"]
        class_id = obj["class_id"]
    # shapes.draw_text(blank, str(obj["id"]), center, BLUE, DEFAULT_FONT_SCALE, thickness=DEFAULT_FONT_THICKNESS)
    if class_id == 2:  # Car
        shapes.draw_text(blank, str(obj["id"]), center, RED, DEFAULT_FONT_SCALE, thickness=DEFAULT_FONT_THICKNESS)
    elif class_id == 3:  # Car
        shapes.draw_text(blank, str(obj["id"]), center, GREEN, DEFAULT_FONT_SCALE, thickness=DEFAULT_FONT_THICKNESS)
    elif class_id == 5:  # Car
        shapes.draw_text(blank, str(obj["id"]), center, YELLOW, DEFAULT_FONT_SCALE, thickness=DEFAULT_FONT_THICKNESS)
    elif class_id == 7:  # Car
        shapes.draw_text(blank, str(obj["id"]), center, BLUE, DEFAULT_FONT_SCALE, thickness=DEFAULT_FONT_THICKNESS)


frame_counter = 0
data_counter = 0
capture_interval = 1
# Data handles
classifier_data = []
checkpoint_data = {}
frame_wise_data = []

# Classifier data: Frame-wise object detected
with open(datapath + data_1) as f:
    classifier_data = json.load(f)

with open(datapath + data_2) as f:
    frame_wise_data = json.load(f)

with open(datapath + data_3) as f:
    checkpoint_data = json.load(f)


def hud_boxes(blank, data, idx):
    global frame_trucks, frame_bikes, frame_cars, frame_buses
    # Reset
    frame_trucks = 0
    frame_bikes = 0
    frame_buses = 0
    frame_cars = 0
    # Draw bounds and refill
    for obj in data[idx]["objects"]:
        coords = obj["relative_coordinates"]
        if obj["class_id"] == 2:  # car
            frame_cars += 1
            shapes.draw_rect(blank, (coords["center_x"], coords["center_y"]), RED)
        elif obj["class_id"] == 3:  # motorbike
            frame_bikes += 1
            shapes.draw_rect(blank, (coords["center_x"], coords["center_y"]), GREEN)
        elif obj["class_id"] == 5:  # bus
            frame_buses += 1
            shapes.draw_rect(blank, (coords["center_x"], coords["center_y"]), YELLOW)
        elif obj["class_id"] == 7:  # truck
            frame_trucks += 1
            shapes.draw_rect(blank, (coords["center_x"], coords["center_y"]), BLUE)


def handle_mouse_click(event, x, y, flags, param):
    global note_mouse_start, note_mouse_end, video_size, IS_MOUSE_DOWN
    if is_paused:
        if event == cv.EVENT_LBUTTONDOWN:
            IS_MOUSE_DOWN = True
            note_mouse_start = util.pixel_to_unit_coord(video_size, (x, y))
        elif event == cv.EVENT_LBUTTONUP:
            IS_MOUSE_DOWN = False
            note_mouse_end = util.pixel_to_unit_coord(video_size, (x, y))


cv.namedWindow('frame')
# Event Callbacks
cv.setMouseCallback('frame', handle_mouse_click)

# Main loop
cap = cv.VideoCapture(datapath + video_path)
video_size = cap.get(3), cap.get(4)


while cap.isOpened():
    if not is_paused:
        ret, frame = cap.read()
        tmp = VehicleCountGroup(1, 2, 3, 4)
        fps = cv.CAP_PROP_FPS
        # Draw classifier data
        hud_boxes(frame, classifier_data, data_counter)

        # Draw checkpoints
        overlay.hud_checkpoints(frame, checkpoint_data)
        frame_wise_overlay(frame, frame_wise_data[frame_counter])

        show_hud(frame, tmp)
        cv.imshow('frame', frame)
        frame_counter += 1
        data_counter += 1
        if cv.waitKey(1) & 0xFF == ord('p'):
            is_paused = True
        if cv.waitKey(1) & 0xFF == ord('q'):
            break
    else: # Paused state
        paused_hud(frame)
        if cv.waitKey(1) & 0xFF == ord('p'):
            is_paused = False
        if cv.waitKey(1) & 0xFF == ord('q'):
            break


cap.release()
cv.destroyAllWindows()
