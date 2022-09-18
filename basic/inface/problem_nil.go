package inface

import (
	"fmt"
	"unsafe"

	"github.com/CoderBenson/go_study/basic/stru"
)

type tflag uint8

type nameOff int32
type typeOff int32
type textOff int32

type name struct {
	bytes *byte
}

type _type struct {
	size       uintptr
	ptrdata    uintptr // size of memory prefix holding all pointers
	hash       uint32
	tflag      tflag
	align      uint8
	fieldAlign uint8
	kind       uint8
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal func(unsafe.Pointer, unsafe.Pointer) bool
	// gcdata stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, gcdata is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	gcdata    *byte
	str       nameOff
	ptrToThis typeOff
}

type interfacetype struct {
	typ     _type
	pkgpath name
	mhdr    []imethod
}

type imethod struct {
	name nameOff
	ityp typeOff
}

type itab struct {
	inter *interfacetype
	_type *_type
	hash  uint32 // copy of _type.hash. Used for type switches.
	_     [4]byte
	fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}

type iface struct {
	tab  *itab
	data unsafe.Pointer
}

type eface struct {
	_type *_type
	data  unsafe.Pointer
}

func TestNil() *stru.Student {
	// var result Eatable
	if 2 == 3 {
		return nil
	}
	return nil
}

func TestInterface() {
	var e Eatable
	e = TestNil() // stru.NewStudent(20, "xiaoming", true, 99)
	fmt.Println(e == nil)
	// e := TestNil()
	face := (*iface)(unsafe.Pointer(&e))
	fmt.Println(*face)
}
