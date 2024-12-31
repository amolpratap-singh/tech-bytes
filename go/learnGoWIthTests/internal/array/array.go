package array

import "fmt"


func Sum(numbers []int) int {
	sumNumber := 0

	for _, number := range numbers {
		sumNumber += number
	}

	return sumNumber
}

func SumAll(numbersToSum ...[]int) []int {
	lengthOfNumbers := len(numbersToSum)
	sum := make([]int, lengthOfNumbers)
	fmt.Println("Input got and constructed %v ", sum)
	for i, numbers := range numbersToSum {
		sum[i] = Sum(numbers)
	}
	return sum

}