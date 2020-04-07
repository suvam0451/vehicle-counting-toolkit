# Traffic Detection Golang

Poring my traffic detection algorithm codes to golang.

## Goals

The goals of this project are to come up with that is to real-time

- Is able to run in real-time along live camera feed.
- Adapt to indian traffic sitation as much as possible.

![Segment Detection](https://i.imgur.com/Y0sq99i.png?1)
https://imgur.com/a/93jSf2K

## Approach

#### CUDA computation

- We then detect **which segment the points falls to** for a pre-specified number of y-sliced regions.

*The following image, for example shows a 16 part division.*

![Segment Detection](https://i.imgur.com/Y0sq99i.png?1)

- Next, since we already know the trajectories of vehicles, we can take the corresponding points for a vehicle group *(Applying a K-means clustering, if needed)* 
and fit the points to a line. We can then join these lines to get an approximation for expected trajectory for a lane.

We have a few strategies for how we connect line segments generated from the previous steps.

1. Taking the centroids of lines in each segment and joining their mid-points
2. Attempt to join the ends of lines in each segment by equaitable shift in angles for each line,
with a degree of relaxation.
3. Do not attempt to align the ends of lines at all. The calculations will remain true to the input setof data points.
 
### Libraries:  

#### Darknet

