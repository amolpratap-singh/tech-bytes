package array

import (
	"reflect"
	"testing"
)


func TestSum(t *testing.T) {
	t.Run("Collection of 5 integers", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}
		got := Sum(numbers)
		want := 15

		if got != want {
			t.Errorf("got %d want %d given %v", got, want, numbers)
		}
	})

	t.Run("collection of any size", func(t *testing.T) {
		numbers := []int{1, 3, 5}
		got := Sum(numbers)
		want := 9

		if got != want {
			t.Errorf("got %d want %d given %v", got, want, numbers)
		}
	})
}

func TestSumAll(t *testing.T) {

	t.Run("Sum of two given arrays", func(t *testing.T) {
		arr1 := []int{1, 2}
		arr2 := []int{0, 9}
		got := SumAll(arr1, arr2)
		want := []int{3, 9}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %d want %d given arr1: %v and arr2: %v", got, want, arr1, arr2)
		}
	})

	t.Run("Sum of one slice passed", func(t *testing.T) {
		got := SumAll([]int{1, 1, 1})
		want := []int{3}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got %v want %v for the given slice []int{1,1,1}", got, want)
		}
	})
}