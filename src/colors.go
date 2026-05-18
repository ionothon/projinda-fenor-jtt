package main

import (
	"math"
)

func rgbToHsv(r, g, b uint32) (float64, float64, float64) {
	red := float64(r>>8) / 255.0
	green := float64(g>>8) / 255.0
	blue := float64(b>>8) / 255.0

	max := math.Max(red, math.Max(green, blue))
	min := math.Min(red, math.Min(green, blue))
	delta := max - min

	var hue, sat, val float64

	if delta == 0 {
		hue = 0
	} else if max == red {
		hue = 60 * math.Mod((green-blue)/delta, 6)
	} else if max == green {
		hue = 60 * (((blue - red) / delta) + 2)
	} else if max == blue {
		hue = 60 * (((red - green) / delta) + 4)
	}

	if hue < 0 {
		hue += 360
	}

	if max == 0 {
		sat = 0
	} else {
		sat = delta / max
	}

	val = max

	return hue, sat, val
}

func classifyColor(h, s, v float64) string {
	if v <= 0.15 {
		return "Black"
	}
	if s < 0.15 {
		if v > 0.85 {
			return "White"
		}
		return "Grey"
	}

	if h >= 345 || h < 20 {
		return "Red"
	} else if h >= 20 && h < 45 {
		return "Orange"
	} else if h >= 45 && h < 70 {
		return "Yellow"
	} else if h >= 70 && h < 145 {
		return "Green"
	} else if h >= 145 && h < 195 {
		return "Cyan"
	} else if h >= 195 && h < 255 {
		return "Blue"
	} else if h >= 255 && h < 305 {
		return "Purple"
	} else if h >= 305 && h < 345 {
		return "Pink"
	}
	return "Unknown"
}
