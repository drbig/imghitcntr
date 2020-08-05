package main

import (
	"testing"
)

func BenchmarkGenImage(b *testing.B) {
	setupDefaultColors()
	setupDigitsMask()

	for i := 0; i < b.N; i++ {
		genImage(i, colorBG, colorFG)
	}
}
