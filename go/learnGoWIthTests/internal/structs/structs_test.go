package structs

import (
	"testing"
)

func TestPerimeter(t *testing.T) {
	t.Run("Test case to calculate the perimeter of object", func(t *testing.T) {
		got := Perimeter(10.0, 10.0)
		want := 40.0
	
		if got != want {
			t.Errorf("Perimeter calculation failed as got: %0.2f and want: %0.2f", got, want)
		}

	})
}

func TestArea(t *testing.T) {
	t.Run("Test case to calculate the area of object", func(t *testing.T) {
		got := Area(12.0, 6.0)
		want := 72.0

		if got != want {
			t.Errorf("Area calculation failed as got : %0.2f and want %0.2f", got, want)
		}
	})
}