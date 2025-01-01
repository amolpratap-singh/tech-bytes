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

func TestSumAllSlice(t *testing.T) {
	t.Run("Sum of two slice for the new refactor function", func(t *testing.T) {
		arr1 := []int{1, 5}
		arr2 := []int{0, 11}
		got := SumAllSlice(arr1, arr2)
		want := []int{6, 11}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got %v want %v for the given slice %v and %v", got, want, arr1, arr2)
		}

	})
}

func TestSumAllTails(t *testing.T) {

	checkSums := func(t testing.TB, got, want []int) {
		t.Helper()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Test case failed due to requored: %v but got: %v", want, got)
		}
	}

	t.Run("Test case to validate single arr tail element sum", func(t *testing.T) {
		got := SumAllTails([]int{1, 2, 3, 4, 5})
		want := []int{14}
		checkSums(t, got, want)
	})

	t.Run("Test case to validate more than one arr tail element sum", func(t *testing.T) {
		got := SumAllTails([]int{1, 3, 5, 7}, []int{2, 4, 6, 8}, []int{0, 1, 2, 3})
		want := []int{15, 18, 6}
		checkSums(t, got, want)
	})

	t.Run("Test case to validate for the empty element in the arr", func(t *testing.T) {
		got := SumAllTails([]int{}, []int{1, 2, 3})
		want:= []int{0, 5}
		checkSums(t, got, want)
	})
}
