#pragma once

#include <stdio.h>

#include <fstream>
#include <iostream>
#include <opencv2/highgui.hpp>
#include <opencv2/opencv.hpp>
#include <string>

#include "geometry.cpp"

using namespace cv;
using namespace std;

using json = nlohmann::json;

class MouseHandler {
   private:
    bool isDragging = false;

   public:
    static void onMouse(int event, int x, int y, int, void*) {
        cout << "Mouse click event..." << endl;
        if (event != EVENT_LBUTTONDOWN) {
            cout << "Left mouse click..." << endl;
            return;
        }

        // Point seed = Point(x, y);
        // int lo = ffillMode == 0 ? 0 : loDiff;
        // int up = ffillMode == 0 ? 0 : upDiff;
        // int flags = connectivity + (newMaskVal << 8) +
        //             (ffillMode == 1 ? FLOODFILL_FIXED_RANGE : 0);
        // int b = (unsigned)theRNG() & 255;
        // int g = (unsigned)theRNG() & 255;
        // int r = (unsigned)theRNG() & 255;
        // Rect ccomp;

        // Scalar newVal = isColor ? Scalar(b, g, r)
        //                         : Scalar(r * 0.299 + g * 0.587 + b * 0.114);
        // Mat dst = isColor ? image : gray;
        // int area;

        // if (useMask) {
        //     threshold(mask, mask, 1, 128, THRESH_BINARY);
        //     area = floodFill(dst, mask, seed, newVal, &ccomp,
        //                      Scalar(lo, lo, lo), Scalar(up, up, up), flags);
        //     imshow("mask", mask);
        // } else {
        //     area = floodFill(dst, seed, newVal, &ccomp, Scalar(lo, lo, lo),
        //                      Scalar(up, up, up), flags);
        // }

        // imshow("image", dst);
        // cout << area << " pixels were repainted\n";
    }
};