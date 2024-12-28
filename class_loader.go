package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

type ClassLoader struct {
	reader io.Reader
	err    error
}

func (classLoader *ClassLoader) loadClass() Class {
	class := Class{
		magic:        classLoader.loadU4(),
		minorVersion: classLoader.loadU2(),
		majorVersion: classLoader.loadU2(),
		constantPool: classLoader.loadConstantPool(),
		accessFlags:  classLoader.loadU2(),
	}
	class.name = class.constantPool.resolveClass(classLoader.loadU2())
	class.superClass = class.constantPool.resolveClass(classLoader.loadU2())
	class.interfaces = classLoader.loadInterfaces(class.constantPool)
	class.fields = classLoader.loadFields(class.constantPool)
	fields := classLoader.loadFields(class.constantPool)
	methods := make([]Method, len(fields))
	for i, field := range fields {
		methods[i] = Method(field)
	}
	class.methods = methods

	class.attributes = classLoader.loadAttributes(class.constantPool)

	return class
}

func (classLoader *ClassLoader) loadFields(constantPool ConstantPool) (fields []Field) {
	fieldCount := classLoader.loadU2()

	for i := u2(0); i < fieldCount; i++ {
		fields = append(fields, Field{
			accessFlags: classLoader.loadU2(),
			name:        constantPool.resolveString(classLoader.loadU2()),
			descriptor:  constantPool.resolveString(classLoader.loadU2()),
			attributes:  classLoader.loadAttributes(constantPool),
		})
	}

	return fields
}

func (classLoader *ClassLoader) loadAttributes(constantPool ConstantPool) (attributes []Attribute) {
	attributesCount := classLoader.loadU2()

	for i := u2(0); i < attributesCount; i++ {
		attributes = append(attributes, Attribute{
			name: constantPool.resolveString(classLoader.loadU2()),
			info: classLoader.readBytes(int(classLoader.loadU4())),
		})
	}

	return attributes
}

func (classLoader *ClassLoader) loadInterfaces(constantPool ConstantPool) (interfaces []Interface) {
	interfaceCount := classLoader.loadU2()

	for i := u2(0); i < interfaceCount; i++ {
		interfaces = append(interfaces, Interface{name: constantPool.resolveString(classLoader.loadU2())})
	}

	return interfaces
}

func (classLoader *ClassLoader) loadConstantPool() (constantPool ConstantPool) {
	constantPoolCount := classLoader.loadU2()

	//The constant_pool table is indexed from 1 to constant_pool_count - 1.
	for i := u2(1); i < constantPoolCount; i++ {
		constant := Constant{tag: ConstantPoolTags(classLoader.loadU1())}

		switch constant.tag {

		case CONSTANT_Integer:
			constant.info = ConstantType{Integer: int(classLoader.loadU4())}
		case CONSTANT_String:
			utfLength := classLoader.loadU2()
			constant.info = ConstantType{String: string(classLoader.readBytes(int(utfLength)))}
		case CONSTANT_StringIndex:
			constant.info = ConstantType{StringIndex: classLoader.loadU2()}
		case CONSTANT_Class:
			constant.info = ConstantType{ClassIndex: classLoader.loadU2()}
		case CONSTANT_Methodref:
			methodRef := MethodRef{
				ClassIndex:       classLoader.loadU2(),
				NameAndTypeIndex: classLoader.loadU2(),
			}
			constant.info = ConstantType{MethodRef: methodRef}
		case CONSTANT_Fieldref:
			fieldRef := FieldRef{
				ClassIndex:       classLoader.loadU2(),
				NameAndTypeIndex: classLoader.loadU2(),
			}
			constant.info = ConstantType{FieldRef: fieldRef}
		case CONSTANT_NameAndType:
			nameAndType := NameAndType{
				nameIndex:  classLoader.loadU2(),
				descriptor: classLoader.loadU2(),
			}
			constant.info = ConstantType{NameAndType: nameAndType}
		default:
			classLoader.err = fmt.Errorf("unsupported tag: %d", constant.tag)
		}
		constantPool = append(constantPool, constant)
	}
	return constantPool
}

func (classLoader *ClassLoader) readBytes(n int) []byte {
	byteArray := make([]byte, n, n)
	// only store the first error
	if classLoader.err == nil {
		_, classLoader.err = io.ReadFull(classLoader.reader, byteArray)
	}
	return byteArray
}

func (classLoader *ClassLoader) loadU1() u1 { return classLoader.readBytes(1)[0] }

func (classLoader *ClassLoader) loadU2() u2 {
	return binary.BigEndian.Uint16(classLoader.readBytes(2))
}
func (classLoader *ClassLoader) loadU4() u4 {
	return binary.BigEndian.Uint32(classLoader.readBytes(4))
}
func (classLoader *ClassLoader) loadU8() u8 {
	return binary.BigEndian.Uint64(classLoader.readBytes(8))
}
