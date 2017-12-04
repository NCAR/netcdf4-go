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
func (dim Dim) GetSize()(lenDim int, err error){
	cLenDim, err := NcInqDimLen(dim.groupId, dim.myId)
	lenDim = int(cLenDim)
	return
}
//
//
//// returns true if this dimension is unlimited.
//bool NcxxDim::isUnlimited() const
//{
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

// gets the name of the dimension.

func (dim Dim) GetName()(name string, err error){
	name, err = NcInqDimname(dim.groupId, dim.myId)
	return
}

// renames this dimension.

func (dim Dim) ReName(name string)( err error){
	err = NcRenameDim(dim.groupId, dim.myId, name)
	return
}
