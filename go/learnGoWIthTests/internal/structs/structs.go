package structs

type Rectangle struct {
	Width float64
	Height float64
}

type Circle struct {
	Radius float64
}

func Perimeter(rectange Rectangle) float64 {
	return 2 * (rectange.Width + rectange.Height)
}

func Area(rectange Rectangle) float64 {
	return rectange.Width * rectange.Height
}
