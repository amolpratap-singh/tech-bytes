package array

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
	for i, numbers := range numbersToSum {
		sum[i] = Sum(numbers)
	}
	return sum
}

func SumAllSlice(numbersToSum ...[]int) []int {
	var sums []int

	for _, number := range numbersToSum {
		sums = append(sums, Sum(number))
	}
	return sums
}

func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int 

	for _, number := range numbersToSum {

		if len(number) == 0 {
			sums = append(sums, 0)
		} else {
			sums = append(sums, Sum(number[1:]))
		}
	}

	return sums
}