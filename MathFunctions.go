package main

import (
	"math"
)

// * default TimeIncrement function
func TimeIncrement(x_ui, dx_ui, tprev_s, tcur_s, dt_s float32) (xprime_ui float32) {
	
	// compute xprime
	Δt_s  := tcur_s - tprev_s
	Δx_ui := Δt_s * (dx_ui / dt_s)
	xprime_ui = clampUI(x_ui + Δx_ui)
	
	// return
	return
	
}

// * default EnviroUnityToNN function 
func EnviroUnityToNN(d_unity, a float32) (d_NN float32) {
	
	// compute d_NN
	d_NN = 1 - float32(math.Pow(math.E, float64(-1 * d_unity * a)))
	
	// return
	return
	
}
