package main

import (
	"fmt"
	"os"
	"strings"
)

type Frame struct {
	operandStack       []interface{}
	localVariables     []interface{} // function parameters basically
	code               []byte
	instructionPointer uint32
	class              Class
}

func (frame *Frame) Exec() interface{} {

	for {
		op := frame.code[frame.instructionPointer]

		switch op {
		case 0: // noop

		case 4: // iconst_1
			frame.operandStack = append(frame.operandStack, 1)
		case 5: // iconst_2
			frame.operandStack = append(frame.operandStack, 2)
		case 6: // iconst_3
			frame.operandStack = append(frame.operandStack, 3)
		case 7: // iconst_4
			frame.operandStack = append(frame.operandStack, 4)
		case 17: // sipush
		case 26: // iload_0
			frame.operandStack = append(frame.operandStack, frame.localVariables[0])
		case 27: // iload_1
			frame.operandStack = append(frame.operandStack, frame.localVariables[1])
		case 42: //	aload0
			frame.operandStack = append(frame.operandStack, frame.localVariables[0])
		case 43: // aload1
			frame.operandStack = append(frame.operandStack, frame.localVariables[1])
		case 60: // istore_1
			last := frame.operandStack[len(frame.operandStack)-1].(int)
			frame.operandStack = frame.operandStack[:len(frame.operandStack)-1]
			frame.localVariables[1] = last
		case 61: // istore2
			last := frame.operandStack[len(frame.operandStack)-1].(int)
			frame.operandStack = frame.operandStack[:len(frame.operandStack)-1]
			frame.localVariables[2] = last
		case 76: // astore_1
			last := frame.operandStack[len(frame.operandStack)-1]
			frame.operandStack = frame.operandStack[:len(frame.operandStack)-1]
			frame.localVariables[1] = last
		case 87: // pop
			frame.operandStack = frame.operandStack[:len(frame.operandStack)-1]
		case 89: // dup
			last := frame.operandStack[len(frame.operandStack)-1]
			frame.operandStack = append(frame.operandStack, last)
		case 96: // iadd
			first := frame.operandStack[len(frame.operandStack)-1].(int)
			second := frame.operandStack[len(frame.operandStack)-2].(int)
			frame.operandStack = frame.operandStack[:len(frame.operandStack)-2]
			frame.operandStack = append(frame.operandStack, first+second)
		case 104: // imul
			first := frame.pop().(int)
			second := frame.pop().(int)
			frame.operandStack = append(frame.operandStack, first*second)
		case 172: // ireturn
			return nil
		case 177: // return
			frame.operandStack = append(frame.operandStack, nil)
			return nil
		case 180: // getField
			frame.instructionPointer++
			indexByte1 := frame.code[frame.instructionPointer]
			frame.instructionPointer++
			indexByte2 := frame.code[frame.instructionPointer]
			_ = (indexByte1 << 8) | indexByte2

			objectRef := frame.pop().(Class)
			frame.push(objectRef.fields[0].value)
			//frame.operandStack = frame.operandStack[:len(frame.operandStack)-1]
			//frame.operandStack = append(frame.operandStack, frame.class.constantPool.resolveFieldRef(u2(fieldRefIndex)))
		case 181: // putField
			frame.instructionPointer++
			indexByte1 := frame.code[frame.instructionPointer]
			frame.instructionPointer++
			indexByte2 := frame.code[frame.instructionPointer]
			_ = (indexByte1 << 8) | indexByte2

			value := frame.pop().(int)
			objectRef := frame.pop().(Class)

			objectRef.fields[0].value = value
		case 182, 183: // invokespecial
			frame.instructionPointer++
			indexByte1 := frame.code[frame.instructionPointer]
			frame.instructionPointer++
			indexByte2 := frame.code[frame.instructionPointer]
			methodRefIndex := (indexByte1 << 8) | indexByte2

			resolved := frame.class.constantPool.resolveMethodRef(u2(methodRefIndex))

			// now we need to execute this again
			var method Method
			if resolved.class == frame.class.name {
				method = frame.class.findMethod(resolved.name)

				start := strings.Index(resolved.methodType, "(")
				end := strings.Index(resolved.methodType, ")")
				params := resolved.methodType[start+1 : end]
				returnTypes := resolved.methodType[end+1:]

				methodCodeAttribute := method.findCodeAttribute()
				args := make([]interface{}, 0)
				for _ = range len(params) {
					arg := frame.operandStack[len(frame.operandStack)-1]
					frame.operandStack = frame.operandStack[:len(frame.operandStack)-1]
					args = append(args, arg)
				}
				objectRef := frame.operandStack[len(frame.operandStack)-1]
				frame.operandStack = frame.operandStack[:len(frame.operandStack)-1]
				args = append([]interface{}{objectRef.(Class)}, args...)

				methodFrame := methodCodeAttribute.toFrame(objectRef.(Class), args...)
				methodFrame.Exec()

				if len(returnTypes) != 0 && returnTypes != "V" {
					frame.operandStack = append(frame.operandStack, methodFrame.operandStack[len(methodFrame.operandStack)-1])
				}

			} else if resolved.class == "java/lang/Object" {
			} else {
				panic("we do not know this class")
			}

			fmt.Printf("%d", methodRefIndex)

		case 184: // (0xb8)    invokestatic
			frame.instructionPointer++
			indexByte1 := frame.code[frame.instructionPointer]
			frame.instructionPointer++
			indexByte2 := frame.code[frame.instructionPointer]
			methodRefIndex := (indexByte1 << 8) | indexByte2

			resolved := frame.class.constantPool.resolveMethodRef(u2(methodRefIndex))

			// now we need to execute this again
			method := frame.class.findMethod(resolved.name)

			// the variables are already on the stack, does the internal stack get passed along????

			start := strings.Index(resolved.methodType, "(")
			end := strings.Index(resolved.methodType, ")")
			params := resolved.methodType[start+1 : end]
			returnTypes := resolved.methodType[end+1:]

			methodCodeAttribute := method.findCodeAttribute()
			args := make([]interface{}, 0)
			for _ = range len(params) {
				arg := frame.operandStack[len(frame.operandStack)-1]
				frame.operandStack = frame.operandStack[:len(frame.operandStack)-1]
				args = append(args, arg)
			}
			methodFrame := methodCodeAttribute.toFrame(frame.class, args...)
			methodFrame.Exec()

			if len(returnTypes) != 0 && returnTypes != "V" {
				frame.operandStack = append(frame.operandStack, methodFrame.operandStack[len(methodFrame.operandStack)-1])
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
			frame.operandStack = append(frame.operandStack, newClass)

		default:
			fmt.Printf("wooops unimplemented op code %d\n", op)
		}

		frame.instructionPointer++

	}

	return nil
}
