package main

import (
	"fmt"
	"os"
)

func main() {
	file, _ := os.Open("Main.class")
	loader := &ClassLoader{reader: file}

	class := loader.loadClass()
	err := file.Close()
	if err != nil {
		return
	}

	var mainCode *CodeAttribute
	for _, method := range class.methods {
		if method.name == "main" {
			mainCode = (&method.attributes[0]).toCodeAttribute()
			break
		}
	}

	frame := mainCode.toFrame(&class)

	jvm := JVM{classes: make([]*Class, 0)}
	jvm.classes = append(jvm.classes, &class)
	jvm.Exec(&frame)

	if frame.stack[len(frame.stack)-1] == nil {
		fmt.Printf("JVM exited with status code 0")
	}

}
