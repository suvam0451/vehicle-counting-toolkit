
#include "geometry.h"

int Line::RelativePositionToLine(Point A) {
    if (A.Y == A.X * M + C) return 0;
    if (A.Y > A.X * M + C) return 1;

    return -1;
}