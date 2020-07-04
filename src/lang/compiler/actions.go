package compiler

import "strings"
import "strconv"
import "reflect"
import "fmt"
import "os"
import "encoding/gob"

import . "lang/interpreter"

var operations = []string{ "+", "-", "*", "/", "^", "%", "&", "|", "=", "!=", ">", "<", ">=", "<=", ")", "(", "~~", "~~~", ":" }

func convToAct(_val []interface{}, dir, name string) []Action {
  var val []Action

  if reflect.TypeOf(_val[0]).String() == "compiler.Lex" {

    var num []Lex

    for _, v := range _val {
      num = append(num, v.(Lex))
    }

    val, _ = Actionizer(num, true, dir, name)

  } else {

    for _, v := range _val {
      val = append(val, v.(Action))
    }

  }

  return val
}

func getLeft(index int, exp []interface{}, dir, name string) ([]Action, []interface{}) {

  var _num1 []interface{}

  //_num1 loop
  for o := index - 1; o >= 0; o-- {

    _num1 = append(_num1, exp[o])
  }

  reverseInterface(_num1)

  num1 := convToAct(_num1, dir, name)

  return num1, _num1
}

func getRight(index int, exp []interface{}, dir, name string) ([]Action, []interface{}) {
  var _num2 []interface{}

  //_num2 loop
  for o := index + 1; o < len(exp); o++ {

    _num2 = append(_num2, exp[o])
  }

  num2 := convToAct(_num2, dir, name)

  return num2, _num2
}

func calcExp(index int, exp []interface{}, dir, name string) ([]Action, []Action, []interface{}, []interface{}) {

  num1, _num1 := getLeft(index, exp, dir, name)
  num2, _num2 := getRight(index, exp, dir, name)

  return num1, num2, _num1, _num2
}

func callCalcParams(i *int, lex []Lex, len_lex int, dir, filename string) ([][]Action, [][]Action, []SubCaller, bool) {

  cbCnt := 0
  glCnt := 0
  bCnt := 0
  pCnt := 1

  indexes := [][]Lex{[]Lex{}}
  var putIndexes [][]Action

  if lex[*i].Name == "." {

    cbCnt = 0
    glCnt = 0
    bCnt = 0
    pCnt = 0

    for o := (*i) + 1; o < len_lex; o++ {
      if lex[o].Name == "{" {
        cbCnt++
      }
      if lex[o].Name == "[:" {
        glCnt++
      }
      if lex[o].Name == "[" {
        bCnt++
      }
      if lex[o].Name == "(" {
        pCnt++
      }

      if lex[o].Name == "}" {
        cbCnt--
      }
      if lex[o].Name == ":]" {
        glCnt--
      }
      if lex[o].Name == "]" {
        bCnt--
      }
      if lex[o].Name == ")" {
        pCnt--
      }

      if lex[o].Name == "." {
        indexes = append(indexes, []Lex{})
      } else {

        (*i)++

        indexes[len(indexes) - 1] = append(indexes[len(indexes) - 1], lex[o])

        if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {

          if o < len_lex - 1 && lex[o + 1].Name == "." {
            continue
          } else {
            break
          }

        }
      }
    }

    for _, v := range indexes {
      temp, _ := Actionizer(v[1:len(v) - 1], true, dir, filename)
      putIndexes = append(putIndexes, temp)
    }

    (*i)++
  }

  var params_ [][]Action
  var subcaller []SubCaller

  var isProc = false

  if *i < len_lex && lex[*i].Name == "(" {

    params := [][]Lex{[]Lex{}}

    cbCnt = 0
    glCnt = 0
    bCnt = 0
    pCnt = 0

    for o := *i; o < len_lex; o++ {
      if lex[o].Name == "{" {
        cbCnt++;
      }
      if lex[o].Name == "}" {
        cbCnt--;
      }

      if lex[o].Name == "[:" {
        glCnt++;
      }
      if lex[o].Name == ":]" {
        glCnt--;
      }

      if lex[o].Name == "[" {
        bCnt++;
      }
      if lex[o].Name == "]" {
        bCnt--;
      }

      if lex[o].Name == "(" {
        pCnt++;
      }
      if lex[o].Name == ")" {
        pCnt--;
      }

      if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
        break
      }

      if o == *i {
        continue
      }

      //detect a new argument
      if lex[o].Name == "," && cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 1 {
        params = append(params, []Lex{})
        continue
      }

      params[len(params) - 1] = append(params[len(params) - 1], lex[o])
    }

    for _, v := range params {

      if len(v) == 0 {
        continue
      }

      temp, _ := Actionizer(v, true, dir, filename)
      params_ = append(params_, temp)
    }

    pCnt_ := 0
    skip_nums := 0

    for o := *i; o < len_lex; o++ {
      if lex[o].Name == "(" {
        pCnt_++
      }
      if lex[o].Name == ")" {
        pCnt_--
      }

      skip_nums++;

      if pCnt_ == 0 {
        break
      }
    }

    isProc = true

    (*i)+=skip_nums

    //detect a subcaller
    //subcaller means test_fn()() //<-- the last () is the subcaller
    //a subcaller can also be test_fn().[0] //<-- the .[0] is a subcaller
    if *i < len_lex {

      if lex[*i].Name == "(" || lex[*i].Name == "." {

        paramsSub, indexesSub, subVal, isProcSub := callCalcParams(i, lex, len_lex, dir, filename)

        subcaller = append(subcaller, SubCaller{ indexesSub, paramsSub, isProcSub })
        subcaller = append(subcaller, subVal...)
      }
    }
  }

  return params_, putIndexes, subcaller, isProc
}

//function to actionize the callers (#~ and @~)
func callCalc(i *int, lex []Lex, len_lex int, dir, filename string) ([][]Action, [][]Action, []SubCaller, string) {

  var name = lex[*i + 2].Name

  (*i)+=3

  params_, putIndexes, subcaller, _ := callCalcParams(i, lex, len_lex, dir, filename)

  return params_, putIndexes, subcaller, name
}

func fnCalc(i *int, lex []Lex, len_lex int, dir, name string) ([]Action, []string, string) {

  var params []string
  var fnName string
  var logic []Action

  if *i + 1 < len(lex) && lex[*i + 1].Name == "~" {
    fnName = lex[*i + 2].Name

    for o := (*i) + 4; o < len_lex; o++ {
      if lex[o].Name == ")" {
        break
      }

      if lex[o].Name == "," {
        (*i)++
        continue
      }

      params = append(params, lex[o].Name)
    }
    *i+=(len(params) + 5)

    var logic_ = []Lex{}

    cbCnt := 0

    for o := *i; o < len_lex; o++ {
      if lex[o].Name == "{" {
        cbCnt++
      }

      if lex[o].Name == "}" {
        cbCnt--
      }

      logic_ = append(logic_, lex[o])

      if cbCnt == 0 {
        break
      }
    }

    (*i)+=len(logic_) - 1

    logic, _ = Actionizer(logic_, false, dir, name)
  } else {
    params = []string{}
    fnName = ""

    for o := (*i) + 2; o < len_lex; o+=2 {
      if lex[o].Name == ")" {
        break
      }

      params = append(params, lex[o].Name)
    }
    *i+=(3 + len(params))

    var logic_ = []Lex{}

    cbCnt := 0

    for o := *i; o < len_lex; o++ {
      if lex[o].Name == "{" {
        cbCnt++
      }

      if lex[o].Name == "}" {
        cbCnt--
      }

      logic_ = append(logic_, lex[o])

      if cbCnt == 0 {
        break
      }
    }

    (*i)+=len(logic_) + 1

    logic, _ = Actionizer(logic_, false, dir, name)
  }

  return logic, params, fnName
}

