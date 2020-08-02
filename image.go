package main

func numToDigits(number int) (digits []byte) {
	for divider := 1; true; divider *= 10 {
		digit := int((number / divider) % 10)
		if digit == 0 && divider > number {
			break
		}
		digits = append(digits, byte(digit))
	}
	return digits
}
