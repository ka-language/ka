#ifndef OMM_LOG_FORMAT_HPP_
#define OMM_LOG_FORMAT_HPP_

#include <iostream>
#include <windows.h>
#include <stdio.h>
#include <vector>
#include <map>

#include "casts.hpp"
#include "../bind.h"

namespace omm {

  void log_format(Action in, const CliParams cli_params, std::map<std::string, Variable> vars, int hash_spacing, std::string doPrint) {

    if (in.Type == "hash") {
      std::map<std::string, vector<Action>> hashvals = in.Hash_Values;

      if (hashvals.size() == 0) cout << "[::]" << (doPrint == "print" ? "" : "\n");
      else {
        cout << "[:" << endl;

        for (std::pair<std::string, std::vector<Action>> it : hashvals) {
          std::string key = it.first;
          std::vector<Action> _value = it.second;

          cout << std::string(hash_spacing, ' ') << key << ": ";
          log_format(_value[0], cli_params, vars, hash_spacing + 2, "log");
        }

        cout << std::string(hash_spacing - 2, ' ') << ":]" << (doPrint == "print" ? "" : "\n");
      }
    } else if (in.Type == "array") {
      std::map<std::string, std::vector<Action>> hashvals = in.Hash_Values;

      if (hashvals.size() == 0) cout << "[]" << (doPrint == "print" ? "" : "\n");
      else {
        cout << "[" << endl;

        for (std::pair<std::string, std::vector<Action>> it : hashvals) {
          std::string key = it.first;
          std::vector<Action> _value = it.second;

          std::cout << std::string(hash_spacing, ' ') << key << ": ";
          log_format(_value[0], cli_params, vars, hash_spacing + 2, "log");
        }

        std::cout << std::string(hash_spacing - 2, ' ') << "]" << (doPrint == "print" ? "" : "\n");
      }
    } else if (in.Type == "process" || in.Type == "group") std::cout << "{(PROCESS~ | GROUP~) " << "PARAM COUNT: " << in.Params.size() << "}" << (doPrint == "print" ? "" : "\n");
    else if (in.Type == "thread") std::cout << "{Promise for proc " << in.Name.substr(1 /* remove the $ */ ) << "}";
    else if (in.Name == "operation") {
      log_format(in.First[0], cli_params, vars, hash_spacing, "print");

      std::string op = in.Type;
      cout << " " << GetOp(&op[0]) << " ";
      log_format(in.Second[0], cli_params, vars, hash_spacing, "print");

    } else {

      //convert value to a string
      //then print it

      std::string val = cast(in, "string").ExpStr[0];

      std::cout << val << (doPrint == "print" ? "" : "\n");
    }

    std::cout << std::flush; //flush the stdout
  }

}

#endif
