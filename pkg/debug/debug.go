package debug

import "fmt"

func Debug(args ...interface{}) {
	fmt.Print("\n")
	for _, a := range args {
		fmt.Println(a)
	}
	fmt.Print("\n")
}
