package netcdf4

import (
	"fmt"
)

type Var struct {
	nullObject bool

	myId ID

	groupId ID
}

func NewVarNull() (v Var) {
	v.nullObject = true
	v.myId = -1
	v.groupId = -1
	return
}

func NewVar(group Group, vID ID) (v Var) {
	v.nullObject = false
	v.myId = vID
	v.groupId, _ = group.ID()
	return
}

// Gets parent group.
func (v Var) GetParentGroup() Group {
	return NewGroup(v.groupId)
}

// Get the variable id.
func (v Var) GetId() ID {
	return v.myId
}

func (v Var) GetGrpId() ID {
	return v.groupId
}

func (v Var) DataLength() (int, error) {
	ncDims, err := v.GetDims()
	if err != nil {
		return 0, err
	}

	if len(ncDims) == 0 { //scalar
		return 1, nil
	}
	n := 1
	for _, dim := range ncDims {
		dimLen, err := dim.GetSize() //consider the unlimited dims
		if err != nil {
			return 0, err
		}
		n *= dimLen
	}
	return n, nil
}

///////////////////////////////////////////
//  Information about the variable type
///////////////////////////////////////////

// Gets the NcxxType object with a given name.
func (v Var) GetType() (Type, error) {

	// if this variable has not been defined, return a NULL type
	if v.IsNull() {
		return NewTypeNull(), fmt.Errorf("getType NULL")
	}

	// first get the typeid
	xtypep, err := NcInqVartype(v.groupId, v.myId)
	if err != nil {
		return NewTypeNull(), err
	}

	if xtypep == Byte.GetId() {
		return Byte, nil
	}
	if xtypep == Ubyte.GetId() {
		return Ubyte, nil
	}
	if xtypep == Char.GetId() {
		return Char, nil
	}
	if xtypep == Short.GetId() {
		return Short, nil
	}
	if xtypep == Ushort.GetId() {
		return Ushort, nil
	}
	if xtypep == Int.GetId() {
		return Int, nil
	}
	if xtypep == Uint.GetId() {
		return Uint, nil
	}
	if xtypep == Int64.GetId() {
		return Int64, nil
	}
	if xtypep == Uint64.GetId() {
		return Uint64, nil
	}
	if xtypep == Float.GetId() {
		return Float, nil
	}
	if xtypep == Double.GetId() {
		return Double, nil
	}
	if xtypep == String.GetId() {
		return String, nil
	}

	//multimap<string,NcxxType>::const_iterator it;
	//multimap<string,NcxxType>
	//types(NcxxGroup(groupId).getTypes(NcxxGroup::ParentsAndCurrent));
	//for(it=types.begin(); it!=types.end(); it++) {
	//if(it->second.getId() == xtypep) return it->second;
	//}
	// we will never reach here
	return NewTypeNull(), nil
}

// Gets the set of Ncdim objects.
func (v Var) GetDims() ([]Dim, error) {

	dimCount, dimIds, err := NcInqVardimid(v.groupId, v.myId)
	if err != nil {
		return []Dim(nil), err
	}

	ncDims := make([]Dim, dimCount)
	for i := 0; i < dimCount; i++ {
		ncDims[i] = NewDim(v.GetParentGroup(), dimIds[i])
	}

	return ncDims, nil
}

// Gets the i'th Dim object.
func (v Var) GetDim(i int) (Dim, error) {
	ncDims, err := v.GetDims()
	if err != nil {
		return NewDimNull(), err
	}
	if i >= len(ncDims) || i < 0 {
		return NewDimNull(), fmt.Errorf("error: index out of range: index = %d, size = %d", i, len(ncDims))
	}
	return ncDims[i], nil
}

// Gets the number of dimensions.
func (v Var) GetDimCount() (int, error) {
	// get the number of dimensions
	return NcInqVarndims(v.groupId, v.myId)
}

/*! Returns true if this object is null (i.e. it has no contents); otherwise returns false. */
func (v Var) IsNull() bool {
	return v.nullObject
}

///////////////////////////////////
// Other Basic variable info
///////////////////////////////////

// The name of this variable.
func (v Var) GetName() (string, error) {
	return NcInqVarname(v.groupId, v.myId)
}

///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////
//  data writing
///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// Write a scalar into the netCDF variable.

// Write string scalar

//func (v Var) putStringScalar(dataVal string) {
//vector<size_t> index;
//index.push_back(0);
//putVal(index, dataVal);
//}

//// Write char scalar
//
//void NcxxVar::putVal(char dataVal) const {
//putVal(&dataVal);
//}

// There are four kinds of writing
// nc_put_var1_xxx  >>write one datum
// nc_put_var_xxx   >>Write an entire variable with one call.
// nc_put_vara_xxx   >>Write an array of values to a variable.
// nc_put_varm_xxx   >>Write a mapped array of values to a variable.

// The problem is the memory layout of the data in C and go are different, thus, it is difficult to use the C API in go
// Write the entire data into the netCDF variable.

func (v Var) PutValAll(data interface{}) {
	//CheckDataMode(groupId);
	//varType, _ := v.GetType()
	//if varType.IsComplex() {
	//	nc_put_var(groupId, myId, dataValues)
	//} else {
	NcPutVarDouble(v.groupId, v.myId, data)
	//}
}

//func  (v Var) checkData(data interface{}) error {
//	// check the length
//	if reflect.TypeOf(data).Kind()==reflect.Slice{
//		if reflect.ValueOf(data).Len()!=
//	}
//
//}

// Data reading

// Reads the entire data of the netCDF variable.

// The name of this variable.
//func (v Var) GetVal() (interface{}, error) {
//
//}
