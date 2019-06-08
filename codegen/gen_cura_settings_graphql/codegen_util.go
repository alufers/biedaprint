package main

import (
	"encoding/json"
	"fmt"
)

func makeDirective(name string, arguments map[string]string) string {
	argumentsStr := ""
	if arguments != nil {
		argumentsStr += "\n"
		for k, v := range arguments {
			argumentsStr += fmt.Sprintf("    %v: %v\n", k, v)
		}
	}
	return fmt.Sprintf("@%v(%v)", name, argumentsStr)
}

func makeStringLiteral(str string) string {
	data, _ := json.Marshal(str)
	return string(data)
}
