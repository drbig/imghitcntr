package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"
)

var (
	colorBG color.RGBA
	colorFG color.RGBA
)

func setupDefaultColors() {
	var ok bool

	colorBG, ok = parseHexColor(flagBgColor)
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: can't parse BG color: \"%s\"\n", flagBgColor)
		flag.Usage()
		os.Exit(1)
	}

	colorFG, ok = parseHexColor(flagFgColor)
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: can't parse FG color: \"%s\"\n", flagFgColor)
		flag.Usage()
		os.Exit(1)
	}
}

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

// This is after https://stackoverflow.com/a/54200713
func parseHexColor(s string) (c color.RGBA, ok bool) {
	ok = true
	c.A = 0xff

	if s[0] != '#' {
		return c, false
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		ok = false
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		ok = false
	}
	return
}
