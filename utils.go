package main

func (frame *Frame) pop() interface{} {
	value := (frame.stack)[len(frame.stack)-1]
	frame.stack = frame.stack[:len(frame.stack)-1]
	return value
}

func (frame *Frame) push(value interface{}) {
	frame.stack = append(frame.stack, value)
}