//export Actionizer
func Actionizer(lex []Lex, doExpress bool, dir, name string) ([]Action, map[string][]Action) {
  var actions = []Action{}
  var len_lex = len(lex)
  var variables = make(map[string][]Action) //to store the variables declared in this scope

  for i := 0; i < len_lex; i++ {

    if doExpress { //if it is an expression
      var exp []interface{}

      cbCnt := 0
      glCnt := 0
      bCnt := 0
      pCnt := 0

      for o := i; o < len_lex; o++ {
        if lex[o].Name == "{" {
          cbCnt++;
        }
        if lex[o].Name == "}" {
          cbCnt--;
        }

        if lex[o].Name == "[:" {
          glCnt++;
        }
        if lex[o].Name == ":]" {
          glCnt--;
        }

        if lex[o].Name == "[" {
          bCnt++;
        }
        if lex[o].Name == "]" {
          bCnt--;
        }

        if lex[o].Name == "(" {
          pCnt++;
        }
        if lex[o].Name == ")" {
          pCnt--;
        }

        if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
          break
        }

        exp = append(exp, lex[o])

        i++
      }

      for ;interfaceContainOperations(exp, "|") || interfaceContainOperations(exp, "&") || interfaceContainOperations(exp, "!|") || interfaceContainOperations(exp, "!&") || interfaceContainOperations(exp, "$|") || interfaceContainOperations(exp, "!$|"); {
        indexes := map[string]int{
          "|": interfaceIndexOfOperations("|", exp),
          "&": interfaceIndexOfOperations("&", exp),
          "!|": interfaceIndexOfOperations("!|", exp),
          "!&": interfaceIndexOfOperations("!&", exp),
          "$|": interfaceIndexOfOperations("$|", exp),
          "!$|": interfaceIndexOfOperations("!$|", exp),
        }

        //get max index
        var max = [2]interface{}{}

        for k, v := range indexes {
          if v != -1 {
            max = [2]interface{}{ k, v }
          }
        }

        for k, v := range indexes {
          if (v != -1 && v < max[1].(int)) || max[1].(int) == -1 {
            max = [2]interface{}{ k, v }
          }
        }

        switch max[0].(string) {
          case "|":
            index := max[1].(int)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "or", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case "&":
            index := max[1].(int)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "and", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case "!|":
            index := max[1].(int)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "nor", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case "!&":
            index := max[1].(int)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "nand", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case "$|":
            index := max[1].(int)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "xor", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case "!$|":
            index := max[1].(int)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "xnor", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
        }
      }

      for ;interfaceContainOperations(exp, "=") || interfaceContainOperations(exp, "!=") || interfaceContainOperations(exp, ">") || interfaceContainOperations(exp, "<") || interfaceContainOperations(exp, ">=") || interfaceContainOperations(exp, "<=") || interfaceContainOperations(exp, "~~") || interfaceContainOperations(exp, "~~~"); {
        indexes := map[string]int{
          "=": interfaceIndexOfOperations("=", exp),
          "!=": interfaceIndexOfOperations("!=", exp),
          ">": interfaceIndexOfOperations(">", exp),
          "<": interfaceIndexOfOperations("<", exp),
          ">=": interfaceIndexOfOperations(">=", exp),
          "<=": interfaceIndexOfOperations("<=", exp),
          "~~": interfaceIndexOfOperations("~~", exp),
          "~~~": interfaceIndexOfOperations("~~~", exp),
        }

        //get max index
        var max = [2]interface{}{}

        for k, v := range indexes {
          if v != -1 {
            max = [2]interface{}{ k, v }
          }
        }

        for k, v := range indexes {
          if (v != -1 && v < max[1].(int)) || max[1].(int) == -1 {
            max = [2]interface{}{ k, v }
          }
        }

        switch max[0].(string) {
          case "=":
            index := max[1].(int)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "equals", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case "!=":
            index := max[1].(int)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "notEqual", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case ">":
            index := max[1].(int)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "greater", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case "<":
            index := max[1].(int)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "less", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case ">=":
            index := max[1].(int)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "greaterOrEqual", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case "<=":
            index := max[1].(int)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "lessOrEqual", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case "~~":
            index := max[1].(int)

            var degree_ []interface{}
            doDeg := false

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            for o := index + 1; o < len(exp); o++ {
              if exp[o] == "{" {
                cbCnt++;
              }
              if exp[o] == "}" {
                cbCnt--;
              }

              if exp[o] == "[:" {
                glCnt++;
              }
              if exp[o] == ":]" {
                glCnt--;
              }

              if exp[o] == "[" {
                bCnt++;
              }
              if exp[o] == "]" {
                bCnt--;
              }

              if exp[o] == "(" {
                pCnt++;
              }
              if exp[o] == ")" {
                pCnt--;
              }

              if reflect.TypeOf(exp[o]).String() == "compiler.Lex" && exp[o].(Lex).Name == ":" && cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
                doDeg = true
                break
              }

              degree_ = append(degree_, exp[o])
            }

            var degree = []Action{}
            var addDeg = 0

            if doDeg {
              degree = convToAct(degree_, dir, name)
              addDeg = len(degree_) + 1
            }

            num1, _num1 := getLeft(index, exp, dir, name)
            num2, _num2 := getRight(index + addDeg, exp, dir, name)

            var act_exp = Action{ "similar", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, degree, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + addDeg + 1:]...)

            exp = exp_
          case "~~~":
            index := max[1].(int)

            var degree_ []interface{}
            doDeg := false

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            for o := index + 1; o < len(exp); o++ {
              if exp[o] == "{" {
                cbCnt++;
              }
              if exp[o] == "}" {
                cbCnt--;
              }

              if exp[o] == "[:" {
                glCnt++;
              }
              if exp[o] == ":]" {
                glCnt--;
              }

              if exp[o] == "[" {
                bCnt++;
              }
              if exp[o] == "]" {
                bCnt--;
              }

              if exp[o] == "(" {
                pCnt++;
              }
              if exp[o] == ")" {
                pCnt--;
              }

              if reflect.TypeOf(exp[o]).String() == "compiler.Lex" && exp[o].(Lex).Name == ":" && cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
                doDeg = true
                break
              }

              degree_ = append(degree_, exp[o])
            }

            var degree = []Action{}
            var addDeg = 0

            if doDeg {
              degree = convToAct(degree_, dir, name)
              addDeg = len(degree_) + 1
            }

            num1, _num1 := getLeft(index, exp, dir, name)
            num2, _num2 := getRight(index + addDeg, exp, dir, name)

            var act_exp = Action{ "strictSimilar", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, degree, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + addDeg + 1:]...)

            exp = exp_
        }
      }

      for ;interfaceContainOperations(exp, "+") || interfaceContainOperations(exp, "-"); {

        if (interfaceIndexOfOperations("+", exp) < interfaceIndexOfOperations("-", exp) && interfaceIndexOfOperations("+", exp) != -1) || interfaceIndexOfOperations("-", exp) == -1 {
          index := interfaceIndexOfOperations("+", exp)

          num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

          var act_exp = Action{ "add", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

          exp_ := append(exp[:index - len(_num1)], act_exp)
          exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

          exp = exp_
        } else {
          index := interfaceIndexOfOperations("-", exp)

          num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

          var act_exp = Action{ "subtract", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

          exp_ := append(exp[:index - len(_num1)], act_exp)
          exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

          exp = exp_
        }

      }

      for ;interfaceContainOperations(exp, "*") || interfaceContainOperations(exp, "/") || interfaceContainOperations(exp, "%"); {

        indexes := map[string]int{
          "*": interfaceIndexOfOperations("*", exp),
          "/": interfaceIndexOfOperations("/", exp),
          "%": interfaceIndexOfOperations("%", exp),
        }

        //get max index
        var max = [2]interface{}{}

        for k, v := range indexes {
          if v != -1 {
            max = [2]interface{}{ k, v }
          }
        }

        for k, v := range indexes {
          if (v != -1 && v < max[1].(int)) || max[1].(int) == -1 {
            max = [2]interface{}{ k, v }
          }
        }

        switch max[0].(string) {
          case "*":
            index := interfaceIndexOfOperations("*", exp)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "multiply", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case "/":
            index := interfaceIndexOfOperations("/", exp)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "divide", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case "%":
            index := interfaceIndexOfOperations("%", exp)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "modulo", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
        }

      }

      for ;interfaceContainOperations(exp, "^"); {
        index := interfaceIndexOfOperations("^", exp)

        num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

        var act_exp = Action{ "exponentiate", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

        exp_ := append(exp[:index - len(_num1)], act_exp)
        exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

        exp = exp_
      }

      for ;interfaceContainOperations(exp, "!"); {

        index := interfaceIndexOfOperations("!", exp)

        var num []interface{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := index + 1; o < len(exp); o++ {

          if exp[o].(Lex).Name == "{" {
            cbCnt++
          }
          if exp[o].(Lex).Name == "}" {
            cbCnt--
          }

          if exp[o].(Lex).Name == "[" {
            bCnt++
          }
          if exp[o].(Lex).Name == "]" {
            bCnt--
          }

          if exp[o].(Lex).Name == "[:" {
            glCnt++
          }
          if exp[o].(Lex).Name == ":]" {
            glCnt--
          }

          if exp[o].(Lex).Name == "(" {
            pCnt++
          }
          if exp[o].(Lex).Name == ")" {
            pCnt--
          }

          if arrayContainInterface(operations, exp[o]) {
            break
          }

          num = append(num, exp[o])
        }

        numAct := convToAct(num, dir, name)

        var act_exp = Action{ "not", "operation", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, numAct, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

        exp_ := append(exp[:index], act_exp)
        exp_ = append(exp_, exp[index + len(num) + 1:]...)

        exp = exp_
      }

      var fn_indexes []int

      for ;interfaceContainWithProcIndex(exp, "(", fn_indexes); {

        index := interfaceIndexOfWithProcIndex("(", exp, fn_indexes)

        if index - 1 != -1 && (reflect.TypeOf(exp[index - 1]).String() != "compiler.Lex" || ((strings.HasPrefix(exp[index - 1].(Lex).Name, "$") || exp[index - 1].(Lex).Name == "]")))  {
          fn_indexes = append(fn_indexes, index)
          continue
        }

        var pExp []Lex

        pCnt := 0

        for o := index; o < len(exp); o++ {
          if exp[o].(Lex).Name == "(" {
            pCnt++;
          }
          if exp[o].(Lex).Name == ")" {
            pCnt--;
          }

          pExp = append(pExp, exp[o].(Lex))

          if pCnt == 0 {
            break
          }
        }

        pExp = pExp[1:len(pExp) - 1]

        pExpAct, _ := Actionizer(pExp, true, dir, name)

        scbCnt := 0
        sglCnt := 0
        sbCnt := 0
        spCnt := 0

        indexes := [][]Lex{}

        if !(index + len(pExp) + 2 >= len(exp)) {
          if exp[index + len(pExp) + 2].(Lex).Name == "." {
            for o := index + len(pExp) + 2; o < len_lex; o++ {
              if exp[o].(Lex).Name == "{" {
                scbCnt++
              }
              if exp[o].(Lex).Name == "}" {
                scbCnt--
              }

              if exp[o].(Lex).Name == "[" {
                sbCnt++
              }
              if exp[o].(Lex).Name == "]" {
                sbCnt--
              }

              if exp[o].(Lex).Name == "[:" {
                sglCnt++
              }
              if exp[o].(Lex).Name == ":]" {
                sglCnt--
              }

              if exp[o].(Lex).Name == "(" {
                spCnt++
              }
              if exp[o].(Lex).Name == ")" {
                spCnt--
              }

              if exp[o].(Lex).Name == "." {
                indexes = append(indexes, []Lex{})
              } else {

                i++

                indexes[len(indexes) - 1] = append(indexes[len(indexes) - 1], exp[o].(Lex))

                if scbCnt == 0 && sglCnt == 0 && sbCnt == 0 && spCnt == 0 {

                  if o < len(exp) - 1 && exp[o + 1].(Lex).Name == "." {
                    continue
                  } else {
                    break
                  }

                }
              }
            }

            var putIndexes [][]Action

            for _, v := range indexes {

              v = v[1:len(v) - 1]
              temp, _ := Actionizer(v, true, dir, name)
              putIndexes = append(putIndexes, temp)
            }

            pExpAct[0].Type = "expressionIndex"
            pExpAct[0].Indexes = putIndexes
          }
        }

        exp = append([]interface{}{ pExpAct[0] }, exp...)
      }

      if len(exp) == 0 {
        break
      }

      if reflect.TypeOf(exp[0]).String() == "compiler.Lex" {

        //variale that grets convved to a []Lex
        var toa []Lex

        for _, v := range exp {
          toa = append(toa, v.(Lex))
        }

        temp, _ := Actionizer(toa, false, dir, name)
        exp[0] = temp[0]
      }

      actions = append(actions, exp[0].(Action))
    }

    if i >= len_lex {
      break
    }

    switch lex[i].Name {
      case "newlineN":
        actions = append(actions, Action{ "newline", "", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
      case "local":
        exp_ := []Lex{}

        //getting nb semicolons
        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 4; o < len_lex; o++ {

          if lex[o].Name == "{" {
            cbCnt++;
          }
          if lex[o].Name == "}" {
            cbCnt--;
          }

          if lex[o].Name == "[:" {
            glCnt++;
          }
          if lex[o].Name == ":]" {
            glCnt--;
          }

          if lex[o].Name == "[" {
            bCnt++;
          }
          if lex[o].Name == "]" {
            bCnt--;
          }

          if lex[o].Name == "(" {
            pCnt++;
          }
          if lex[o].Name == ")" {
            pCnt--;
          }

          if cbCnt != 0 || glCnt != 0 || bCnt != 0 || pCnt != 0 {
            exp_ = append(exp_, lex[o])
            continue
          }

          if lex[o].Name == "newlineS" {
            break
          }

          exp_ = append(exp_, lex[o])
        }

        exp, _ := Actionizer(exp_, true, dir, name)

        variables[lex[i + 2].Name] = exp
        actions = append(actions, Action{ "local", lex[i + 2].Name, "", exp, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
        i+=(4 + len(exp_))
      case "alt":

        var alter = Action{ "alt", "", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

        pCnt := 0

        cond_ := []Lex{}

        for o := i + 1; o < len_lex; o++ {
          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          cond_ = append(cond_, lex[o])

          if pCnt == 0 {
            break
          }
        }

        i+=len(cond_) + 1

        cond, _ := Actionizer(cond_, true, dir, name)

        for do := true; do || lex[i].Name == "=>"; do = false {

          adder := 0

          if lex[i].Name != "=>" {
            adder = 1
          }

          cbCnt := 0

          actions_ := []Lex{}

          for o := i + adder; o < len_lex; o++ {
            if lex[o].Name == "{" {
              cbCnt++
            }
            if lex[o].Name == "}" {
              cbCnt--
            }

            actions_ = append(actions_, lex[o])

            if cbCnt == 0 {
              break
            }
          }

          i+=len(actions_)
          actions, _ := Actionizer(actions_, true, dir, name)

          alter.Condition = append(alter.Condition, Condition{ "alt", cond, actions })
          i++

          if i >= len_lex {
            break
          }
        }

        actions = append(actions, alter)

      case "global":
        exp_ := []Lex{}

        //getting nb semicolons
        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 4; o < len_lex; o++ {

          if lex[o].Name == "{" {
            cbCnt++;
          }
          if lex[o].Name == "}" {
            cbCnt--;
          }

          if lex[o].Name == "[:" {
            glCnt++;
          }
          if lex[o].Name == ":]" {
            glCnt--;
          }

          if lex[o].Name == "[" {
            bCnt++;
          }
          if lex[o].Name == "]" {
            bCnt--;
          }

          if lex[o].Name == "(" {
            pCnt++;
          }
          if lex[o].Name == ")" {
            pCnt--;
          }

          if cbCnt != 0 || glCnt != 0 || bCnt != 0 || pCnt != 0 {
            exp_ = append(exp_, lex[o])
            continue
          }

          if lex[o].Name == "newlineS" {
            break
          }

          exp_ = append(exp_, lex[o])
        }

        exp, _ := Actionizer(exp_, true, dir, name)

        variables[lex[i + 2].Name] = exp
        actions = append(actions, Action{ "global", lex[i + 2].Name, "", exp, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
        i+=(4 + len(exp_))
      case "log":
        exp_ := []Lex{}

        //getting nb semicolons
        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {

          if lex[o].Name == "{" {
            cbCnt++;
          }
          if lex[o].Name == "}" {
            cbCnt--;
          }

          if lex[o].Name == "[:" {
            glCnt++;
          }
          if lex[o].Name == ":]" {
            glCnt--;
          }

          if lex[o].Name == "[" {
            bCnt++;
          }
          if lex[o].Name == "]" {
            bCnt--;
          }

          if lex[o].Name == "(" {
            pCnt++;
          }
          if lex[o].Name == ")" {
            pCnt--;
          }

          if cbCnt != 0 || glCnt != 0 || bCnt != 0 || pCnt != 0 {
            exp_ = append(exp_, lex[o])
            continue
          }

          if lex[o].Name == "newlineS" {
            break
          }

          exp_ = append(exp_, lex[o])
        }

        exp, _ := Actionizer(exp_, true, dir, name)

        actions = append(actions, Action{ "log", "", "", exp, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
        i+=(2 + len(exp_))
      case "print":
        exp_ := []Lex{}

        //getting nb semicolons
        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {

          if lex[o].Name == "{" {
            cbCnt++;
          }
          if lex[o].Name == "}" {
            cbCnt--;
          }

          if lex[o].Name == "[:" {
            glCnt++;
          }
          if lex[o].Name == ":]" {
            glCnt--;
          }

          if lex[o].Name == "[" {
            bCnt++;
          }
          if lex[o].Name == "]" {
            bCnt--;
          }

          if lex[o].Name == "(" {
            pCnt++;
          }
          if lex[o].Name == ")" {
            pCnt--;
          }

          if cbCnt != 0 || glCnt != 0 || bCnt != 0 || pCnt != 0 {
            exp_ = append(exp_, lex[o])
            continue
          }

          if lex[o].Name == "newlineS" {
            break
          }

          exp_ = append(exp_, lex[o])
        }

        exp, _ := Actionizer(exp_, false, dir, name)

        actions = append(actions, Action{ "print", "", "", exp, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
        i+=(2 + len(exp_))
      case "{":
        exp_ := []Lex{}

        //getting nb semicolons
        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i; o < len_lex; o++ {

          if lex[o].Name == "{" {
            cbCnt++;
          }
          if lex[o].Name == "}" {
            cbCnt--;
          }

          if lex[o].Name == "[:" {
            glCnt++;
          }
          if lex[o].Name == ":]" {
            glCnt--;
          }

          if lex[o].Name == "[" {
            bCnt++;
          }
          if lex[o].Name == "]" {
            bCnt--;
          }

          if lex[o].Name == "(" {
            pCnt++;
          }
          if lex[o].Name == ")" {
            pCnt--;
          }

          if cbCnt != 0 || glCnt != 0 || bCnt != 0 || pCnt != 0 {
            exp_ = append(exp_, lex[o])
            continue
          }

          exp_ = append(exp_, lex[o])

          if cbCnt == 0 {
            break
          }

          if lex[o].Name == "newlineS" {
            break
          }
        }

        exp_ = exp_[1:]
        exp_ = exp_[:len(exp_) - 1]

        exp, _ := Actionizer(exp_, false, dir, name)

        actions = append(actions, Action{ "group", "", "", exp, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
        i+=(len(exp_) + 1)
      case "function":

        putFalsey := make(map[string][]Action)

        logic, params, fnName := fnCalc(&i, lex, len_lex, dir, name)

        act := Action{ "function", fnName, "", logic, params, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, putFalsey, "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }
        variables[fnName] = []Action{ act }
        actions = append(actions, act)
      case "fargc":

        i+=2
        count := lex[i].Name
        i++

        //if it is not a number of a parameter list
        //they are just as follows
        /*

        fargc ~ (string, number, ...whatever datatypes) {

        }

        */
        //used to properly overload functions
        if getType(count) != "number" && count != "(" {

          //throw an error
          colorprint("Error while actionizing in " + lex[i].Dir + "!\n", 12)
          fmt.Println("Expected either a numeric value or a parameter list after fargc but instead got", count, "which is of type", getType(count), "\n\nError occured on line", lex[i].Line, "\nFound near:", strings.TrimSpace(lex[i].Exp))

          //exit the process
          os.Exit(1)
        }

        var types []string

        if count == "(" {

          for ;i < len_lex; i++ {

            if lex[i].Name == "," {
              continue
            }

            if lex[i].Name == ")" {
              i++
              break
            }

            types = append(types, lex[i].Name[1:])

            if !isType(lex[i].Name[1:]) {

              //throw an error
              colorprint("Error while actionizing in " + lex[i].Dir + "!\n", 12)
              fmt.Println("Expected a type value instead of", lex[i].Name[1:], "\n\nError occured on line", lex[i].Line, "\nFound near:", strings.TrimSpace(lex[i].Exp))

              //exit the process
              os.Exit(1)
            }
          }
        }

        if lex[i].Name != "{" {

          //throw an error
          colorprint("Error while actionizing in " + lex[i].Dir + "!\n", 12)
          fmt.Println("Expected { instead of", lex[i].Name, "\nFound near:", strings.TrimSpace(lex[i].Exp))

          //exit the process
          os.Exit(1)
        }

        cbCnt := 0

        var exp []Lex

        for ;i < len_lex; i++ {

          if lex[i].Name == "{" {
            cbCnt++
            continue
          }
          if lex[i].Name == "}" {
            cbCnt--
            continue
          }

          exp = append(exp, lex[i])

          if cbCnt == 0 {
            break
          }
        }
        i--

        actionized, _ := Actionizer(exp, false, dir, name)

        if count == "(" {

          actions = append(actions, Action{ "fargc_paramlist", "", "", actionized, types, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
        } else {

          if getType(count) != "number" || strings.ContainsRune(count, '.') {
            //throw an error
            colorprint("Error while actionizing in " + lex[i].Dir + "!\n", 12)
            fmt.Println("Expected an integer instead of", count, "\nFound near:", strings.TrimSpace(lex[i - 1].Exp))

            //exit the process
            os.Exit(1)
          }

          ommNumCountInt, _ := BigNumConverter(count) //convert to omm number

          actions = append(actions, Action{ "fargc_number", "", "", actionized, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, ommNumCountInt, []int64{}, OmmThread{} })
        }
      case "#":

        params_, putIndexes, subcaller, name := callCalc(&i, lex, len_lex, dir, name)

        actions = append(actions, Action{ "#", name, "", []Action{}, []string{}, params_, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, putIndexes, make(map[string][]Action), "private", subcaller, []int64{}, []int64{}, OmmThread{} })
      case "@":

        params_, putIndexes, subcaller, name := callCalc(&i, lex, len_lex, dir, name)

        actions = append(actions, Action{ "@", name, "", []Action{}, []string{}, params_, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, putIndexes, make(map[string][]Action), "private", subcaller, []int64{}, []int64{}, OmmThread{} })
      case "return":

        returner_ := []Lex{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {
          if lex[o].Name == "{" {
            cbCnt++
          }
          if lex[o].Name == "}" {
            cbCnt--
          }

          if lex[o].Name == "[:" {
            glCnt++
          }
          if lex[o].Name == ":]" {
            glCnt--
          }

          if lex[o].Name == "[" {
            bCnt++
          }
          if lex[o].Name == "]" {
            bCnt--
          }

          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
            break
          }

          returner_ = append(returner_, lex[o])
        }

        returner, _ := Actionizer(returner_, true, dir, name)

        actions = append(actions, Action{ "return", "", "", returner, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
        i+=len(returner_) + 2
      case "if":

        conditions := []Condition{}

        if lex[i].Name == "if" {
          var cond_ = []Lex{}
          pCnt := 0

          for o := i + 1; o < len_lex; o++ {
            if lex[o].Name == "(" {
              pCnt++
            }
            if lex[o].Name == ")" {
              pCnt--
            }

            cond_ = append(cond_, lex[o])

            if pCnt == 0 {
              break
            }
          }

          cond, _ := Actionizer(cond_, true, dir, name)

          cbCnt := 0
          glCnt := 0
          bCnt := 0
          pCnt = 0

          actions_ := []Lex{}

          var curlyBraceCond = lex[i + 1 + len(cond_)].Name == "{"

          for o := i + 1 + len(cond_); o < len_lex; o++ {
            if lex[o].Name == "{" {
              cbCnt++
            }
            if lex[o].Name == "}" {
              cbCnt--
            }

            if lex[o].Name == "[:" {
              glCnt++
            }
            if lex[o].Name == ":]" {
              glCnt--
            }

            if lex[o].Name == "[" {
              bCnt++
            }
            if lex[o].Name == "]" {
              bCnt--
            }

            if lex[o].Name == "(" {
              pCnt++
            }
            if lex[o].Name == ")" {
              pCnt--
            }

            actions_ = append(actions_, lex[o])

            if !curlyBraceCond {

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
                break
              }
            } else if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
              break
            }
          }

          acts, _ := Actionizer(actions_, false, dir, name)

          var if_ = Condition{ "if", cond, acts }

          conditions = append(conditions, if_)

          i+=(1 + len(cond_) + len(actions_))
        }

        if !(i >= len_lex) {
          for ;lex[i].Name == "elseif"; {

            var cond_ = []Lex{}
            pCnt := 0

            for o := i + 1; o < len_lex; o++ {
              if lex[o].Name == "(" {
                pCnt++
              }
              if lex[o].Name == ")" {
                pCnt--
              }

              cond_ = append(cond_, lex[o])

              if pCnt == 0 {
                break
              }
            }

            cond, _ := Actionizer(cond_, true, dir, name)

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt = 0

            actions_ := []Lex{}

            var curlyBraceCond = lex[i + 1 + len(cond_)].Name == "{"

            for o := i + 1 + len(cond_); o < len_lex; o++ {
              if lex[o].Name == "{" {
                cbCnt++
              }
              if lex[o].Name == "}" {
                cbCnt--
              }

              if lex[o].Name == "[:" {
                glCnt++
              }
              if lex[o].Name == ":]" {
                glCnt--
              }

              if lex[o].Name == "[" {
                bCnt++
              }
              if lex[o].Name == "]" {
                bCnt--
              }

              if lex[o].Name == "(" {
                pCnt++
              }
              if lex[o].Name == ")" {
                pCnt--
              }

              actions_ = append(actions_, lex[o])

              if !curlyBraceCond {

                if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
                  break
                }
              } else if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
                break
              }
            }

            acts, _ := Actionizer(actions_, false, dir, name)

            var elseif_ = Condition{ "elseif", cond, acts }

            conditions = append(conditions, elseif_)

            i+=(1 + len(cond_) + len(actions_))

          }
        }

        if !(i >= len_lex) {
          actions_ := []Lex{}

          for ;lex[i].Name == "else"; {
            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            //allow for user to write: else ~ <do something>;
            var curlyBraceCond = lex[i + 1].Name == "{"

            for o := i + 1; o < len_lex; o++ {
              if lex[o].Name == "{" {
                cbCnt++
              }
              if lex[o].Name == "}" {
                cbCnt--
              }

              if lex[o].Name == "[:" {
                glCnt++
              }
              if lex[o].Name == ":]" {
                glCnt--
              }

              if lex[o].Name == "[" {
                bCnt++
              }
              if lex[o].Name == "]" {
                bCnt--
              }

              if lex[o].Name == "(" {
                pCnt++
              }
              if lex[o].Name == ")" {
                pCnt--
              }

              if !curlyBraceCond {

                if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
                  break
                }
              } else if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
                break
              }

              actions_ = append(actions_, lex[o])
            }

            if !curlyBraceCond {
              actions_ = actions_[1:]
            }

            actions, _ := Actionizer(actions_, false, dir, name)

            var else_ = Condition{ "else", []Action{}, actions }

            conditions = append(conditions, else_)

            i+=(1 + len(actions_))

            if i >= len_lex {
              break
            }
          }
        }

        actions = append(actions, Action{ "conditional", "", "", []Action{}, []string{}, [][]Action{}, conditions, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
        i--
      case "import":

        var fileDir = lex[i + 2].Name

        //remove the quotes
        fileDir = fileDir[1:len(fileDir) - 1]

        var files = []map[string]string{}

        //see if user wants to import a file from the stdlib
        if strings.HasPrefix(fileDir, "?~") {

          if strings.HasPrefix(fileDir[2:], "/") {
            files = ReadFileJS("./stdlib" + fileDir[2:])
          } else {
            files = ReadFileJS("./stdlib/" + fileDir[2:])
          }

        } else {
          files = ReadFileJS(dir + fileDir)
        }

        var lexxed = []map[string]interface{}{}

        i+=2

        if i + 1 < len_lex && lex[i + 1].Name == "newlineS" {
          i++
        }

        for _, v := range files {

          if arrayContain(imported, v["FileName"]) {
            continue
          }

          imported = append(imported, v["FileName"])

          curlex := Lexer(v["Content"], dir, strings.TrimPrefix(v["FileName"], dir) /* remove the directory part of the filename */)

          curmap := map[string]interface{}{
            "FileName": v["FileName"],
            "Content": curlex,
          }

          lexxed = append(lexxed, curmap)
        }

        var actionizedFiles [][]Action

        for _, v := range lexxed {

          if strings.HasSuffix(v["FileName"].(string), ".oat") {

            readfile, _ := os.Open(v["FileName"].(string))

            var decoded []Action

            decoder := gob.NewDecoder(readfile)
            e := decoder.Decode(&decoded)

            if e != nil {
              colorprint("Error while actionizing " + dir + name + ", ", 12)
              fmt.Println(v["FileName"], "was detected as an oat, but is not oat compatible.")
              os.Exit(1)
            }

            readfile.Close()

            actionizedFiles = append(actionizedFiles, decoded)

          } else {
            temp, _ := Actionizer(v["Content"].([]Lex), false, dir, name)
            actionizedFiles = append(actionizedFiles, temp)
          }
        }

        actions = append(actions, Action{ "import", "", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, actionizedFiles, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
      case "break":
        actions = append(actions, Action{ "break", "", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
      case "skip":
        actions = append(actions, Action{ "skip", "", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
      case "loop":

        var condition_ = []Lex{}

        pCnt := 0

        for o := i + 1; o < len_lex; o++ {
          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          condition_ = append(condition_, lex[o])

          if pCnt == 0 {
            break
          }
        }

        condition, _ := Actionizer(condition_, true, dir, name)
        action_ := []Lex{}

        var curlyBraceCond = lex[i + 1 + len(condition_)].Name == "{"

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt = 0

        for o := i + 1 + len(condition_); o < len_lex; o++ {
          if lex[o].Name == "{" {
            cbCnt++
          }
          if lex[o].Name == "}" {
            cbCnt--
          }

          if lex[o].Name == "[:" {
            glCnt++
          }
          if lex[o].Name == ":]" {
            glCnt--
          }

          if lex[o].Name == "[" {
            bCnt++
          }
          if lex[o].Name == "]" {
            bCnt--
          }

          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          action_ = append(action_, lex[o])

          if !curlyBraceCond {

            if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
              break
            }
          } else if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
            break
          }
        }

        action, _ := Actionizer(action_, false, dir, name)

        actions = append(actions, Action{ "loop", "", "", action, []string{}, [][]Action{}, []Condition{ { "loop", condition, action } }, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
        i+=len(condition_) + len(action_)
      case "[:":
        var phrase = []Lex{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i; o < len_lex; o++ {
          if lex[o].Name == "{" {
            cbCnt++
          }
          if lex[o].Name == "}" {
            cbCnt--
          }

          if lex[o].Name == "[:" {
            glCnt++
          }
          if lex[o].Name == ":]" {
            glCnt--
          }

          if lex[o].Name == "[" {
            bCnt++
          }
          if lex[o].Name == "]" {
            bCnt--
          }

          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          phrase = append(phrase, lex[o])

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
            break
          }
        }

        i+=len(phrase)

        phrase = phrase[1:len(phrase) - 1]

        var _translated =  [][][]Lex{ [][]Lex{ []Lex{}, []Lex{} } }

        cbCnt = 0
        glCnt = 0
        bCnt = 0
        pCnt = 0

        cur := 0

        for _, v := range phrase {

          if v.Name == "{" {
            cbCnt++
          }
          if v.Name == "}" {
            cbCnt--
          }

          if v.Name == "[:" {
            glCnt++
          }
          if v.Name == ":]" {
            glCnt--
          }

          if v.Name == "[" {
            bCnt++
          }
          if v.Name == "]" {
            bCnt--
          }

          if v.Name == "(" {
            pCnt++
          }
          if v.Name == ")" {
            pCnt--
          }

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && v.Name == ":" {
            cur = 1
            continue
          }
          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && v.Name == "," {
            cur = 0
            _translated = append(_translated, [][]Lex{ []Lex{}, []Lex{} })
            continue
          }

          _translated[len(_translated) - 1][cur] = append(_translated[len(_translated) - 1][cur], v)
        }

        var translated = make(map[string][]Action)

        for _, v := range _translated {

          if len(v[0]) <= 0 {
            break
          }

          var name_ = v[0][0].Name

          if strings.HasPrefix(v[0][0].Name, "'") {
            name_ = name_[1:len(name_) - 1]
          }

          if strings.HasPrefix(v[0][0].Name, "$") {
            name_ = name_[1:]
          }

          hashVal, _ := Actionizer(v[1], true, dir, name)

          translated[name_] = hashVal
        }

        i--

        if i >= len_lex {
          actions = append(actions, Action{ "hash", "hashed_value", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, translated, "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
          break
        }

        if lex[i].Name == "." {
          indexes := [][]Lex{ []Lex{} }

          cbCnt = 0
          glCnt = 0
          bCnt = 0
          pCnt = 0

          for o := i + 1; o < len_lex; o++ {
            if lex[o].Name == "{" {
              cbCnt++
            }
            if lex[o].Name == "[:" {
              glCnt++
            }
            if lex[o].Name == "[" {
              bCnt++
            }
            if lex[o].Name == "(" {
              pCnt++
            }

            if lex[o].Name == "}" {
              cbCnt--
            }
            if lex[o].Name == ":]" {
              glCnt--
            }
            if lex[o].Name == "]" {
              bCnt--
            }
            if lex[o].Name == ")" {
              pCnt--
            }

            if lex[o].Name == "." {
              indexes = append(indexes, []Lex{})
            } else {

              i++

              indexes[len(indexes) - 1] = append(indexes[len(indexes) - 1], lex[o])

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {

                if o < len_lex - 1 && lex[o + 1].Name == "." {
                  continue
                } else {
                  break
                }

              }
            }
          }

          var putIndexes [][]Action

          for _, v := range indexes {

            v = v[1:len(v) - 1]
            temp, _ := Actionizer(v, true, dir, name)
            putIndexes = append(putIndexes, temp)
          }

          i+=3

          actions = append(actions, Action{ "hashIndex", "", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, putIndexes, translated, "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
        } else {
          actions = append(actions, Action{ "hash", "hashed_value", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, translated, "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
        }
      case "[":
        var phrase = []Lex{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i; o < len_lex; o++ {

          if lex[o].Name == "{" {
            cbCnt++
          }
          if lex[o].Name == "}" {
            cbCnt--
          }

          if lex[o].Name == "[:" {
            glCnt++
          }
          if lex[o].Name == ":]" {
            glCnt--
          }

          if lex[o].Name == "[" {
            bCnt++
          }
          if lex[o].Name == "]" {
            bCnt--
          }

          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          phrase = append(phrase, lex[o])

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
            break
          }
        }

        i+=len(phrase)

        phrase = phrase[1:len(phrase) - 1]

        var arr [][]Action

        for o := 0; o < len(phrase); o++ {

          var sub []Lex

          cbCnt := 0
          glCnt := 0
          bCnt := 0
          pCnt := 0

          for j := o; j < len(phrase); j++ {

            if phrase[j].Name == "{" {
              cbCnt++
            }
            if phrase[j].Name == "}" {
              cbCnt--
            }

            if phrase[j].Name == "[:" {
              glCnt++
            }
            if phrase[j].Name == ":]" {
              glCnt--
            }

            if phrase[j].Name == "[" {
              bCnt++
            }
            if phrase[j].Name == "]" {
              bCnt--
            }

            if phrase[j].Name == "(" {
              pCnt++
            }
            if phrase[j].Name == ")" {
              pCnt--
            }

            if phrase[j].Name == "," && cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
              break
            }
            sub = append(sub, phrase[j])
          }

          o+=len(sub)
          temp, _ := Actionizer(sub, true, dir, name)
          arr = append(arr, temp)
        }

        hashedArr := make(map[string][]Action)

        cur := 0

        for _, v := range arr {
          hashedArr[strconv.Itoa(cur)] = v
          cur++
        }

        if i >= len_lex {
          actions = append(actions, Action{ "array", "hashed_value", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, arr, [][]Action{}, hashedArr, "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
          break
        }

        if lex[i].Name == "." {
          indexes := [][]Lex{ []Lex{} }

          cbCnt = 0
          glCnt = 0
          bCnt = 0
          pCnt = 0

          for o := i + 1; o < len_lex; o++ {
            if lex[o].Name == "{" {
              cbCnt++
            }
            if lex[o].Name == "[:" {
              glCnt++
            }
            if lex[o].Name == "[" {
              bCnt++
            }
            if lex[o].Name == "(" {
              pCnt++
            }

            if lex[o].Name == "}" {
              cbCnt--
            }
            if lex[o].Name == ":]" {
              glCnt--
            }
            if lex[o].Name == "]" {
              bCnt--
            }
            if lex[o].Name == ")" {
              pCnt--
            }

            if lex[o].Name == "." {
              indexes = append(indexes, []Lex{})
            } else {

              i++

              indexes[len(indexes) - 1] = append(indexes[len(indexes) - 1], lex[o])

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {

                if o < len_lex - 1 && lex[o + 1].Name == "." {
                  continue
                } else {
                  break
                }

              }
            }
          }

          var putIndexes [][]Action

          for _, v := range indexes {

            v = v[1:len(v) - 1]
            temp, _ := Actionizer(v, true, dir, name)
            putIndexes = append(putIndexes, temp)
          }

          i+=3

          actions = append(actions, Action{ "arrayIndex", "", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, arr, putIndexes, hashedArr, "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
        } else {
          actions = append(actions, Action{ "array", "hashed_value", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, arr, [][]Action{}, hashedArr, "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
        }
      case "each":
        var condition_ = []Lex{}

        pCnt := 0

        for o := i + 1; o < len_lex; o++ {
          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          condition_ = append(condition_, lex[o])

          if pCnt == 0 {
            break
          }
        }

        i+=len(condition_) + 1

        condition_ = condition_[1:len(condition_) - 1]

        var _iterator []Lex

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt = 0

        var stopIterIndex int

        for k, v := range condition_ {

          if v.Name == "{" {
            cbCnt++
          }
          if v.Name == "}" {
            cbCnt--
          }

          if v.Name == "[:" {
            glCnt++
          }
          if v.Name == ":]" {
            glCnt--
          }

          if v.Name == "[" {
            bCnt++
          }
          if v.Name == "]" {
            bCnt--
          }

          if v.Name == "(" {
            pCnt++
          }
          if v.Name == ")" {
            pCnt--
          }

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && v.Name == "," {

            //index where the iterator stopped
            stopIterIndex = k
            break
          }

          _iterator = append(_iterator, v)
          stopIterIndex = k
        }

        iterator, _ := Actionizer(_iterator, true, dir, name)

        var var1 string
        var var2 string

        if stopIterIndex + 1 >= len(condition_) || stopIterIndex + 3 >= len(condition_) {
          //throw an error
          colorprint("Error while actionizing in " + lex[i].Dir + "!\n", 12)
          fmt.Println("Expected two variables after iterator in \"each\"", "\n\nError occured on line", lex[i].Line, "\nFound near:", strings.TrimSpace(lex[i].Exp))

          //exit the process
          os.Exit(1)
        } else {
          var1, var2 = condition_[stopIterIndex + 1].Name, condition_[stopIterIndex + 3].Name
        }

        cbCnt = 0
        glCnt = 0
        bCnt = 0
        pCnt = 0

        var exp []Lex

        var curlyBraceCond = lex[i].Name == "{"

        for o := i; o < len_lex; o++ {
          if lex[o].Name == "{" {
            cbCnt++
          }
          if lex[o].Name == "}" {
            cbCnt--
          }

          if lex[o].Name == "[:" {
            glCnt++
          }
          if lex[o].Name == ":]" {
            glCnt--
          }

          if lex[o].Name == "[" {
            bCnt++
          }
          if lex[o].Name == "]" {
            bCnt--
          }

          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          exp = append(exp, lex[o])

          if !curlyBraceCond {

            if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
              break
            }
          } else if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
            break
          }
        }

        i+=len(exp) + 1
        actionized, _ := Actionizer(exp, false, dir, name)
        actions = append(actions, Action{ "each", "", "", actionized, []string{ var1, var2 }, [][]Action{}, []Condition{}, iterator, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })

      case "await":

        awaiter := []Lex{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {
          if lex[o].Name == "{" {
            cbCnt++
          }
          if lex[o].Name == "}" {
            cbCnt--
          }

          if lex[o].Name == "[:" {
            glCnt++
          }
          if lex[o].Name == ":]" {
            glCnt--
          }

          if lex[o].Name == "[" {
            bCnt++
          }
          if lex[o].Name == "]" {
            bCnt--
          }

          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
            break
          }

          awaiter = append(awaiter, lex[o])
        }

        awaiterAct, _ := Actionizer(awaiter, true, dir, name)

        actions = append(actions, Action{ "await", "", "", awaiterAct, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
        i+=len(awaiter) + 2

      default:

        valPuts := func(lex []Lex, i int) int {

          if i >= len_lex {
            return 1
          }

          val := lex[i].Name

          i++

          switch getType(val) {

            case "string": {

              noQ := val[1:len(val) - 1] //the string will be given like this: "hello", but omm needs to store them like this: hello
              hashedString := make(map[string][]Action)

              cur := 0

              for _, v := range noQ {
                hashedIndex := make(map[string][]Action)
                hashedString[strconv.Itoa(cur)] = []Action{ Action{ "rune", "exp_value", string(v), []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, hashedIndex, "public", []SubCaller{}, []int64{}, []int64{}, OmmThread{} } }
                cur++
              }

              actions = append(actions, Action{ "string", "exp_value", noQ, []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, hashedString, "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
            }
            case "rune":
              hashed := make(map[string][]Action)

              noQ := val[1:len(val) - 1] //see noQ for string

              actions = append(actions, Action{ "rune", "exp_value", noQ, []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, hashed, "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
            case "number":

              hashed := make(map[string][]Action)

              integer, decimal := BigNumConverter(val)

              actions = append(actions, Action{ "number", "exp_value", "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, hashed, "private", []SubCaller{}, integer, decimal, OmmThread{} })
            case "boolean":

              hashed := make(map[string][]Action)

              actions = append(actions, Action{ "boolean", "exp_value", val, []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, hashed, "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
            case "falsey":

              hashed := make(map[string][]Action)

              actions = append(actions, Action{ "falsey", "exp_value", val, []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, hashed, "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
            case "none":

              if strings.HasPrefix(val, "$") {

                actions = append(actions, Action{ "variable", val, val, []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
              } else {

                //get it? 42?
                actions = append(actions, Action{ "none", "exp_value", val, []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
              }
          }

          return 0
        }

        if i + 1 < len_lex {

          if lex[i + 1].Name == "->" {

            var val_ []Lex

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            for o := i + 2; o < len_lex; o++ {
              if lex[o].Name == "{" {
                cbCnt++
              }
              if lex[o].Name == "}" {
                cbCnt--
              }

              if lex[o].Name == "[:" {
                glCnt++
              }
              if lex[o].Name == ":]" {
                glCnt--
              }

              if lex[o].Name == "[" {
                bCnt++
              }
              if lex[o].Name == "]" {
                bCnt--
              }

              if lex[o].Name == "(" {
                pCnt++
              }
              if lex[o].Name == ")" {
                pCnt--
              }

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && arrayContainInterface(operations, lex[o].Name) {
                break
              }

              val_ = append(val_, lex[o])
            }

            val, _ := Actionizer(val_, true, dir, name)

            getValueType := func(val string) string {
              switch getType(val) {
                case "string": fallthrough
                case "number": fallthrough
                case "boolean": fallthrough
                case "falsey":
                  return "exp_value"
                case "array": fallthrough
                case "hash":
                  return "hashed_value"
              }

              return "exp_value"
            }

            actions = append(actions, Action{ "cast", lex[i].Name, getValueType(lex[i].Name), val, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
            i+=len(val_) + 2
            continue
          }

          if (lex[i + 1].Name == "++" || lex[i + 1].Name == "--") && strings.HasPrefix(lex[i].Name, "$") {

            actions = append(actions, Action{ lex[i + 1].Name, lex[i].Name, "", []Action{}, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
            i++
            continue;
          }

          if (lex[i + 1].Name == "+=" || lex[i + 1].Name == "-=" || lex[i + 1].Name == "*=" || lex[i + 1].Name == "/=" || lex[i + 1].Name == "%=" || lex[i + 1].Name == "^=") && strings.HasPrefix(lex[i].Name, "$") {

            var by_ []Lex

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            for o := i + 2; o < len_lex; o++ {
              if lex[o].Name == "{" {
                cbCnt++
              }
              if lex[o].Name == "}" {
                cbCnt--
              }

              if lex[o].Name == "[:" {
                glCnt++
              }
              if lex[o].Name == ":]" {
                glCnt--
              }

              if lex[o].Name == "[" {
                bCnt++
              }
              if lex[o].Name == "]" {
                bCnt--
              }

              if lex[o].Name == "(" {
                pCnt++
              }
              if lex[o].Name == ")" {
                pCnt--
              }

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
                break
              }

              by_ = append(by_, lex[o])
            }

            by, _ := Actionizer(by_, true, dir, name)

            actions = append(actions, Action{ lex[i + 1].Name, lex[i].Name, "", by, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
            continue;
          }

          doPutIndex := false

          icbCnt := 0
          iglCnt := 0
          ibCnt := 0
          ipCnt := 0

          for o := i; o < len_lex; o++ {
            if lex[o].Name == "{" {
              icbCnt++
            }
            if lex[o].Name == "}" {
              icbCnt--
            }

            if lex[o].Name == "[:" {
              iglCnt++
            }
            if lex[o].Name == ":]" {
              iglCnt--
            }

            if lex[o].Name == "[" {
              ibCnt++
            }
            if lex[o].Name == "]" {
              ibCnt--
            }

            if lex[o].Name == "(" {
              ipCnt++
            }
            if lex[o].Name == ")" {
              ipCnt--
            }

            if icbCnt == 0 && iglCnt == 0 && ibCnt == 0 && ipCnt == 0 && lex[o].Name == "newlineS" {
              break
            }

            if icbCnt == 0 && iglCnt == 0 && ibCnt == 0 && ipCnt == 0 && lex[o].Name == ":" {
              doPutIndex = true
              break
            }
          }

          var indexes [][]Action
          varname := lex[i].Name //store the variable name in variable (because we cannot get it later)

          if lex[i + 1].Name == "." && lex[i + 2].Name == "[" && doPutIndex {

            _indexes := [][]Lex{}

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            for o := i + 1; o < len_lex; i, o = i + 1, o + 1 {
              if lex[o].Name == "{" {
                cbCnt++
              }
              if lex[o].Name == "}" {
                cbCnt--
              }

              if lex[o].Name == "[:" {
                glCnt++
              }
              if lex[o].Name == ":]" {
                glCnt--
              }

              if lex[o].Name == "[" {
                bCnt++
              }
              if lex[o].Name == "]" {
                bCnt--
              }

              if lex[o].Name == "(" {
                pCnt++
              }
              if lex[o].Name == ")" {
                pCnt--
              }

              if lex[o].Name == "." {
                _indexes = append(_indexes, []Lex{})
                continue
              }

              _indexes[len(_indexes) - 1] = append(_indexes[len(_indexes) - 1], lex[o])

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o + 1].Name == ":" {
                break
              }
            }

            for _, v := range _indexes {
              temp, _ := Actionizer(v[1:len(v) - 1], true, dir, name)
              indexes = append(indexes, temp)
            }

            i++
          }

          if lex[i + 1].Name == ":" && (strings.HasPrefix(lex[i].Name, "$") || lex[i].Name == "]") {

            exp_ := []Lex{}

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            for o := i + 2; o < len_lex; o++ {

              if lex[o].Name == "{" {
                cbCnt++
              }
              if lex[o].Name == "}" {
                cbCnt--
              }

              if lex[o].Name == "[:" {
                glCnt++
              }
              if lex[o].Name == ":]" {
                glCnt--
              }

              if lex[o].Name == "[" {
                bCnt++
              }
              if lex[o].Name == "]" {
                bCnt--
              }

              if lex[o].Name == "(" {
                pCnt++
              }
              if lex[o].Name == ")" {
                pCnt--
              }

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
                break
              }

              exp_ = append(exp_, lex[o]);
            }

            exp, _ := Actionizer(exp_, true, dir, name)

            act := Action{ "let", varname, "", exp, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, indexes, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} }

            //set the variable
            variables[varname] = exp
            actions = append(actions, act)
            i+=(len(exp_)) - 1
            continue
          }

          if lex[i + 1].Name == "." {

            val := lex[i].Name

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            indexes := [][]Lex{ []Lex{} }

            cbCnt = 0
            glCnt = 0
            bCnt = 0
            pCnt = 0

            for o := i + 2; o < len_lex; o++ {
              if lex[o].Name == "{" {
                cbCnt++
              }
              if lex[o].Name == "[:" {
                glCnt++
              }
              if lex[o].Name == "[" {
                bCnt++
              }
              if lex[o].Name == "(" {
                pCnt++
              }

              if lex[o].Name == "}" {
                cbCnt--
              }
              if lex[o].Name == ":]" {
                glCnt--
              }
              if lex[o].Name == "]" {
                bCnt--
              }
              if lex[o].Name == ")" {
                pCnt--
              }

              if lex[o].Name == "." {
                indexes = append(indexes, []Lex{})
              } else {

                i++

                indexes[len(indexes) - 1] = append(indexes[len(indexes) - 1], lex[o])

                if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {

                  if o < len_lex - 1 && lex[o + 1].Name == "." {
                    continue
                  } else {
                    break
                  }

                }
              }
            }

            var putIndexes [][]Action

            for _, v := range indexes {

              v = v[1:len(v) - 1]
              temp, _ := Actionizer(v, true, dir, name)
              putIndexes = append(putIndexes, temp)
            }

            i+=3

            if strings.HasPrefix(val, "$") {
              actVal, _ := Actionizer([]Lex{ Lex{ val, "", 0, "", "", dir } }, true, dir, name)

              actions = append(actions, Action{ "variableIndex", "", "", actVal, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, putIndexes, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
            } else {
              actVal, _ := Actionizer([]Lex{ Lex{ val, "", 0, "", "", dir } }, true, dir, name)

              actions = append(actions, Action{ "expressionIndex", "", "", actVal, []string{}, [][]Action{}, []Condition{}, []Action{}, []Action{}, []Action{}, [][]Action{}, putIndexes, make(map[string][]Action), "private", []SubCaller{}, []int64{}, []int64{}, OmmThread{} })
            }

          }

          valPuts(lex, i)

        } else {

          valPuts(lex, i)
        }
      }
  }

  return actions, variables
}
