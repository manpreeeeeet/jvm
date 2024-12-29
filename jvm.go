package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type JVM struct {
	classes []*Class
}

func (jvm *JVM) addClass(className string) (*Class, error) {
	class, err := jvm.getClass(className)
	if class != nil {
		return class, nil
	}

	file, err := os.Open(className + ".class")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	loader := &ClassLoader{reader: file}

	loadedClass := loader.loadClass()
	jvm.classes = append(jvm.classes, &loadedClass)
	return &loadedClass, nil
}

func (jvm *JVM) getClass(className string) (*Class, error) {
	for _, class := range jvm.classes {
		if class.name == className {
			return class, nil
		}
	}
	return nil, errors.New("failed to find class" + className)
}

func (jvm *JVM) executeMethod(className string, methodName string, methodDescriptor string, args ...interface{}) (interface{}, error) {

	if className == "java/lang/Object" && methodName == "<init>" && methodDescriptor == "()V" {
		return nil, nil
	}

	class, err := jvm.getClass(className)
	if err != nil {
		return nil, err
	}

	method, err := class.findMethod(methodName, methodDescriptor)
	if err != nil {
		return nil, err
	}

	code, err := method.findCodeAttribute()
	if err != nil {
		return nil, err
	}
	frame := code.toFrame(class, args...)
	result := jvm.Exec(&frame)

	argString := ""
	for i, arg := range args {
		if i > 0 {
			argString += ","
		}
		switch arg.(type) {
		case *Object:
			argString += "this"
		default:
			argString += fmt.Sprintf("%v", arg)
		}
	}
	log.Printf("%s.%s(%s) --> %v\n", className, methodName, argString, result)
	return result, nil
}
