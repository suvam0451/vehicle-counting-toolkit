#pragma once

#include <stdio.h>

#include <fstream>
#include <iostream>
#include <opencv2/highgui.hpp>
#include <opencv2/opencv.hpp>
#include <string>

#include "geometry.cpp"
#include "include/json/json.hpp"
#include "mouseevents.cpp"

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

void drawLine(Mat& target, float start_x, float start_y, float end_x,
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

void loadJSON(string filepath, json& target) {
    ifstream tmp(filepath);
    target << tmp;
    tmp.close();
}

// Read the config file and raw bound lines
void drawLinesFromConfig(cv::Mat& frameRef, json& configData) {
    // Draw the original detection object
    for (json& objs : configData["splits"]) {
        // Y axis is inversed
        float start_x = objs[0]["X"];
        float start_y = objs[0]["Y"];
        float end_x = objs[1]["X"];
        float end_y = objs[1]["Y"];
        drawLine(frameRef, start_x, 1.0 - start_y, end_x, 1.0 - end_y);
    }
}

int main(int argc, char** argv) {
    if (argc <= 2) {
        cout << "Expecting a image file to be passed to program" << endl;
        return -1;
    }

    if (argc <= 3) {
        cout << "Provide frame tagged/tracking tagged data as argument 4/5"
             << endl;
        return -1;
    }

    // "../data/veh_G.json"
    string __videofile = argv[1];
    string __framedata = argv[2];
    string __trackingdata = argv[3];
    string __configfile = "config.json";

    VideoCapture cap(__videofile);
    double fps = cap.get(cv::CAP_PROP_FPS);

    cout << "FPS of the video: " << fps << endl;
    bool toggleHelpText = true;
    bool togglePause = false;
    bool programActive = true;

    if (!cap.isOpened()) {
        cout << "Error opening video stream" << endl;
    }

    VehicleCountGroup grp;

    // Original frame-by-frame data
    json frame_tagged;
    loadJSON(__trackingdata, frame_tagged);

    // Reading classifier data
    json original_data;
    loadJSON(__framedata, original_data);

    // Reading classifier data
    json config_data;
    loadJSON(__configfile, config_data);

    MouseHandler __mouse;

    long long int frame_count = 0;
    long long int data_index = 0;                      // Frame skipping applied
    int frame_skip = config_data["video_frame_skip"];  // 1 -> Normal playback
    int sampling_frequency =
        config_data["sampling_frequency"];  // 1 -> Normal playback

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

                if (frame.empty()) break;

                showHUD(frame, grp);
                if (toggleHelpText) showHelpText(frame, 0);

                // frame skipping
                // if (frame_count % 3 == 0) {
                //     data_index++;
                // }
                data_index += frame_skip;

                // Draw the original detection object
                for (json& objs : original_data[data_index]["objects"]) {
                    drawRectangle(frame,
                                  objs["relative_coordinates"]["center_x"],
                                  objs["relative_coordinates"]["center_y"],
                                  objs["class_id"]);
                }

                drawLinesFromConfig(frame, config_data);

                // for (json& objs : frame_tagged[data_index]["frames"]) {
                //     drawDetectionObject(frame, objs["center_x"],
                //                         objs["center_y"], objs["id"]);
                // }

                // cout << "hello !!!" << endl;
                // Display the resulting frame
                // namedWindow("Frame", cv::WINDOW_NORMAL);
                // resizeWindow("Frame", 300, 300);
                imshow("Frame", frame);
                setMouseCallback("image", __mouse.onMouse, 0);
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
                    frame_count += frame_skip;
                }
                if (c == 39)  // Forward
                {
                    cout << "hello world, forward";
                    cap.set(1, -3);
                    data_index--;
                    frame_count -= frame_skip;
                }
                if (c == 27) {
                    programActive = false;
                    break;
                }
            }
        }
        if (!programActive) break;
    }

    cap.release();
    destroyAllWindows();
    return 0;
}