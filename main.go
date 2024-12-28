package main

import (
	"fmt"
	"os"
)

func main() {
	file, _ := os.Open("Main.class")
	loader := &ClassLoader{reader: file}

	class := loader.loadClass()

	var mainCode CodeAttribute
	for _, method := range class.methods {
		if method.name == "main" {
			mainCode = (&method.attributes[0]).toCodeAttribute()
			break
		}
	}

	frame := mainCode.toFrame(class)
	frame.Exec()

	if frame.operandStack[len(frame.operandStack)-1] == nil {
		fmt.Printf("JVM exited with status code 0")
	}

}
