package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"go.zoe.im/injgo"
)

const (
	helpContent = `injgo is a injector written in Go.
USAGE:
  injgo PROCESS_NAME/PROCESS_ID Libraies...

  EXAMPLES:
    1. Inject test.dll to process Calc.exe
      $ injgo Calc.exe test.dll

    2. Inject test.dll and demo.dll to process with PID: 1888
	  $ injgo 1888 test.dll demo.dll
`
)

func main() {
	if len(os.Args) <= 2 {
		fmt.Println(helpContent)
		return
	}

	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		p, err := injgo.FindProcessByName(os.Args[1])
		if err != nil {
			fmt.Println("can't find process:", os.Args[1], "error:", err)
			return
		}
		pid = p.ProcessID
	}

	// find pid and or

	fmt.Println("injector ", pid)
	for _, name := range os.Args[2:] {
		// check if file exits
		name, _ = filepath.Abs(name)
		err = injgo.Inject(pid, name, false)
		if err != nil {
			fmt.Println("inject ", name, "error:", err)
		} else {
			fmt.Println("inject ", name, "success")
		}
	}
}
