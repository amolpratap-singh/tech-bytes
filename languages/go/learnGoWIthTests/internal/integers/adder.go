package integers

import "fmt"

// Add takes two integers and returns the sum of them.
func Add(x, y int) int {
	return x + y
}

func main() {
	fmt.Printf("Code for additon of two  numbers %v", Add(1, 3))
}
