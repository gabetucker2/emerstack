// Range returns the min, max (and associated indexes, -1 = no values) for the tensor.
	// This is needed for display and is thus in the core api in optimized form
	// Other math operations can be done using gonum/floats package.
	Range() (min, max float64, minIdx, maxIdx int)
	
	
// Clamp clamps x to the provided closed interval [a, b]
func Clamp(x, a, b float32) float32 {
	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}

// ClampInt clamps x to the provided closed interval [a, b]
func ClampInt(x, a, b int) int {
	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}

see popcode1d.go for example

		val = mat32.Clamp(val, pc.Min, pc.Max)



in srn.go 


ss.ConfigWorldTsrs()  // Needs to be added to Config program

func(ss *Sim) ConfigWorldTsrs() {
	if ss.EnvpTsr == nil {
		ss.EnvpTsr = etensor.NewFloat32([]int{1,7},nil,nil)
	}
	if ss.EnvcTsr == nil {
		ss.EnvcTsr = etensor.NewFloat32([]int{1,7},nil,nil)
	}
	if ss.IntpTsr == nil {
		ss.IntpTsr = etensor.NewFloat32([]int{1,7},nil,nil)	
	}
	if ss.IntcTsr == nil {
		ss.IntcTsr = etensor.NewFloat32([]int{1,7},nil,nil)	
	}
	if ss.EnviroTsr == nil {
		ss.EnviroTsr = etensor.NewFloat32([]int{1,7},nil,nil)	
	}	
	if ss.InteroTsr == nil {
		ss.InteroTsr = etensor.NewFloat32([]int{1,7},nil,nil)	
}
}
