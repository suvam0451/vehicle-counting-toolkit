import cv2 as cv


# Gets dimensions of image {x, y}
def get_dims(blank):
    return blank.shape[1], blank.shape[0]


def screen_coords(blank, coords):
    dims = get_dims(blank)
    return int(dims[0] * coords[0]), int(dims[1] * coords[1])


def draw_circle(blank, x=0.5, y=0.5, radius=1, thickness=0.1, color=(0, 0, 255)):
    cv.circle(blank, (blank.shape[1] * x, blank.shape[0] * y), radius, color)


def draw_line(blank, xa=0.0, ya=0.0, xb=1.0, yb=1.0, color=(255, 255, 255), thickness=1):
    start = screen_coords(blank, (xa, ya))
    end = screen_coords(blank, (xb, yb))
    cv.line(blank, start, end, color, thickness)


def draw_text(blank, text, start=(0.5, 0.5), color=(0, 0, 255), scale=1, font=cv.FONT_HERSHEY_SIMPLEX, thickness=1):
    cv.putText(blank, text, screen_coords(blank, start), cv.FONT_HERSHEY_SIMPLEX, scale,
               color, thickness)


def draw_rect(blank, loc, color, offset=20, thickness=2):
    mid_point = screen_coords(blank, loc)
    start_point = (mid_point[0] - offset, mid_point[1] - offset)
    end_point = (mid_point[0] + offset, mid_point[1] + offset)
    cv.rectangle(blank, start_point, end_point, color, thickness)