package netcdf4
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
	v.groupId, _ = group.GetId()
	return
}

/*! Returns true if this object is null (i.e. it has no contents); otherwise returns false. */
func (v Var) IsNull() bool {
	return v.nullObject
}