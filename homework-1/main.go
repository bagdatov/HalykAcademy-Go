package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// SortSlice is an implemention of bubble sort method,
// time complexity is O(n^2).
func SortSlice(slice []int) {
	for i := 0; i < len(slice)-1; i++ {
		for j := 0; j < len(slice)-i-1; j++ {

			if slice[j] > slice[j+1] {
				slice[j], slice[j+1] = slice[j+1], slice[j]
			}
		}
	}
}

// IncrementOdd is adding one to odd indexes of the slice.
func IncrementOdd(slice []int) {
	for i := range slice {
		if i%2 != 0 {
			slice[i]++
		}
	}
}

// PrintSlice is using fmt package to print slice without brackets.
func PrintSlice(slice []int) {
	fmt.Println(slice)
}

// ReverseSlice is reversing order of values in the slice.
func ReverseSlice(slice []int) {
	middle := len(slice) / 2
	end := len(slice) - 1

	for i := range slice {
		if i >= middle {
			break
		}
		slice[i], slice[end-i] = slice[end-i], slice[i]
	}
}

// appendFunc is merging functions that are passed as arguments.
func appendFunc(dst func([]int), src ...func([]int)) func([]int) {
	return func(slice []int) {
		dst(slice)
		for _, f := range src {
			f(slice)
		}
	}
}

func main() {
	// slice := []int{-1, 2, 3, 4, 5, 0, 7, 8}
	// slice := []int{-1, -2}
	// slice := []int{0}
	// slice := []int{}
	// slice := []int{-9223372036854775808, 9223372036854775807}

	slice := make([]int, 100)

	// seeding time for pseudo-random numbers
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		slice[i] = rand.Intn(math.MaxInt16)
	}

	fmt.Println("PrintSlice function:")
	PrintSlice(slice)
	fmt.Println()

	fmt.Println("SortSlice function:")
	SortSlice(slice)
	PrintSlice(slice)
	fmt.Println()

	fmt.Println("IncrementOdd function:")
	IncrementOdd(slice)
	PrintSlice(slice)
	fmt.Println()

	fmt.Println("ReverseSlice function:")
	ReverseSlice(slice)
	PrintSlice(slice)
	fmt.Println()

	superFunc := appendFunc(IncrementOdd, ReverseSlice, PrintSlice)
	fmt.Println("appendFunc function:")
	superFunc(slice)

}
