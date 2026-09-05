package structs

import (
	"fmt"
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

	checkArea := func(t testing.TB, shape Shape, want float64) {
		t.Helper()
		got := shape.Area()
		fmt.Printf("shape got as %v \n", shape)
		if got != want {
			t.Errorf("Area Calculated for %v which got: %0.2f and want: %g", shape, got, want)
		}
	}
	t.Run("Test case to calculate the area of object", func(t *testing.T) {
		rectangle := Rectangle{12.0, 6.0}
		//got := rectangle.Area()
		want := 72.0

		checkArea(t, rectangle, want)
	})

	t.Run("Test case to calcuate the area of circle", func(t *testing.T) {
		circle := Circle{10.0}
		//got := circle.Area()
		want := 314.1592653589793

		checkArea(t, circle, want)
	})

	// Table Driven Test
	areaTests := []struct {
		name    string
		shape   Shape
		hasArea float64
	}{
		{name: "Rectangle", shape: Rectangle{Width: 12, Height: 6}, hasArea: 72.0},
		{name: "Circle", shape: Circle{Radius: 10}, hasArea: 314.1592653589793},
		{name: "Triangle", shape: Triangle{Base: 12, Height: 6}, hasArea: 36.0},
	}

	for _, tt := range areaTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			got := tt.shape.Area()
			if got != tt.hasArea {
				t.Errorf("%#v got %g want %g", tt.shape, got, tt.hasArea)
			}
		})

	}

}
