package interpreter

func add(num1, num2 Action, cli_params CliParams) Action {

  /* TABLE OF TYPES:

    string + (* - array - hash) = string
    array + (* - none) = array
    num + num = num
    hash + hash = hash
    boolean + boolean = boolean
    num + boolean = num
    default = falsey
  */

  type1 := num1.Type
  type2 := num2.Type

  var final Action

  if (type1 == "string" && (type2 != "array" && type2 != "hash")) || (type2 == "string" && (type1 != "array" && type2 != "hash")) {
    //detect case `string + (* - array - hash) = string`
  }

  return final
}
