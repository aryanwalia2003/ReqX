package scripting

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// ConsoleAPI provides standard logging methods to the JS VM.
type ConsoleAPI struct{}

func (c *ConsoleAPI) Log(args ...interface{}) {
	fmt.Print(color.BlueString("[JS LOG] "))
	fmt.Println(args...)
}

func (c *ConsoleAPI) Warn(args ...interface{}) {
	fmt.Print(color.YellowString("[JS WARN] "))
	fmt.Println(args...)
}

func (c *ConsoleAPI) Error(args ...interface{}) {
	fmt.Print(color.RedString("[JS ERR] "))
	fmt.Println(args...)
}

func (c *ConsoleAPI) Dir(obj map[string]interface{}) { //this takes a map of strings and interfaces and prints it in a formatted way, this belongs to the console api
	fmt.Println("[JS DIR]")
	b, _ := json.MarshalIndent(obj, "", "  ")
	fmt.Println(string(b))
} //this will print the object in a formatted way

func (c *ConsoleAPI) Table(args ...interface{}) { //this takes a slice of interfaces and prints it in a table format, this belongs to the console api
	// A lightweight implementation of console.table
	fmt.Println("[JS TABLE]")
	for _, arg := range args {
		if mapArg, ok := arg.(map[string]interface{}); ok {
			for k, v := range mapArg {
				fmt.Printf("| %-20s | %-30v |\n", k, v)
			}
		} else {
			fmt.Printf("| %-50v |\n", arg)
		}
	}
	fmt.Println(strings.Repeat("-", 55))
}
