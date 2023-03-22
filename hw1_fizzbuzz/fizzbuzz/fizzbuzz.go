package fizzbuzz

import "strconv"

func FizzBuzz(i int) string {
	switch rem := i % 15; {
	case rem == 0:
		return "FizzBuzz"
	case rem%5 == 0:
		return "Buzz"
	case rem%3 == 0:
		return "Fizz"
	default:
		return strconv.Itoa(i)
	}
}
