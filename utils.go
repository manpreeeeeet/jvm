package main

func (frame *Frame) pop() interface{} {
	value := (frame.operandStack)[len(frame.operandStack)-1]
	frame.operandStack = frame.operandStack[:len(frame.operandStack)-1]
	return value
}

func (frame *Frame) push(value interface{}) {
	frame.operandStack = append(frame.operandStack, value)
}
