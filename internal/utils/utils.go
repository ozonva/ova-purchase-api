package utils

import (
	"errors"
	"math"
)

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func Split(input []int, size int) ([][]int, error) {
	if size <= 0 {
		return nil, errors.New("size must be positive")
	}
	batches := int(math.Ceil(float64(len(input)) / float64(size)))
	result := make([][]int, 0, batches)
	for i := 0; i < batches; i++ {
		left := i * size
		right := min(len(input), left+size)
		source := input[left:right]
		buf := make([]int, len(source))
		copy(buf, source)
		result = append(result, buf)
	}
	return result, nil
}

func Reverse(input map[string]string) (map[string]string, error) {
	result := make(map[string]string)
	for k, v := range input {
		if _, ok := result[v]; ok {
			return nil, errors.New("key already exists")
		}
		result[v] = k
	}
	return result, nil
}

func Filter(input []int, exclude []int) []int {
	if len(exclude) == 0 {
		return input
	}
	set := make(map[int]bool)
	for _, value := range exclude {
		set[value] = true
	}
	result := make([]int, 0)
	for _, value := range input {
		if _, ok := set[value]; !ok {
			result = append(result, value)
		}
	}
	return result
}
