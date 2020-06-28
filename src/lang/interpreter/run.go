package interpreter

import "fmt"

type CliParams map[string]map[string]interface{}

//number sizes

//export DigitSize
const DigitSize = 1;
//export MAX_DIGIT
const MAX_DIGIT = 10 * DigitSize - 1
//export MIN_DIGIT
const MIN_DIGIT = -1 * MAX_DIGIT

//////////////

//export RunInterpreter
func RunInterpreter(actions []Action, cli_params map[string]map[string]interface{}, dir string) {

  var t = one
  var q = one
  t.Decimal = []int64{ 1, 2 }
  q.Decimal = []int64{ 1, 2 }

  fmt.Println(isEqual(t, q))

  var vars = make(map[string]Variable)

  for k, v := range goprocs {
    vars["$" + k] = Variable{
      Type: "goproc",
      Name: "$" + k,
      GoProc: v,
    }
  }

  interpreter(actions, CliParams(cli_params), vars, false, []Action{}, dir)

  for _, v := range threads {
    v.WaitFor()
  }
}
