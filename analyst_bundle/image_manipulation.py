import cv2 as cv


def blur_uniform(blank, radius=3):
    cv.GaussianBlur(blank, (radius, radius), cv.BORDER_DEFAULT)




# Can pass the blurred result for even better edges
def canny(blank, mini=125, maxi=175):
    return cv.Canny(blank, 125, 175)


# Edges get thickened
def dilate(blank, radius=7, iterations=1):
    return cv.dilate(blank, (radius, radius), iterations)


