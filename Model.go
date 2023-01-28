// * set up dependencies
package main

import (
	. "github.com/gabetucker2/gostack"
)

// * main structure setup
func SetupModel() {
	
	////////////////////////////////////////////////////////////////////
	// DO NOT EDIT
	//  logistical
	*currentParameterIdx = 0
	propertyKeys := []string {"layerValues", "bs_ui", "dx_ui", "timeIncrements", "relations", "actions"}
	//  default vals
	tprev_s = 0
	tcur_s = 0
	////////////////////////////////////////////////////////////////////
	// UNITS
	// s  := seconds
	// ui := unit interval [0, 1]
	////////////////////////////////////////////////////////////////////
	// EDIT BELOW
	dt_s = 1
	
	// initialize our Parameters variable
	Parameters = MakeStack(
	
		map[string]*Stack {
		
			"affiliation" : MakeStack(
				
				// property keys
				propertyKeys,
				
				// property vals
				[]any {
					
					// layers
					MakeStack(
						[]string {"enviro", "intero"}, // layer names
						[]*float32 {&enviro.Values[*currentParameterIdx], &intero.Values[IncrementCurrentParameterIdx(currentParameterIdx)]}, // layer addresses
						// (what's going on here is we need to procedurally update our currentParameterIdx value so that we don't need to type in [0], [1], etc every time)
						// (but we can't do so from inside a function call, so we sneakily do it by calling a function with a return value)
					),

					// bs_ui
					MakeStack(
						[]string {"intero"}, // layers
						[]float32 {0.5}, // corresponding values
					),

					// dx_ui
					-0.167,

					// timeIncrements
					MakeStack(
						[]string {"enviro", "intero"},
						[]func(float32, float32, float32, float32, float32) float32 {TimeIncrement, TimeIncrement},
					),

					// relations (assuming dt_s) (assuming change in this => how much do others change?)
					MakeStack(
						[]string {
							"Feel a sense of achievement when you are fulfilled from talking to others", // 1
							"Get social anxiety when others are in your environment", // 2
							"Remove social anxiety when you have a fulfilling social interaction", // 3
						},
						[]*Relation {
							MakeRelation("intero", "achievement", "intero", 0.05), // 1
							MakeRelation("enviro", "socialSituation", "intero", 0.07), // 2
							MakeRelation("intero", "socialSituation", "intero", -0.12), // 3
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
									MakeRequirement("socialSituation", "intero", Max, 0.5), // must not be too anxious
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
									MakeRequirement("socialSituation", "intero", Max, 0.34), // must not be too anxious
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

			"achievement" : MakeStack(

				// property keys
				propertyKeys,

				// property vals
				[]any {

					// layers
					MakeStack(
						[]string {"enviro", "intero"}, // layer names
						[]*float32 {&enviro.Values[*currentParameterIdx], &intero.Values[IncrementCurrentParameterIdx(currentParameterIdx)]}, // layer addresses
						// (what's going on here is we need to procedurally update our currentParameterIdx value so that we don't need to type in [0], [1], etc every time)
						// (but we can't do so from inside a function call, so we sneakily do it by calling a function with a return value)
					),

					// bs_ui
					0.7,

					// dx_ui
					-0.125,

					// timeIncrements
					MakeStack(
						[]string {"enviro", "intero"},
						[]func(float32, float32, float32, float32, float32) float32 {TimeIncrement, TimeIncrement},
					),

					// relations (assuming dt_s) (assuming change in this => how much do others change?)
					MakeStack(
						[]string {
							"Want to talk to others after making an achievement", // 1
							"Have more people around you after making an achievement", // 2
						},
						[]*Relation {
							MakeRelation("intero", "affiliation", "intero", -0.08), // 1
							MakeRelation("enviro", "affiliation", "enviro", 0.05), // 2
						},
					),

					// actions
					MakeStack(
						[]string {
							"Achieve goal", // 1
						},
						[]*Action {

							// 1
							MakeAction(
								// requirements for this action to be performed
								[]*Requirement {
									MakeSimpleRequirement("intero", Max, 0.4), // must want an achievement
								},
								// updates if the action is performed
								[]*Update {
									MakeSimpleUpdate("intero", 0.6), // fulfillment
									MakeSimpleUpdate("enviro", -0.2), // can't reachieve right after
								},
								// cost to perform action
								0.5,
							),

						},
					),

				},

			),

			"hunger" : MakeStack(

				// property keys
				propertyKeys,

				// property vals
				[]any {

					// layers
					MakeStack(
						[]string {"enviro", "intero"}, // layer names
						[]*float32 {&enviro.Values[*currentParameterIdx], &intero.Values[IncrementCurrentParameterIdx(currentParameterIdx)]}, // layer addresses
						// (what's going on here is we need to procedurally update our currentParameterIdx value so that we don't need to type in [0], [1], etc every time)
						// (but we can't do so from inside a function call, so we sneakily do it by calling a function with a return value)
					),

					// bs_ui
					0.4,

					// dx_ui
					-0.083,

					// timeIncrements
					MakeStack(
						[]string {"enviro", "intero"},
						[]func(float32, float32, float32, float32, float32) float32 {TimeIncrement, TimeIncrement},
					),

					// relations (assuming dt_s) (assuming change in this => how much do others change?)
					MakeStack(
						[]string {
							"Less inclined to talk to others while hungry", // 1
						},
						[]*Relation {
							MakeRelation("intero", "affiliation", "intero", 0.008), // 1
						},
					),

					// actions
					MakeStack(
						[]string {
							"Eat meal", // 1
						},
						[]*Action {

							// 1
							MakeAction(
								// requirements for this action to be performed
								[]*Requirement {
									MakeSimpleRequirement("intero", Max, 0.8), // must be hungry
								},
								// updates if the action is performed
								[]*Update {
									MakeSimpleUpdate("intero", 0.2), // fulfillment
								},
								// cost to perform action
								0.21,
							),

						},
					),

				},

			),

			"sex" : MakeStack(

				// property keys
				propertyKeys,

				// property vals
				[]any {

					// layers
					MakeStack(
						[]string {"enviro", "intero"}, // layer names
						[]*float32 {&enviro.Values[*currentParameterIdx], &intero.Values[IncrementCurrentParameterIdx(currentParameterIdx)]}, // layer addresses
						// (what's going on here is we need to procedurally update our currentParameterIdx value so that we don't need to type in [0], [1], etc every time)
						// (but we can't do so from inside a function call, so we sneakily do it by calling a function with a return value)
					),

					// bs_ui
					0.8,

					// dx_ui
					-0.095,

					// timeIncrements
					MakeStack(
						[]string {"enviro", "intero"},
						[]func(float32, float32, float32, float32, float32) float32 {TimeIncrement, TimeIncrement},
					),

					// relations (assuming dt_s) (assuming change in this => how much do others change?)
					MakeStack(
						[]string {
							"Less inclined to talk to others while fulfilled in a special way", // 1
						},
						[]*Relation {
							MakeRelation("intero", "affiliation", "intero", 0.009), // 1
						},
					),

					// actions
					MakeStack(
						[]string {
							"Self-please", // 1
						},
						[]*Action {

							// 1
							MakeAction(
								// requirements for this action to be performed
								[]*Requirement {
									MakeSimpleRequirement("intero", Max, 0.7), // must desire this act
								},
								// updates if the action is performed
								[]*Update {
									MakeSimpleUpdate("intero", 0.23), // fulfillment
								},
								// cost to perform action
								0.16,
							),

						},
					),

				},

			),

			"sleep" : MakeStack(

				// property keys
				propertyKeys,

				// property vals
				[]any {

					// layers
					MakeStack(
						[]string {"enviro", "intero"}, // layer names
						[]*float32 {&enviro.Values[*currentParameterIdx], &intero.Values[IncrementCurrentParameterIdx(currentParameterIdx)]}, // layer addresses
						// (what's going on here is we need to procedurally update our currentParameterIdx value so that we don't need to type in [0], [1], etc every time)
						// (but we can't do so from inside a function call, so we sneakily do it by calling a function with a return value)
					),

					// bs_ui
					0.9,

					// dx_ui
					-0.14,

					// timeIncrements
					MakeStack(
						[]string {"enviro", "intero"},
						[]func(float32, float32, float32, float32, float32) float32 {TimeIncrement, TimeIncrement},
					),

					// relations (assuming dt_s) (assuming change in this => how much do others change?)
					MakeStack(
						[]string {
							"Less inclined to talk to others while sleepy", // 1
						},
						[]*Relation {
							MakeRelation("intero", "affiliation", "intero", 0.018), // 1
						},
					),

					// actions
					MakeStack(
						[]string {
							"Sleep", // 1
						},
						[]*Action {

							// 1
							MakeAction(
								// requirements for this action to be performed
								[]*Requirement {
									MakeSimpleRequirement("intero", Max, 0.6), // must desire this act
								},
								// updates if the action is performed
								[]*Update {
									MakeSimpleUpdate("intero", 0.4), // fulfillment
								},
								// cost to perform action
								0.19,
							),

						},
					),

				},

			),

			"socialSituation" : MakeStack(

				// property keys
				propertyKeys,

				// property vals
				[]any {

					// layers
					MakeStack(
						[]string {"enviro", "intero"}, // layer names
						[]*float32 {&enviro.Values[*currentParameterIdx], &intero.Values[IncrementCurrentParameterIdx(currentParameterIdx)]}, // layer addresses
						// (what's going on here is we need to procedurally update our currentParameterIdx value so that we don't need to type in [0], [1], etc every time)
						// (but we can't do so from inside a function call, so we sneakily do it by calling a function with a return value)
					),

					// bs_ui
					0.5,
					
					// dx_ui
					0,

					// timeIncrements
					MakeStack(
						[]string {"enviro", "intero"},
						[]func(float32, float32, float32, float32, float32) float32 {TimeIncrement, TimeIncrement},
					),

					// relations (assuming dt_s) (assuming change in this => how much do others change?)
					MakeStack(
						[]string {
							"Less inclined to talk to others while socially anxious", // 1
						},
						[]*Relation {
							MakeRelation("intero", "affiliation", "intero", 0.014), // 1
						},
					),

					// actions
					MakeStack(
						[]string {
							"Evade people", // 1
						},
						[]*Action {

							// 1
							MakeAction(
								// requirements for this action to be performed
								[]*Requirement {
									MakeSimpleRequirement("intero", Min, 0.8), // must desire this act
								},
								// updates if the action is performed
								[]*Update {
									MakeSimpleUpdate("intero", -0.3), // fulfillment
								},
								// cost to perform action
								0.18,
							),

						},
					),

				},

			),

			"danger" : MakeStack(

				// property keys
				propertyKeys,

				// property vals
				[]any {

					// layers
					MakeStack(
						[]string {"enviro", "intero"}, // layer names
						[]*float32 {&enviro.Values[*currentParameterIdx], &intero.Values[IncrementCurrentParameterIdx(currentParameterIdx)]}, // layer addresses
						// (what's going on here is we need to procedurally update our currentParameterIdx value so that we don't need to type in [0], [1], etc every time)
						// (but we can't do so from inside a function call, so we sneakily do it by calling a function with a return value)
					),

					// bs_ui
					0.5,

					// dx_ui
					0,

					// timeIncrements
					MakeStack(
						[]string {"enviro", "intero"},
						[]func(float32, float32, float32, float32, float32) float32 {TimeIncrement, TimeIncrement},
					),

					// relations (assuming dt_s) (assuming change in this => how much do others change?)
					MakeStack(
						[]string {
						},
						[]*Relation {
						},
					),

					// actions
					MakeStack(
						[]string {
							"Run", // 1
						},
						[]*Action {

							// 1
							MakeAction(
								// requirements for this action to be performed
								[]*Requirement {
									MakeSimpleRequirement("intero", Min, 0.6), // must desire this act
								},
								// updates if the action is performed
								[]*Update {
									MakeSimpleUpdate("intero", -0.4), // fulfillment
								},
								// cost to perform action
								0.38,
							),

						},
					),

				},

			),

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
					MakeRequirement("socialSituation", "intero", Max, 0.35), // must not be too anxious
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
	
	////////////////////////////////////////////////////////////////////

}
