package main

import (
	"errors"
)

type JVM struct {
	classes []*Class
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

	return result, nil
}
