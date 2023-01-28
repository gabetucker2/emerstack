// * set up dependencies
package main

import (
	. "github.com/gabetucker2/gostack"

)

// * initialize variables
var Parameters *Stack
var ComplexActions *Stack
var currentParameterIdx *int
var tprev_s, tcur_s, dt_s int

// * define enums
type MinOrMax int
const (
	Min MinOrMax = iota
	Max
)

// * define structs
type Relation struct {
	ThisLayer string
	OtherParameter string
	OtherLayer string
	Dx float32
}

type Requirement struct {
	Parameter string // if empty, then we will assume that the current layer is the parameter
	Layer string
	MinOrMax MinOrMax
	Threshold float32
}

type Update struct {
	Parameter string // if empty, then we will assume that the current layer is the parameter
	Layer string
	Dx float32
}

type Action struct {
	Requirements []*Requirement
	Updates []*Update
	Cost float32
}

// * define struct initializer functions
func MakeRelation(thisLayer, otherParameter, otherLayer string, dx float32) *Relation {
	newRelation := new(Relation)
	newRelation.ThisLayer = thisLayer
	newRelation.OtherParameter = otherParameter
	newRelation.OtherLayer = otherLayer
	newRelation.Dx = dx
	return newRelation
}

func MakeSimpleRequirement(layer string, minOrMax MinOrMax, threshold float32) *Requirement {
	newRequirement := new(Requirement)
	newRequirement.Parameter = ""
	newRequirement.Layer = layer
	newRequirement.MinOrMax = minOrMax
	newRequirement.Threshold = threshold
	return newRequirement
}

func MakeRequirement(parameter, layer string, minOrMax MinOrMax, threshold float32) *Requirement {
	newRequirement := new(Requirement)
	newRequirement.Parameter = parameter
	newRequirement.Layer = layer
	newRequirement.MinOrMax = minOrMax
	newRequirement.Threshold = threshold
	return newRequirement
}

func MakeUpdate(parameter, layer string, dx float32) *Update {
	newUpdate := new(Update)
	newUpdate.Parameter = parameter
	newUpdate.Layer = layer
	newUpdate.Dx = dx
	return newUpdate
}

func MakeSimpleUpdate(layer string, dx float32) *Update {
	newUpdate := new(Update)
	newUpdate.Parameter = ""
	newUpdate.Layer = layer
	newUpdate.Dx = dx
	return newUpdate
}

func MakeAction(requirements []*Requirement, updates []*Update, cost float32) *Action {
	newAction := new(Action)
	newAction.Requirements = requirements
	newAction.Updates = updates
	newAction.Cost = cost
	return newAction
}

// * define additional helper functions
func IncrementCurrentParameterIdx(idx *int) int {
	newInt := *idx + 1
	*idx = newInt
	return *idx
}

func PerformActions(actions *Stack, defaultParameterName string) {
	
	for _, _action := range actions.ToArray() {
		
		action := _action.(*Action)

		// TODO: find a way to incorporate additional conditional into whether to perform action
		if true {

			// check if requirements are fulfilled
			condition := true
			for _, requirement := range action.Requirements {
				if requirement.Parameter == "" {
					requirement.Parameter = defaultParameterName
				}
				x := Parameters.Get(FIND_Key, requirement.Parameter).Val.(*Stack).Get(FIND_Key, requirement.Layer).Val.(*float32)
				switch requirement.MinOrMax {
				case Min:
					condition = condition && *x >= requirement.Threshold
				case Max:
					condition = condition && *x <= requirement.Threshold
				}
			}

			// requirements have been fulfilled
			if condition {
				// perform updates
				for _, update := range action.Updates {
					*Parameters.Get(FIND_Key, update.Parameter).Val.(*Stack).Get(FIND_Key, update.Layer).Val.(*float32) -= update.Dx
				}
				
			}

		}
	}
	
}

// * define timeIncrement functions
func TimeIncrement(x_ui, dx_ui, tprev_s, tcur_s, dt_s float32) (xprime_ui float32) {
	
	// compute xprime
	Δt_s  := tcur_s - tprev_s
	Δx_ui := Δt_s * (dx_ui / dt_s)
	xprime_ui = clampUI(x_ui + Δx_ui)
	
	// return
	return
	
}
