

namespace geometry {

class Point {
    int X;
    int Y;

    Point(int _X, int _Y) {
        X = _X;
        Y = _Y;
    }
};

class Line {
   public:
    int X;
    int C;
    int M;
};

static int PositionRelativeToLine() {
    return 0;
}
}  // namespace geometry