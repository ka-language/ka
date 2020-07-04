package interpreter

import "fmt"
import "strings"

func log_format(in Action, hash_spacing int, endl bool) {

  switch in.Type {

    case "hash":
      if len(in.Hash_Values) == 0 {
        fmt.Print("[::]")
        goto end
      }

      fmt.Println("[:")

      for k, v := range in.Hash_Values {

        if v[0].Access == "private" { //if it is private, do not print it
          continue
        }

        fmt.Print(strings.Repeat(" ", hash_spacing) +  k + ": ")
        log_format(v[0], hash_spacing + 2, true)
      }

      fmt.Print(strings.Repeat(" ", hash_spacing - 2) + ":]")

    case "array":
      if len(in.Hash_Values) == 0 {
        fmt.Print("[]")
        goto end
      }

      fmt.Println("[")

      for k, v := range in.Hash_Values {

        if v[0].Access == "private" { //if it is private, do not print it
          fmt.Println(strings.Repeat(" ", hash_spacing) + "::private::")
          continue
        }

        fmt.Print(strings.Repeat(" ", hash_spacing) +  k + ": ")
        log_format(v[0], hash_spacing + 2, true)
      }

      fmt.Print(strings.Repeat(" ", hash_spacing - 2) + "]")
    case "group":
      fmt.Print("{...}")
    case "function":
      fmt.Print("{...} ", "PARAM COUNT:", len(in.Params))
    case "operation":
      log_format(in.First[0], hash_spacing, false)
      fmt.Print("", getOp(in.Type), "")
      log_format(in.Second[0], hash_spacing, false)
    case "thread":
      fmt.Print("CHANNEL: use the `await` keyword to get the value")
    default:
      //cast to a string, then print
      str := cast(in, "string")
      fmt.Print(str.ExpStr)

  }

  end:
  if endl { //if it was "logged", print a newline
    fmt.Println()
  }
}
