package model

import (
	"math"
	"sort"
)

func GetAllowedColors(positionX, positionY int) []*Color {

	var allowedColors = []*Color{

		// WHITE
		{
			Id: 0,
			R:  65535,
			G:  65535,
			B:  65535,
		},

		// LIGHT_GRAY
		{
			Id: 3,
			R:  47545,
			G:  47545,
			B:  47545,
		},

		// RED
		{
			Id: 4,
			R:  60395,
			G:  4626,
			B:  5911,
		},

		// ORANGE
		{
			Id: 5,
			R:  65278,
			G:  25957,
			B:  7453,
		},

		// YELLOW
		{
			Id: 6,
			R:  65535,
			G:  57311,
			B:  13364,
		},

		// GREEN
		{
			Id: 7,
			R:  9509,
			G:  50372,
			B:  10794,
		},

		// MINT
		{
			Id: 8,
			R:  11051,
			G:  65278,
			B:  36494,
		},

		// SKYBLUE
		{
			Id: 9,
			R:  0,
			G:  43690,
			B:  64507,
		},

		// SEABLUE
		{
			Id: 10,
			R:  2570,
			G:  9252,
			B:  51143,
		},

		// PURPLE
		{
			Id: 11,
			R:  38293,
			G:  4626,
			B:  44461,
		},

		// PINK
		{
			Id: 12,
			R:  55512,
			G:  24415,
			B:  39835,
		},

		// BIEGE
		{
			Id: 13,
			R:  65278,
			G:  41634,
			B:  34438,
		},

		// BROWN
		{
			Id: 14,
			R:  38293,
			G:  18504,
			B:  11051,
		},

		// BLACK
		{
			Id: 1,
			R:  0,
			G:  0,
			B:  0,
		},

		// DARK_GRAY
		{
			Id: 2,
			R:  17990,
			G:  17990,
			B:  17990,
		},

		// DARK_RED
		{
			Id: 15,
			R:  26728,
			G:  3341,
			B:  3341,
		},

		// DARK_ORANGE
		{
			Id: 16,
			R:  47545,
			G:  12336,
			B:  4626,
		},

		// DARK_YELLOW
		{
			Id: 17,
			R:  58853,
			G:  38807,
			B:  9509,
		},

		// DARK_GREEN
		{
			Id: 18,
			R:  2056,
			G:  15677,
			B:  6939,
		},

		// DARK_MINT
		{
			Id: 19,
			R:  4112,
			G:  28013,
			B:  21588,
		},

		// DARK_SKYBLUE
		{
			Id: 20,
			R:  0,
			G:  19789,
			B:  37008,
		},

		// DARK_SEABLUE
		{
			Id: 21,
			R:  1542,
			G:  3855,
			B:  22616,
		},

		// DARK_PURPLE
		{
			Id: 22,
			R:  18761,
			G:  1799,
			B:  23644,
		},

		// DARK_PINK
		{
			Id: 23,
			R:  31611,
			G:  12079,
			B:  18761,
		},

		// DARK_BIEGE
		{
			Id: 24,
			R:  50372,
			G:  27499,
			B:  18504,
		},

		// DARK_BROWN
		{
			Id: 25,
			R:  22616,
			G:  10794,
			B:  4883,
		},
	}

	for i, col := range allowedColors {
		if i < len(allowedColors)/2 {
			col.X = positionX + i*24
			col.Y = positionY
		} else {
			col.X = positionX + (i-len(allowedColors)/2)*24
			col.Y = positionY + 24
		}
	}

	sort.Slice(allowedColors, func(i, j int) bool {
		return allowedColors[i].Id < allowedColors[j].Id
	})

	return allowedColors

}

type Color struct {
	Id int
	R  uint32
	G  uint32
	B  uint32
	X  int
	Y  int
}

func (c *Color) Distance(other *Color) float64 {

	var r,
		g, b uint32

	if c.R > other.R {
		r = c.R - other.R
	} else {
		r = other.R - c.R
	}

	if c.G > other.G {
		g = c.G - other.G
	} else {
		g = other.G - c.G
	}

	if c.B > other.B {
		b = c.B - other.B
	} else {
		b = other.B - c.B
	}

	res := math.Pow(float64(r), 2) + math.Pow(float64(g), 2) + math.Pow(float64(b), 2)

	return math.Sqrt(res)
}

func (c *Color) FindClosest(allowedColors []*Color) *Color {
	var minDistance float64 = math.MaxFloat64
	var minColor *Color = nil
	for _, color := range allowedColors {
		dist := c.Distance(color)
		if dist < minDistance {
			minDistance = dist
			minColor = color
		}
	}
	return minColor
}
