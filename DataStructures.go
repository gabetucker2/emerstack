package main

func SetupDataStructures() {
	
	paramsStack = MakeStack(map[string]Parameter{

		//env
		"friend":          {true, -1, -1},
		"desk":            {true, -1, -1},
		"food":            {true, -1, -1},
		"mate":            {true, -1, -1},
		"bed":             {true, -1, -1},
		"socialsituation": {true, -1, -1},
		"danger":          {true, -1, -1},
	
		//inp
		"affiliation":   {false, 0.02, -0.167},
		"achievement":   {false, 0.0208, -0.083},
		"hunger":        {false, 0.014, -0.143},
		"sex":           {false, 0.0012, -0.333},
		"sleep":         {false, 0.005, -0.0104},
		"socialanxiety": {false, -1, -1}, // contingent on env - update later
		"fear":          {false, -1, -1}, // contingent on env - update later
		
	})
	envStack = paramsStack.GetMany(FIND_Lambda, func(card *Card) bool {
		return card.Val.(Parameter).envNotInp
	})
	inpStack = paramsStack.GetMany(FIND_Lambda, func(card *Card) bool {
		return !card.Val.(Parameter).envNotInp
	})

}
