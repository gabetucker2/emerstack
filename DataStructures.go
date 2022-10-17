// * setup dependencies
package main

import (
	. "github.com/gabetucker2/gostack"

)

// * initialize variables
var Parameters *Stack
var ComplexActions *Stack
var currentParameterIdx *int

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
func AsymptoteZero(x, dt, dx float32) float32 { // TODO: work on this more later
	return dt * (x - (dx*x))
}

// * main structure setup
func SetupDataStructures() {

	*currentParameterIdx = 0
	
	// initialize our Parameters variable
	Parameters = MakeStack(

		// parameters keys
		[]string {"affiliation", "achievement", "hunger", "sex", "sleep", "socialAnxiety", "fear"},

		// parameters vals
		[]*Stack {

			// affiliation
			MakeStack(

				// property keys
				[]string {"layerValues", "dt", "dx", "timeIncrements", "relations", "actions"},

				// property vals
				[]any {

					// layers
					MakeStack(
						[]string {"enviro", "intero"}, // layer names
						[]*float32 {&enviro.Values[*currentParameterIdx], &intero.Values[IncrementCurrentParameterIdx(currentParameterIdx)]}, // layer addresses
						// (what's going on here is we need to procedurally update our currentParameterIdx value so that we don't need to type in [0], [1], etc every time)
						// (but we can't do so from inside a function call, so we sneakily do it by calling a function with a return value)
					),

					// dt (in seconds)
					1,

					// dx
					-0.167,

					// timeIncrements
					MakeStack(
						[]string {"enviro", "intero"},
						[]func(float32, float32, float32) float32 {AsymptoteZero, AsymptoteZero},
					),

					// relations (assuming dt) (assuming change in this => how much do others change?)
					MakeStack(
						[]string {
							"Feel a sense of achievement when you are fulfilled from talking to others", // 1
							"Get social anxiety when others are in your environment", // 2
							"Remove social anxiety when you have a fulfilling social interaction", // 3
						},
						[]*Relation {
							MakeRelation("intero", "achievement", "intero", 0.05), // 1
							MakeRelation("enviro", "socialAnxiety", "intero", 0.07), // 2
							MakeRelation("intero", "socialAnxiety", "intero", -0.12), // 3
						},
					),

					// actions
					MakeStack(
						[]string {
							"Hang out with friend", // 1
							"Call friend", // 2
						},
						[]*Action {

							// 1
							MakeAction(
								// requirements for this action to be performed
								[]*Requirement {
									MakeSimpleRequirement("intero", Max, 0.6), // must have social battery
									MakeRequirement("socialAnxiety", "intero", Max, 0.5), // must not be too anxious
								},
								// updates if the action is performed
								[]*Update {
									MakeSimpleUpdate("intero", 0.4), // fulfillment
									MakeSimpleUpdate("enviro", -0.1), // friend leaves after
								},
								// cost to perform action
								0.4,
							),

							// 2
							MakeAction(
								// requirements for this action to be performed
								[]*Requirement {
									MakeSimpleRequirement("intero", Min, 0.3), // must have social battery
									MakeRequirement("socialAnxiety", "intero", Max, 0.34), // must not be too anxious
								},
								// updates if the action is performed
								[]*Update {
									MakeSimpleUpdate("intero", 0.4), // fulfillment
								},
								// cost to perform action
								0.18,
							),

						},
					),

				},

			),

			// TODO: fill the rest

		},
	)

	// create complex actions
	ComplexActions = MakeStack(
		[]string {
			"Grab food with friend", // 1
		},
		[]*Action {

			// 1
			MakeAction(
				// requirements for this action to be performed
				[]*Requirement {
					MakeRequirement("affiliation", "intero", Max, 0.7), // must have social battery
					MakeRequirement("socialAnxiety", "intero", Max, 0.35), // must not be too anxious
					MakeRequirement("food", "intero", Max, 0.5), // must not be too full
				},
				// updates if the action is performed
				[]*Update {
					MakeUpdate("affiliation", "intero", 0.3), // fulfillment
					MakeUpdate("affiliation", "enviro", -0.07), // friend leaves after
					MakeUpdate("food", "intero", 0.5), // food satisfied
					MakeUpdate("food", "enviro", -0.3), // food gone after
				},
				// cost to perform action
				0.27,
			),

		},
	)

}
