package initialize

import (
	"encoding/json"
	"os"
)

//this package is used to initialize the program
/*
	- comprehending project structure
	- dependency managment
*/

func Initialize(configFileName string) *Program {

	var prog Program

	//config is the tusk config file

	configFile, e := os.Open(configFileName)

	if e != nil {
		//error
		_ = e
	}

	var config ConfigData
	json.NewDecoder(configFile).Decode(&config)

	prog.Config = config

	//the starting package to store everything
	var startpack = &Directory{
		Name: "tusk",
	}

	parseDirectory(".", startpack)

	prog.Nested = startpack.Nested

	return &prog
}
