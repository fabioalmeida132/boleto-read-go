package Utils

import (
	"regexp"
	"strconv"
)

func Number(code string) string {
	re := regexp.MustCompile(`[^\d]`)
	return re.ReplaceAllString(code, "")
}

func GetInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func Mod11(rawNumber string) string {
	num := Number(rawNumber)
	sum := 0
	weight := 2
	base := 9
	counter := len(num) - 1

	for index := counter; index >= 0; index-- {
		sum += GetInt(num[index:index+1]) * weight
		if weight < base {
			weight++
		} else {
			weight = 2
		}
	}

	digit := 11 - (sum % 11)
	if digit > 9 {
		digit = 0
	}
	if digit == 0 {
		digit = 1
	}

	return strconv.Itoa(digit)
}

func Mod10(rawNumber string) string {
	num := Number(rawNumber)
	sum := 0
	weight := 2
	counter := len(num) - 1

	for counter >= 0 {
		product := GetInt(num[counter:counter+1]) * weight
		if product >= 10 {
			product = 1 + (product - 10)
		}

		sum += product
		if weight == 2 {
			weight = 1
		} else {
			weight = 2
		}

		counter -= 1
	}

	digit := 10 - (sum % 10)
	if digit == 10 {
		digit = 0
	}

	return strconv.Itoa(digit)
}
