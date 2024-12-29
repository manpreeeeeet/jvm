package main

import (
	"fmt"
	"os"
)

type Frame struct {
	stack              []interface{}
	localVariables     []interface{} // function parameters basically
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
			frame.stack = append(frame.stack, 1)
		case 5: // iconst_2
			frame.stack = append(frame.stack, 2)
		case 6: // iconst_3
			frame.stack = append(frame.stack, 3)
		case 7: // iconst_4
			frame.stack = append(frame.stack, 4)
		case 17: // sipush
		case 26: // iload_0
			frame.stack = append(frame.stack, frame.localVariables[0])
		case 27: // iload_1
			frame.stack = append(frame.stack, frame.localVariables[1])
		case 42: //	aload0
			frame.stack = append(frame.stack, frame.localVariables[0])
		case 43: // aload1
			frame.stack = append(frame.stack, frame.localVariables[1])
		case 60: // istore_1
			last := frame.stack[len(frame.stack)-1].(int)
			frame.stack = frame.stack[:len(frame.stack)-1]
			frame.localVariables[1] = last
		case 61: // istore2
			last := frame.stack[len(frame.stack)-1].(int)
			frame.stack = frame.stack[:len(frame.stack)-1]
			frame.localVariables[2] = last
		case 76: // astore_1
			last := frame.stack[len(frame.stack)-1]
			frame.stack = frame.stack[:len(frame.stack)-1]
			frame.localVariables[1] = last
		case 87: // pop
			frame.stack = frame.stack[:len(frame.stack)-1]
		case 89: // dup
			last := frame.stack[len(frame.stack)-1]
			frame.stack = append(frame.stack, last)
		case 96: // iadd
			first := frame.stack[len(frame.stack)-1].(int)
			second := frame.stack[len(frame.stack)-2].(int)
			frame.stack = frame.stack[:len(frame.stack)-2]
			frame.stack = append(frame.stack, first+second)
		case 104: // imul
			first := frame.pop().(int)
			second := frame.pop().(int)
			frame.stack = append(frame.stack, first*second)
		case 172: // ireturn
			return frame.pop()
		case 177: // return
			frame.stack = append(frame.stack, nil)
			return nil
		case 180: // getField
			frame.instructionPointer++
			indexByte1 := frame.code[frame.instructionPointer]
			frame.instructionPointer++
			indexByte2 := frame.code[frame.instructionPointer]
			_ = (indexByte1 << 8) | indexByte2

			objectRef := frame.pop().(*Class)
			frame.push(objectRef.fields[0].value)
			//frame.stack = frame.stack[:len(frame.stack)-1]
			//frame.stack = append(frame.stack, frame.class.constantPool.resolveFieldRef(u2(fieldRefIndex)))
		case 181: // putField
			frame.instructionPointer++
			indexByte1 := frame.code[frame.instructionPointer]
			frame.instructionPointer++
			indexByte2 := frame.code[frame.instructionPointer]
			_ = (indexByte1 << 8) | indexByte2

			value := frame.pop().(int)
			objectRef := frame.pop().(*Class)

			objectRef.fields[0].value = value
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

		case 184: // (0xb8)    invokestatic
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
			frame.instructionPointer++
			indexByte1 := frame.code[frame.instructionPointer]
			frame.instructionPointer++
			indexByte2 := frame.code[frame.instructionPointer]
			classIndex := (indexByte1 << 8) | indexByte2

			className := frame.class.constantPool.resolveString(frame.class.constantPool[classIndex-1].info.ClassIndex)
			fmt.Printf("className %s\n", className)
			file, _ := os.Open("Main.class")
			loader := &ClassLoader{reader: file}

			newClass := loader.loadClass()
			frame.stack = append(frame.stack, &newClass)

		default:
			fmt.Printf("wooops unimplemented op code %d\n", op)
		}

		frame.instructionPointer++

	}

	return nil
}
