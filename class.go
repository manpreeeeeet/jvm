package main

import "errors"

type u1 = uint8
type u2 = uint16
type u4 = uint32
type u8 = uint64

type Class struct {
	magic        u4
	minorVersion u2
	majorVersion u2
	constantPool ConstantPool
	accessFlags  u2
	name         string
	superClass   string
	interfaces   []Interface
	fields       []Field
	methods      []Method
	attributes   []Attribute
}

type Object struct {
	class  *Class
	fields map[string]interface{}
}

type Interface struct {
	name string
}

func (class *Class) findMethod(name string, descriptor string) (*Method, error) {
	for i := 0; i < len(class.methods); i++ {
		if class.methods[i].name == name && class.methods[i].descriptor == descriptor {
			return &class.methods[i], nil
		}
	}
	return nil, errors.New("method not found " + name)
}

func (class *Class) new() *Object {
	return &Object{
		class:  class,
		fields: make(map[string]interface{}),
	}
}
