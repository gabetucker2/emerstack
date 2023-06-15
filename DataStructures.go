// * set up dependencies
package main

import (
	"github.com/emer/etable/etensor"
	. "github.com/gabetucker2/gostack"
)

// * initialize variables
var Parameters, ComplexActions, Layers *Stack
var tprev_s, tcur_s, dt_s int
var propertyKeys []string

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

func MakeAction(requirements []*Requirement, updates []*Update) *Action {
	newAction := new(Action)
	newAction.Requirements = requirements
	newAction.Updates = updates
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
		// TODO: make update conditions more complex, add complexactions, add clampUI
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

// Initialize immutable values
func InitializeImmutables() {
	
	//  logistical
	propertyKeys = []string {"bs_ui", "dxs_ui", "timeIncrements", "relations", "actions"} // "layerValues" is added to the beginning retroactively
	//  default vals
	tprev_s = 0
	tcur_s = 0
	
}

// Finish initializing the model
func FinishInitializing() {
	
	// add a "layerValues" reference to each parameter's corresponding tensor
	for i, _paramStack := range Parameters.ToArray() {
		paramStack := _paramStack.(*Stack)
		layerVals := MakeStack()
		for _, _layerVal := range Layers.ToArray() {
			layerVal := _layerVal.(*etensor.Float32)
			layerVals.Add(layerVal.Values[i])
		}
		insertStack := MakeStack(Layers.GetMany(nil, nil, RETURN_Keys), layerVals)
		paramStack.Add(MakeCard("layerValues", insertStack), ORDER_Before, FIND_First)
	}

	// fill up the tsrsStack property in Sim
	TheSim.tsrsStack = Layers.Clone()
	for _, layer := range Layers.ToArray() {
		var prevTsr *etensor.Float32
		TheSim.tsrsStack.Add(MakeCard(layer.(string) + "Prev", prevTsr))
	}
	
}
