package main

import (
	"fmt"
)

type Frame struct {
	stack              []interface{}
	localVariables     []interface{}
	code               []byte
	instructionPointer uint32
	class              *Class
}

func (frame *Frame) getIndex() u2 {
	indexByte1 := frame.code[frame.instructionPointer+1]
	indexByte2 := frame.code[frame.instructionPointer+2]
	frame.instructionPointer = frame.instructionPointer + 2
	return (u2(indexByte1) << 8) | u2(indexByte2)
}

func (jvm *JVM) Exec(frame *Frame) interface{} {

	for {
		op := frame.code[frame.instructionPointer]

		switch op {
		case 0: // noop

		case 4: // iconst_1
			frame.push(1)
		case 5: // iconst_2
			frame.push(2)
		case 6: // iconst_3
			frame.push(3)
		case 7: // iconst_4
			frame.push(4)
		case 26: // iload_0
			frame.push(frame.localVariables[0])
		case 27: // iload_1
			frame.push(frame.localVariables[1])
		case 42: //	aload0
			frame.push(frame.localVariables[0])
		case 43: // aload1
			frame.push(frame.localVariables[1])
		case 60: // istore_1
			frame.localVariables[1] = frame.pop().(int)
		case 61: // istore2
			frame.localVariables[2] = frame.pop().(int)
		case 76: // astore_1
			frame.localVariables[1] = frame.pop()
		case 87: // pop
			frame.pop()
		case 89: // dup
			last := frame.stack[len(frame.stack)-1]
			frame.stack = append(frame.stack, last)
		case 96: // iadd
			first := frame.pop().(int)
			second := frame.pop().(int)
			frame.stack = append(frame.stack, first+second)
		case 104: // imul
			first := frame.pop().(int)
			second := frame.pop().(int)
			frame.stack = append(frame.stack, first*second)
		case 172: // ireturn
			return frame.pop()
		case 177: // return
			return nil
		case 180: // getField
			field := frame.class.constantPool.resolveFieldRef(frame.getIndex())
			objectRef := frame.pop().(*Object)
			frame.push(objectRef.fields[field.name])
		case 181: // putField
			field := frame.class.constantPool.resolveFieldRef(frame.getIndex())
			value := frame.pop().(int)
			objectRef := frame.pop().(*Object)
			objectRef.fields[field.name] = value
		case 182, 183: // invokespecial, invokeVirtual

			methodRefIndex := frame.getIndex()
			methodRef := frame.class.constantPool.resolveMethodRef(methodRefIndex)

			methodDescriptorDetailed, _ := findMethodDescriptorDetailed(methodRef.descriptor)
			params, returnType := methodDescriptorDetailed.params, methodDescriptorDetailed.returnType
			args := frame.stack[len(frame.stack)-len(params)-1:]
			frame.stack = frame.stack[:len(frame.stack)-len(params)-1]

			result, err := jvm.executeMethod(methodRef.class, methodRef.name, methodRef.descriptor, args...)
			if err != nil {
				return nil
			}

			if len(returnType) != 0 && returnType != "V" {
				frame.stack = append(frame.stack, result)
			}

		case 184: // invokestatic

			methodRefIndex := frame.getIndex()
			methodRef := frame.class.constantPool.resolveMethodRef(methodRefIndex)

			methodDescriptorDetailed, _ := findMethodDescriptorDetailed(methodRef.descriptor)
			params, returnType := methodDescriptorDetailed.params, methodDescriptorDetailed.returnType
			args := frame.stack[len(frame.stack)-len(params):]
			frame.stack = frame.stack[:len(frame.stack)-len(params)]

			result, err := jvm.executeMethod(methodRef.class, methodRef.name, methodRef.descriptor, args...)
			if err != nil {
				return nil
			}

			if len(returnType) != 0 && returnType != "V" {
				frame.stack = append(frame.stack, result)
			}

		case 187: // new
			classIndex := frame.getIndex()

			className := frame.class.constantPool.resolveString(frame.class.constantPool[classIndex-1].info.ClassIndex)
			class, err := jvm.getClass(className)
			if err != nil {
				return nil
			}

			frame.push(class.new())

		default:
			fmt.Printf("wooops unimplemented op code %d\n", op)
		}

		frame.instructionPointer++

	}
}
