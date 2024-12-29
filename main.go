package main

import (
	"fmt"
	"log"
)

func main() {
	jvm := JVM{classes: make([]*Class, 0)}
	_, err := jvm.addClass("Main")
	if err != nil {
		return
	}
	_, err = jvm.executeMethod("Main", "main", "([Ljava/lang/String;)V")
	if err != nil {
		fmt.Printf("JVM failed with error %v", err)
		return
	}

	log.Println("JVM exited successfully.")

}
