package main

func clampUI(x float32) float32 {
	return clamp(x, 0, 1)
}

func clamp(x, min, max float32) (xprime float32) {
	xprime = x
	if x > max {
		xprime = max
	} else if x < min {
		xprime = min
	}
	return
}
