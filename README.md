# Vehicle Counting Toolkit

An analyst's tool to assist in traffic volumetric measurement using your preferable choice of object detection model.

## Installation

You can download the latest stable release from [here]()

The bundle consists of batch scripts for working with [darknet](https://github.com/AlexeyAB/darknet)

## Pre-Requisites

You have to download and install [CUDA](https://developer.nvidia.com/cuda-downloads), [cuDNN](https://developer.nvidia.com/cuDNN) and [cmake](https://cmake.org/install/). If you use an ubuntu based operating system, you can run teh following commands:

#### Ubuntu

```sh

```

#### RedHat/CentOS/Fedora

```sh
sudo dnf install https://mirrors.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm https://mirrors.rpmfusion.org/nonfree/fedora/rpmfusion-nonfree-release-$(rpm -E %fedora).noarch.rpm
# Install from F32 channel. The dependencies for F33 channel are currently broken
sudo dnf config-manager --add-repo https://developer.download.nvidia.com/compute/cuda/repos/fedora32/x86_64/cuda-fedora32.repo
sudo dnf clean all
sudo dnf module disable nvidia-driver
sudo dnf -y install cuda
# Other tools
sudo dnf -y install cmake cmake-gui
```

For running the python scripts, the recommended method is miniconda

#### Ubuntu

#### RedHat/CentOS/Fedora

```sh
sudo dnf -y install conda
conda init zsh # or conda init bash
# restart shell
conda create -n {your-env-name-choice} python=3.7
conda activate {your-env-name-choice}
conda install -c conda-forge opencv
conda install numpy scikit-learn matplotlib
```

## Examples

The following example demonstrate how to use the tools specifically for tinyYOLO and the COCO dataset.

This project is aimed towards creation of feasible real-time solutions for determining traffic flow characteristics and
statistics data generation using available lane information and general day-to-day vehicle behavior.

The goals of this project are to come up with that is to real-time

- Is able to run in real-time along live camera feed.
- Adapt to indian traffic sitation as much as possible.

## Case study

In the earliest phases, we should be able to extract information from roads such as these categorized by :-

- Clearly delineated lanes _(With all vehicles respecting the traffic)_
- Camera feed is static

![Ideal road scenario](https://i.imgur.com/gpMsysy.jpg)

## Approach

#### CUDA computation

- We then detect **which segment the points falls to** for a pre-specified number of y-sliced regions.

_The following image, for example shows a 16 part division._

![Segment Detection](https://i.imgur.com/Y0sq99i.png?1)

- Next, since we already know the trajectories of vehicles, we can take the corresponding points for a vehicle group _(Applying a K-means clustering, if needed)_
  and fit the points to a line. We can then join these lines to get an approximation for expected trajectory for a lane.

We have a few strategies for how we connect line segments generated from the previous steps.

1. Taking the centroids of lines in each segment and joining their mid-points
2. Attempt to join the ends of lines in each segment by equaitable shift in angles for each line,
   with a degree of relaxation.
3. Do not attempt to align the ends of lines at all. The calculations will remain true to the input setof data points.

### Libraries:

#### Darknet

### Structure of this repository

<details>
     <summary>About the GO project</summary>
     
- The home directory of the repository is a golang packages that can be used to run the tokenizer passes.
- The yaml file dictates the number of iterations and parameter input for each iteration.
- The yolo_mark folder has a copy of windows build of [yolo_mark](https://github.com/AlexeyAB/Yolo_mark).
You can use this to slice images from videos or tag images for genearting models.

</details>

<details>
     <summary>The CUDA library written in python using numba</summary>

- "python" folder has all the libraries along-with a sample main.py file to demonstrate all the algorithms.
- Note that you should **use miniconda/anaconda** to get your libraries so that no version mismatch errors occur.

This works with anaconda/miniconda.

```batch
conda create -n yourenvname python=x.x anaconda
conda install numba opencv matplotlib numpy scikit-learn
conda install -c anaconda cudatoolkit
```

If you prefer vanilla python with pip install, then here is the list of packages used _(I used python3.7)_.

```
pip install opencv numpy matplotlib numba
```

</details>

### Instruction for developers

<details>
     <summary>Preliminary: Bulk generation of image slices / list generation</summary>

**NOTE: All the paths mentioned here are relative to the /bin folder**

- Start with the /bin folder. Copy over images to /bin/input folder
- Run the `GenerateImages.ps1` powershell file. This will create an /intermediate folder and insert **.txt files** with lists of generated image per video file in /input folder.
  The images themselves will be outputted to /imagesets folder.
- Run the `DarknetProcess.ps1` powershell file. This will create an /output folder and start inserting **.json files** with detection data per video file in /input folder.

### With powershell

```powershell
powershell
.\GenerateImages.ps1
.\DarknetProcess.ps1
```

### With powershell Core

```powershell
pwsh
.\GenerateImages.ps1
.\DarknetProcess.ps1
```

</details>

<details>
     <summary>CUDA setup </summary>

Before running the CUDA scripts, we have to setup our CUDA environments and install required packages. This section will just list the commands but for a more step-by-step guide, please read this guide instead.

After following that guide, depending on the shell you are using, activate the conda environment and run the `matplotTag.py` script from /visualizers folder.

### Powershell

```powershell
powershell
conda activate traffic_tools
python matplotTag.py
```

### Powershell Core

```powershell
pwsh
conda activate traffic_tools
python matplotTag.py
```

### CMD

```cmd
conda activate traffic_tools
python matplotTag.py
```

</details>

<details>
     <summary>conda setup </summary>

```sh
conda -n py38 python=3.8
conda activate py38
conda install numba matplotlib scikit-learn opencv
```

</details>
<details>
     <summary>CUDA scripts: </summary>

1. Copy the `yolo_mark.exe` in /bin folder to directory with your video files. The image sampling can be done by the following command. An interval of 10 is recommended for no GPU and you can go as less as 4 if you have a GPU and videos < 5min length

```
yolo_mark.exe outpath cap_video videofile.mp4 10
```

</details>

## Libraries and licensing

This repo uses binaries produced by a lot of third-party projects. Here is the list. If you use this repository for anything, make sure to go into these repositories
