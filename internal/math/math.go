package math

import "sort"

func Divisors(num int64) []int64 {
	var divisors []int64
	for i := int64(1); i*i <= num; i++ {
		if num%i == 0 {
			divisors = append(divisors, i)
			if i != num/i {
				divisors = append(divisors, num/i)
			}
		}
	}
	sort.Slice(divisors, func(i, j int) bool {
		return divisors[i] < divisors[j]
	})
	return divisors
}
