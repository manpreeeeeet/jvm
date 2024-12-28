package main

import "encoding/binary"

type Attribute struct {
	name string
	info []byte
}

type ExceptionItem struct {
	startPc   u2
	endPc     u2
	handlerPc u2
	catchType u2
}

type CodeAttribute struct {
	name           string
	maxStack       u2
	maxLocals      u2
	code           []byte
	exceptionTable []ExceptionItem
	attributes     []Attribute
}

func (codeAttribute *CodeAttribute) toFrame(class Class, args ...interface{}) Frame {

	frame := Frame{
		code:               codeAttribute.code,
		localVariables:     make([]interface{}, codeAttribute.maxLocals, codeAttribute.maxLocals),
		instructionPointer: 0,
		class:              class,
	}

	for i := 0; i < len(args); i++ {
		frame.localVariables[i] = args[i]
	}

	return frame
}

func (attribute *Attribute) toCodeAttribute() CodeAttribute {
	codeAttribute := CodeAttribute{
		name:      attribute.name,
		maxStack:  binary.BigEndian.Uint16(attribute.info[0:2]),
		maxLocals: binary.BigEndian.Uint16(attribute.info[2:4]),
	}
	codeLength := binary.BigEndian.Uint32(attribute.info[4:8])
	codeAttribute.code = attribute.info[8 : 8+codeLength]
	return codeAttribute
}
