package main

type Field struct {
	accessFlags u2
	name        string
	descriptor  string
	attributes  []Attribute
	value       interface{}
}

type Method Field

func (method *Method) findCodeAttribute() CodeAttribute {
	for _, attribute := range method.attributes {
		if attribute.name == "Code" {
			return attribute.toCodeAttribute()
		}
	}
	panic("not a code attribute")
}
