#pragma once

#include <stdio.h>
#include <fstream>
#include <iostream>
#include <opencv2/highgui.hpp>
#include <opencv2/opencv.hpp>
#include <string>
#include "geometry.cpp"

#include "include/json/json.hpp"

using namespace cv;
using namespace std;

using json = nlohmann::json;

// Font height for plotting text for opencv
#define font_height 24
#define FONT_WEIGHT 1.75

struct VehicleCountGroup {
    int motorcycleCount;
    int carcount;
    int truckCount;
    int busCount;

    VehicleCountGroup()
        : motorcycleCount(0), carcount(0), truckCount(0), busCount(0){};
};

enum OperationModes { CarCounting, LanePrediction, CarCountingSugmentWise };

struct TextElement {
    string text;
    Scalar color;
    double fontScale;

    TextElement() : text(""), color(118, 185, 0), fontScale(0.75){};
    TextElement(string str) : text(str), color(118, 185, 0), fontScale(0.75){};
    TextElement(string str, double val)
        : text(str), color(118, 185, 0), fontScale(val){};
};

// function to display HUD and stats
void showHUD(Mat& target, VehicleCountGroup& values) {
    string motorcycleText =
        "Number of motorcycles: " + to_string(values.motorcycleCount);
    string carText = "Number of Cars: " + to_string(values.carcount);
    string truckText = "Number of Trucks: " + to_string(values.truckCount);
    string busText = "Number of Buses: " + to_string(values.busCount);

    int stringIndex = 0;

    putText(target, motorcycleText,
            Point(10, target.rows / 2 + font_height * stringIndex++),
            cv::FONT_HERSHEY_DUPLEX, 0.75, CV_RGB(118, 185, 0), FONT_WEIGHT);
    putText(target, carText,
            Point(10, target.rows / 2 + font_height * stringIndex++),
            cv::FONT_HERSHEY_DUPLEX, 0.75, CV_RGB(118, 185, 0), FONT_WEIGHT);
    putText(target, truckText,
            Point(10, target.rows / 2 + font_height * stringIndex++),
            cv::FONT_HERSHEY_DUPLEX, 0.75, CV_RGB(118, 185, 0), FONT_WEIGHT);
    putText(target, busText,
            Point(10, target.rows / 2 + font_height * stringIndex++),
            cv::FONT_HERSHEY_DUPLEX, 0.75, CV_RGB(118, 185, 0), FONT_WEIGHT);

    stringIndex++;
    stringIndex++;

    putText(target, "[ ] : Cars",
            Point(10, target.rows / 2 + font_height * stringIndex++),
            cv::FONT_HERSHEY_DUPLEX, 0.75, CV_RGB(255, 0, 0), FONT_WEIGHT);
    putText(target, "[ ] : Motorbike",
            Point(10, target.rows / 2 + font_height * stringIndex++),
            cv::FONT_HERSHEY_DUPLEX, 0.75, CV_RGB(0, 255, 0), FONT_WEIGHT);
    putText(target, "[ ] : Trucks",
            Point(10, target.rows / 2 + font_height * stringIndex++),
            cv::FONT_HERSHEY_DUPLEX, 0.75, CV_RGB(0, 0, 255), FONT_WEIGHT);
    putText(target, "[ ] : Bicycle",
            Point(10, target.rows / 2 + font_height * stringIndex++),
            cv::FONT_HERSHEY_DUPLEX, 0.75, CV_RGB(0, 255, 255), FONT_WEIGHT);
}

void showHelpText(Mat& target, int values) {
    // 60 pixel offset from bottom
    Point basePosition(target.cols - 320, target.rows - 60);
    int stringIndex = 0;
    double fontScale = 0.64;

    vector<TextElement> msgs = {TextElement("Press [P] to pause", fontScale),
                                TextElement("Press [ESC] to exit", fontScale)};

    for (auto it : msgs) {
        putText(target, it.text,
                basePosition + Point(0, font_height * stringIndex++),
                cv::FONT_HERSHEY_DUPLEX, fontScale, it.color, FONT_WEIGHT);
    }
}

void drawLine(Mat& target,
              float start_x,
              float start_y,
              float end_x,
              float end_y) {
    Point pt1(target.cols * start_x, target.rows * start_y);
    Point pt2(target.cols * end_x, target.rows * end_y);
    line(target, pt1, pt2, Scalar(100), 4);
}

