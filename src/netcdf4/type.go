package netcdf4

// #include <stdlib.h>
// #include <netcdf.h>
import "C"

type NcType C.nc_type

type Type struct {
	nullObject bool
	myId       NcType //the type Id
	groupId    ID     //the group Id

	/*! An ncid associated with a particular open file (returned from nc_open).
	  This is required by many of the functions ncType uses, such as nc_inq_type */
	gFileId ID
}

func NewTypeNull() (t Type) {
	t.nullObject = true
	t.myId = -1
	t.groupId = -1
	t.gFileId = -1
	return
}
func NewType(id NcType) (t Type) {
	t.nullObject = false
	t.myId = id
	t.groupId = 0
	t.gFileId = -1
	return
}


var Byte = NewType(C.NC_BYTE)
var Ubyte = NewType(C.NC_UBYTE)
var Char = NewType(C.NC_CHAR)
var Short = NewType(C.NC_SHORT)
var Ushort = NewType(C.NC_USHORT)
var Int = NewType(C.NC_INT)
var Uint = NewType(C.NC_UINT)
var Int64 = NewType(C.NC_INT64)
var Uint64 = NewType(C.NC_UINT64)
var Float = NewType(C.NC_FLOAT)
var Double = NewType(C.NC_DOUBLE)
var String = NewType(C.NC_STRING)

/*! Returns true if this object is null (i.e. it has no contents); otherwise returns false. */
func (t Type) IsNull() bool {
	return t.nullObject
}