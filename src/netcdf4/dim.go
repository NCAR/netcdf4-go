package netcdf4

type Dim struct {
	nullObject bool

	myId ID

	groupId ID
}

func NewDimNull() (dim Dim) {
	dim.nullObject = true
	dim.myId = -1
	dim.groupId = -1
	return
}

func NewDim(group Group, dimID ID) (dim Dim) {
	dim.nullObject = false
	dim.myId = dimID
	dim.groupId, _ = group.GetId()
	return
}

//// gets the size of the dimension, for unlimited, this is the current number of records.
func (dim Dim) GetSize() (lenDim int, err error) {
	cLenDim, err := NcInqDimLen(dim.groupId, dim.myId)
	lenDim = int(cLenDim)
	return
}

/*! Gets a  NcxxGroup object of the parent group. */

func (dim Dim) GetParentGroup() Group {
	return NewGroup(dim.groupId)
}

// returns true if this dimension is unlimited.
//func (dim Dim) IsUnlimited() bool {
//
//	ncInqUnlimdims(ncId ID) (nunlimdimsp, unlimdimidsp int, err error)
//
//int numlimdims;
//int* unlimdimidsp=NULL;
//// get the number of unlimited dimensions
//ncxxCheck(nc_inq_unlimdims(groupId,&numlimdims,unlimdimidsp),__FILE__,__LINE__);
//if (numlimdims){
//// get all the unlimited dimension ids in this group
//vector<int> unlimdimid(numlimdims);
//ncxxCheck(nc_inq_unlimdims(groupId,&numlimdims,&unlimdimid[0]),__FILE__,__LINE__);
//vector<int>::iterator it;
//// now look to see if this dimension is unlimited
//it = find(unlimdimid.begin(),unlimdimid.end(),myId);
//return it != unlimdimid.end();
//}
//return false;
//}

/*! The netCDF Id of this dimension. */
func (dim Dim) GetId() ID {
	return dim.myId
}
func (dim Dim) GetGrpId() ID {
	return dim.groupId
}
// gets the name of the dimension.

func (dim Dim) GetName() (name string, err error) {
	name, err = NcInqDimname(dim.groupId, dim.myId)
	return
}

// renames this dimension.

func (dim Dim) ReName(name string) (err error) {
	err = NcRenameDim(dim.groupId, dim.myId, name)
	return
}

/*! Returns true if this object is null (i.e. it has no contents); otherwise returns false. */

func (dim Dim) IsNull() bool {
	return dim.nullObject
}

// set to null

func (dim *Dim) SetNull() {
	dim.nullObject = true
}


// find the

func (dim Dim) IsValidDim(group Group) (bool, error) {
	grpId, _ := group.GetId()
	nDims , dimIds, err := NcInqDimids(grpId,true)
	if err!=nil{
		return false, err
	}
	for i:=0; i<nDims;i++  {
		if dim.myId==dimIds[i]{
			return true,nil  //note where two dims can have the same dimid
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