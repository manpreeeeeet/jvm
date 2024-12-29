package main

import "strings"

type Field struct {
	accessFlags u2
	name        string
	descriptor  string
	attributes  []Attribute
	value       interface{}
}

type Method struct {
	accessFlags u2
	name        string
	descriptor  string
	attributes  []Attribute
	value       interface{}
}
type MethodDescriptorDetailed struct {
	params     string
	returnType string
}

func findMethodDescriptorDetailed(descriptor string) (*MethodDescriptorDetailed, error) {
	start := strings.Index(descriptor, "(")
	end := strings.Index(descriptor, ")")

	params := descriptor[start+1 : end]
	returnTypes := descriptor[end+1:]
	return &MethodDescriptorDetailed{
		params:     params,
		returnType: returnTypes,
	}, nil
}

func (method Method) findCodeAttribute() (*CodeAttribute, error) {
	for _, attribute := range method.attributes {
		if attribute.name == "Code" {
			return attribute.toCodeAttribute(), nil
		}
	}
	panic("not a code attribute")
}
