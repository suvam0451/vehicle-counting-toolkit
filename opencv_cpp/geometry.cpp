// OK

namespace geometry {

/* To represent a point on the image */
class Point {
   public:
    int X;
    int Y;

    Point() {}
    Point(int _X, int _Y) {
        X = _X;
        Y = _Y;
    }
};

class Line {
   public:
    int C;
    int M;

    Line() {}

    Line(Point A, Point B) {
        M = (B.Y - A.Y) / (B.X - A.X);
        C = B.Y - M * B.X;
    }

    int RelativePositionToLine(Point A);
};

int Line::RelativePositionToLine(Point A) {
    if (A.Y == A.X * M + C)
        return 0;
    if (A.Y > A.X * M + C)
        return 1;

    return -1;
}

class Segment : Line {
    Point start;
    Point end;

    Segment(Point A, Point B) {
        start = A;
        end = B;
        M = (B.Y - A.Y) / (B.X - A.X);
        C = B.Y - M * B.X;
    }
};

}  // namespace geometry