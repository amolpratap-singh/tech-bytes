package structs

import (
	"testing"
)

func TestPerimeter(t *testing.T) {
	t.Run("Test case to calculate the perimeter of object", func(t *testing.T) {
		rectangle := Rectangle{10.0, 10.0}
		got := Perimeter(rectangle)
		want := 40.0
	
		if got != want {
			t.Errorf("Perimeter calculation failed as got: %0.2f and want: %0.2f", got, want)
		}

	})
}

func TestArea(t *testing.T) {
	t.Run("Test case to calculate the area of object", func(t *testing.T) {
		rectangle := Rectangle{12.0, 6.0}
		got := Area(rectangle)
		want := 72.0

		if got != want {
			t.Errorf("Area calculation failed as got : %0.2f and want %0.2f", got, want)
		}
	})

	t.Run("Test case to calcuate the area of circle", func(t *testing.T) {
		circle := Circle{10.0}
		got := Area(circle)
		want := 314.1592653589793

		if got != want {
			t.Errorf("got %g want %g", got, want)
		}
	})
}