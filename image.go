package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"strings"

	_ "image/gif"
)

const (
	DIGIT_W = 8
	DIGIT_H = 16
)

const DIGITS_B64 = `
R0lGODlhUAAQAKEBAAAAAP///////////yH5BAEKAAIALAAAAABQABAAAAKmhI+py+0PUZghTUtl
xhXQbYDfFYKe+WFcV5JaN7InGx+pLcK0e5Z+jyPxcA2h7mf8zTTIY9CpYiSB08vmphgCe8aco+pM
Xk3MKFfnLYbXLZSM+IJru3MZdkbnabct+Zt6p4T30sVXJ4UWCFbWlIVm+LcHOBI1JWjpd/k2hNkI
qFnGuQPK6BnKBnlqNlha07bgcqe46UYrShhLpiRLJqonGREsPExcAAA7
`

var (
	colorBG   color.RGBA
	colorFG   color.RGBA
	digitsImg *image.Alpha
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

func setupDigitsImg() {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(DIGITS_B64))
	src, _, err := image.Decode(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to parse built-in digits\n")
		os.Exit(2)
	}
	// A paranoid person could add bounds check here if it's 10*DIGIT_W x DIGIT_H

	bounds := src.Bounds() //you have defined that both src and mask are same size, and maskImg is a grayscale of the src image. So we'll use that common size.
	digitsImg = image.NewAlpha(bounds)
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			r, _, _, _ := src.At(x, y).RGBA()
			digitsImg.SetAlpha(x, y, color.Alpha{uint8(r)})
		}
	}
}

func genImage(hits int, bg, fg color.RGBA) (img draw.Image) {
	digits := numToDigits(hits)
	width := len(digits) * DIGIT_W
	bgi := image.NewUniform(bg)
	fgi := image.NewUniform(fg)

	img = image.NewRGBA(image.Rect(0, 0, width, DIGIT_H))
	draw.Draw(img, image.Rect(0, 0, width, DIGIT_H), bgi, image.Point{}, draw.Over)
	for i, digit := range digits {
		maskPoint := image.Point{int(digit) * DIGIT_W, 0}
		rect := image.Rect(i*DIGIT_W, 0, (i+1)*DIGIT_W, DIGIT_H)
		draw.DrawMask(img, rect, fgi, image.Point{}, digitsImg, maskPoint, draw.Over)
	}

	return img
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
