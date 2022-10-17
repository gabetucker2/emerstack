package main

import (
	. "github.com/gabetucker2/gostack"
)

var Parameters *Stack

func SetupDataStructures() {
	
	Parameters = MakeStack(

		// parameters keys
		[]string {"affiliation", "achievement", "hunger", "sex", "sleep", "socialAnxiety", "fear"},

		// parameters vals
		[]*Stack {

			// affiliation
			MakeStack(
				// property keys
				[]string {"enviro", "intero", "motiveBias", "behavior", "dx", "decrimentFunc", "actions", "relations"},

				// property vals
				[]any {
					// pointer to the corresponding environmental value to this parameter stored in enviro
					&enviro.Values[0],
					// pointer to the corresponding interoceptive value to this parameter stored in intero
					&intero.Values[0],
					// motiveBias
					
					// behavior
					// dx
					// decrimentFunc
					// actions
					// relations
				},
			),

			// TODO: fill the restetc

		},
	)

}
