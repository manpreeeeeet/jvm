package main

type ConstantPoolTags u1

const (
	// UTF_8 in the JVM specification
	CONSTANT_String ConstantPoolTags = iota + 1
	_
	CONSTANT_Integer
	CONSTANT_Float
	CONSTANT_Long
	CONSTANT_Double
	CONSTANT_Class              = 7
	CONSTANT_StringIndex        = 8
	CONSTANT_Fieldref           = 9
	CONSTANT_Methodref          = 10
	CONSTANT_InterfaceMethodref = 11
	CONSTANT_NameAndType        = 12
	_                           = 13
	_
	CONSTANT_MethodHandle  = 15
	CONSTANT_MethodType    = 16
	CONSTANT_Dynamic       = 17
	CONSTANT_InvokeDynamic = 18
	CONSTANT_Module        = 19
	CONSTANT_Package       = 20
)

type MethodRef struct {
	ClassIndex       u2
	NameAndTypeIndex u2
}

type MethodRefResolved struct {
	class      string
	name       string
	descriptor string
}

type FieldRef = MethodRef

type FieldRefResolved = MethodRefResolved

type NameAndType struct {
	nameIndex  u2
	descriptor u2
}

type ConstantType struct {
	String      string
	StringIndex u2
	ClassIndex  u2
	Integer     int
	MethodRef   MethodRef
	FieldRef    FieldRef
	NameAndType NameAndType
}
type ConstantPool []Constant

type Constant struct {
	tag  ConstantPoolTags
	info ConstantType
}

func (constantPool ConstantPool) resolveString(index u2) string {
	if constantPool[index-1].tag == CONSTANT_String {
		return constantPool[index-1].info.String
	}
	return ""
}

func (constantPool ConstantPool) resolveFieldRef(index u2) FieldRefResolved {
	if constantPool[index-1].tag == CONSTANT_Methodref {
		methodRef := constantPool[index-1].info.MethodRef

		methodRefResolved := MethodRefResolved{class: constantPool.resolveClass(methodRef.ClassIndex)}

		if constantPool[methodRef.NameAndTypeIndex-1].tag != CONSTANT_NameAndType {
			panic("wrong name and type index")
		}

		nameAndType := constantPool[methodRef.NameAndTypeIndex-1].info.NameAndType
		methodRefResolved.name = constantPool.resolveString(nameAndType.nameIndex)
		methodRefResolved.descriptor = constantPool.resolveString(nameAndType.descriptor)

		return methodRefResolved

	}
	panic("yooo this isn't a field ref")
}

func (constantPool ConstantPool) resolveClass(index u2) string {
	if constantPool[index-1].tag == CONSTANT_Class {
		return constantPool.resolveString(constantPool[index-1].info.ClassIndex)
	}
	panic("not a class index")
}

func (constantPool ConstantPool) resolveMethodRef(index u2) MethodRefResolved {
	if constantPool[index-1].tag == CONSTANT_Methodref {
		methodRef := constantPool[index-1].info.MethodRef

		methodRefResolved := MethodRefResolved{class: constantPool.resolveClass(methodRef.ClassIndex)}

		if constantPool[methodRef.NameAndTypeIndex-1].tag != CONSTANT_NameAndType {
			panic("wrong name and type index")
		}

		nameAndType := constantPool[methodRef.NameAndTypeIndex-1].info.NameAndType
		methodRefResolved.name = constantPool.resolveString(nameAndType.nameIndex)
		methodRefResolved.descriptor = constantPool.resolveString(nameAndType.descriptor)

		return methodRefResolved

	}
	panic("yooo this isn't a method ref")
}
