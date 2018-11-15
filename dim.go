package netcdf4

//Dim is a representative of a Dimension
type Dim struct {
	nullObject bool
	id, group  ID
}

// NewDimNull reutnrs a new dimension where it is configured to be a nul dimension
func NewDimNull() (dim Dim) {
	return Dim{
		nullObject: true,
		id:         -1,
		group:      -1,
	}
	return
}

// NewDim create a new dimension attached to the passed group
func NewDim(group Group, dimID ID) (dim Dim) {
	gid, _ := group.ID()
	return Dim{
		nullObject: false,
		id:         dimID,
		group:      gid,
	}
}

// GetSize gets the size of the dimension, for unlimited, this is the current number of records.
func (dim Dim) GetSize() (lenDim int, err error) {
	cLenDim, err := NcInqDimLen(dim.group, dim.id)
	lenDim = int(cLenDim)
	return
}

// GetParentGroup gets a NcxxGroup object of the parent group.
func (dim Dim) GetParentGroup() Group {
	return NewGroup(dim.group)
}

// returns true if this dimension is unlimited.
//func (dim Dim) IsUnlimited() bool {
//
//	ncInqUnlimdims(ncId ID) (nunlimdimsp, unlimdimidsp int, err error)
//
//int numlimdims;
//int* unlimdimidsp=NULL;
//// get the number of unlimited dimensions
//ncxxCheck(nc_inq_unlimdims(group,&numlimdims,unlimdimidsp),__FILE__,__LINE__);
//if (numlimdims){
//// get all the unlimited dimension ids in this group
//vector<int> unlimdimid(numlimdims);
//ncxxCheck(nc_inq_unlimdims(group,&numlimdims,&unlimdimid[0]),__FILE__,__LINE__);
//vector<int>::iterator it;
//// now look to see if this dimension is unlimited
//it = find(unlimdimid.begin(),unlimdimid.end(),id);
//return it != unlimdimid.end();
//}
//return false;
//}

/*ID returns the he netCDF Id of this dimension. */
func (dim Dim) ID() ID {
	return dim.id
}

//Group return the group ID
func (dim Dim) Group() ID {
	return dim.group
}

// Name returns the name of the dimension.
func (dim Dim) Name() (name string, err error) {
	name, err = NcInqDimname(dim.group, dim.id)
	return
}

// RenameTo attempts to rename the dimension to name
func (dim Dim) RenameTo(name string) (err error) {
	err = NcRenameDim(dim.group, dim.id, name)
	return
}

// IsNull returns true if the object is null. Returns true if this object is null (i.e. it has no content); otherwise returns false.
func (dim Dim) IsNull() bool {
	return dim.nullObject
}

//SetNull forcibly sets the dimension to null
func (dim *Dim) SetNull() {
	dim.nullObject = true
}

// IsValidDim returns true if the dimension is valid ??
func (dim Dim) IsValidDim(group Group) (bool, error) {
	gid, _ := group.ID()
	nDims, dimIds, err := NcInqDimids(gid, true)
	if err != nil {
		return false, err
	}
	for i := 0; i < nDims; i++ {
		if dim.id == dimIds[i] {
			return true, nil //note where two dims can have the same dimid
		}
	}

	//if dim.GetGrpId() == grpId{
	//	return true, nil
	//}
	//groupsM, err := group.GetGroupsM(ParentsGrps)
	//if err!=nil{
	//	return false, err
	//}
	//
	//_, grps :=groupsM.GetAllPair()
	//for _, grp := range grps{
	//	grpIdC, _ := grp.GetId()
	//	if dim.GetGrpId() == grpIdC{
	//		return true, nil
	//	}
	//}
	return false, nil
}
