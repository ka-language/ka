package interpreter

import . "lang/types"

//list of operations
var operations = map[string]func(val1, val2 OmmType, cli_params CliParams, line uint64, file string) OmmType {
	"number + number": number__plus__number,
	"number - number": number__minus__number,
	"number * number": number__times__number,
	"number / number": number__divide__number,
	"number % number": number__mod__number,
	"number ^ number": number__pow__number,
	"number = number": func(val1, val2 OmmType, cli_params CliParams, line uint64, file string) OmmType {

		var final = falsev

		if isEqual(val1.(OmmNumber), val2.(OmmNumber)) {
			final = truev
		}

		return final
	},
	"number != number": func(val1, val2 OmmType, cli_params CliParams, line uint64, file string) OmmType {

		var final = truev

		if isEqual(val1.(OmmNumber), val2.(OmmNumber)) {
			final = falsev
		}

		return final
	},
}