void drawDetectionObject(Mat& target, float X, float Y, int id) {
    Point pt1(target.cols * X - 4.0f, target.rows * Y - 4.0f);
    Point pt2(target.cols * X + 4.0f, target.rows * Y + 4.0f);
    Point pt3(target.cols * X, target.rows * Y);

    rectangle(target, pt1, pt2, CV_RGB(50, 50, 50), 4);
    putText(target, to_string(id), pt3, cv::FONT_HERSHEY_DUPLEX, 1.0,
            CV_RGB(100, 100, 100), FONT_WEIGHT + 2);
}

void drawRectangle(Mat& target, float X, float Y, int id) {
    Scalar color;
    switch (id) {
        case 2:  // car (red)
            color = Scalar(255, 0, 0);
            break;
        case 3:  // motorbike (Green)
            color = Scalar(0, 255, 0);
            break;
        case 7:  // Truck (Blue)
            color = Scalar(0, 0, 255);
            break;
        case 1:  // Bicycle (Yellow)
            color = Scalar(0, 255, 255);
            break;
        default:
            break;
    }
    Point pt1(target.cols * X - 8.0f, target.rows * Y - 8.0f);
    Point pt2(target.cols * X + 8.0f, target.rows * Y + 8.0f);

    rectangle(target, pt1, pt2, color, 4);
}

int testDetectionValidity() {
    Mat image = imread("../image/Gangasagar_00000000.jpg");

    // Reading classifier data
    ifstream my_file_frame_tagged("../data/G_2_0_02.json");
    json j;
    j << my_file_frame_tagged;

    for (json& objs : j[0]["objects"]) {
        drawRectangle(image, objs["center_x"], objs["center_y"],
                      objs["class_id"]);
    }

    imshow("Frame", image);
    waitKey();
    return 0;
}

int main(int argc, char** argv) {
    // testDetectionValidity();
    // return 0;

    if (argc != 2) {
        cout << "Expecting a image file to be passed to program" << endl;
        return -1;
    }

    VideoCapture cap(argv[1]);
    double fps = cap.get(cv::CAP_PROP_FPS);

    cout << "FPS of the video: " << fps << endl;
    bool toggleHelpText = true;
    bool togglePause = false;
    bool programActive = true;

    if (!cap.isOpened()) {
        cout << "Error opening video stream" << endl;
    }

    VehicleCountGroup grp;

    // Reading classifier data
    ifstream my_file_frame_tagged("../data/veh_G.json");
    json frame_tagged;
    frame_tagged << my_file_frame_tagged;
    my_file_frame_tagged.close();

    // Reading classifier data
    ifstream my_file("../data/G_2_0_02.json");
    json original_data;
    original_data << my_file;
    my_file.close();

    long long int frame_count = 0;
    long long int data_index = 0;  // Frame skipping applied

    while (programActive) {
        if (togglePause) {
            // ESC to exit
            char c = (char)waitKey(25);
            if (c == 104) {
                toggleHelpText = !toggleHelpText;
            }
            if (c == 112) {
                togglePause = !togglePause;
            }
            if (c == 27) {
                programActive = false;
                break;
            }
        } else {
            while (!togglePause) {
                Mat frame;
                cap >> frame;

                if (frame.empty())
                    break;

                showHUD(frame, grp);
                if (toggleHelpText)
                    showHelpText(frame, 0);

                // frame skipping
                if (frame_count % 3 == 0) {
                    data_index++;
                }

                for (json& objs : original_data[data_index]["objects"]) {
                    drawRectangle(frame, objs["center_x"], objs["center_y"],
                                  objs["class_id"]);
                }

                for (json& objs : frame_tagged[data_index]["frames"]) {
                    drawDetectionObject(frame, objs["center_x"],
                                        objs["center_y"], objs["id"]);
                }

                drawLine(frame, 0.1, 0.8, 0.85, 0.55);
                drawLine(frame, 0.15, 0.35, 0.3, 0.28);
                drawLine(frame, 0.45, 0.25, 0.70, 0.30);

                // Display the resulting frame
                // namedWindow("Frame", cv::WINDOW_NORMAL);
                // resizeWindow("Frame", 300, 300);
                imshow("Frame", frame);
                frame_count++;

                // ESC to exit
                char c = (char)waitKey(25);
                if (c == 104) {
                    toggleHelpText = !toggleHelpText;
                }
                if (c == 112) {
                    togglePause = !togglePause;
                }
                if (c == 37)  // Back
                {
                    cout << "hello world, back";
                    cap.set(1, 3);
                    data_index++;
                    frame_count += 3;
                }
                if (c == 39)  // Forward
                {
                    cout << "hello world, forward";
                    cap.set(1, -3);
                    data_index--;
                    frame_count -= 3;
                }
                if (c == 27) {
                    programActive = false;
                    break;
                }
            }
        }
        if (!programActive)
            break;
    }

    cap.release();
    destroyAllWindows();
    return 0;
}