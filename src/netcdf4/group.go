package netcdf4

import (
	"fmt"
	"strconv"
)

//////////////////////////////////////////////////////////////////////
//  Netcdf4 support for Go
//
//  Copied from code for c++ by:
//
//    Mike Dixon, EOL, NCAR
//    P.O.Box 3000, Boulder, CO, 80307-3000, USA
//
//  Offical NetCDF codebase is at:
//
//    https://github.com/Unidata/netcdf-cxx4
//
//  Modification for Go made by:
//
//    Hao Huang, Nanjing University
//    Email: hhuangwx@gmail.com
//
//////////////////////////////////////////////////////////////////////

type Group struct {
	nullObject bool

	myId ID

	// option to use the 'proposed_standard_name' attribute instead
	// of 'standard_name'.

	useProposedStandardName bool
}

/*!
  The enumeration list contains the options for selecting groups (used for returned set of Group objects).
*/
type GroupLocation int

const (
	ChildrenGrps           GroupLocation = iota //!< Select from the set of children in the current group.
	ParentsGrps                                 //!< Select from set of parent groups (excludes the current group).
	ChildrenOfChildrenGrps                      //!< Select from set of all children of children in the current group.
	AllChildrenGrps                             //!< Select from set of all children of the current group and beneath.
	ParentsAndCurrentGrps                       //!< Select from set of parent groups(includes the current group).
	AllGrps                                     //!< Select from set of parent groups, current groups and all the children beneath.
)

const groupLocationName = "ChildrenGrpsParentsGrpsChildrenOfChildrenGrpsAllChildrenGrpsParentsAndCurrentGrpsAllGrps"

var groupLocationIndex = [...]uint8{0, 12, 23, 45, 60, 81, 88}

func (i GroupLocation) String() string {
	if i < 0 || i >= GroupLocation(len(groupLocationIndex)-1) {
		return "GroupLocation(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return groupLocationName[groupLocationIndex[i]:groupLocationIndex[i+1]]
}

type Location int

const (
	Current            Location = iota //!< Select from contents of current group.
	Parents                            //!< Select from contents of parents groups.
	Children                           //!< Select from contents of children groups.
	ParentsAndCurrent                  //!< Select from contents of current and parents groups.
	ChildrenAndCurrent                 //!< Select from contents of current and child groups.
	All                                //!< Select from contents of current, parents and child groups.
)

const locationName = "CurrentParentsChildrenParentsAndCurrentChildrenAndCurrentAll"

var locationIndex = [...]uint8{0, 7, 14, 22, 39, 57, 60}

func (i Location) String() string {
	if i < 0 || i >= Location(len(locationIndex)-1) {
		return "Location(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return locationName[locationIndex[i]:locationIndex[i+1]]
}

func NewGroupNull() (group Group) {
	group.nullObject = true
	group.myId = -1
	group.useProposedStandardName = false
	return
}

func NewGroup(groupId ID) (group Group) {
	group.nullObject = false
	group.myId = groupId
	group.useProposedStandardName = false
	return
}

func NewGroupFrom(rhs Group) (group Group) {
	group.nullObject = rhs.nullObject
	group.myId = rhs.myId
	group.useProposedStandardName = rhs.useProposedStandardName
	return
}

func (group Group) NewGroupWith(groupId ID) {
	group.nullObject = false
	group.myId = groupId
	group.useProposedStandardName = false
	return
}

// /////////////
// NcxxGroup-related methods
// /////////////

// Get the group name.
func (group Group) GetName(fullName bool /*false*/) (string, error) {
	if group.IsNull() {
		return "", fmt.Errorf("error: attempt to invoke GetName on a Null group")
	}
	var groupName string
	var err error
	if fullName {
		// return full name of group with foward "/" separarating sub-groups.
		if groupName, err = ncInqGrpnameFull(group.myId); err != nil {
			return groupName, err
		}
	} else {
		// return the (local) name of this group.
		if groupName, err = ncInqGrpname(group.myId); err != nil {
			return groupName, err
		}
	}
	return groupName, nil
}

// returns true if this is the root group.
func (group Group) IsRootGroup() (bool, error) {
	grpName, err := group.GetName(false)
	if err != nil {
		return false, err
	}
	if grpName == "/" {
		return true, nil
	} else {
		return false, nil
	}
}

// Get the parent group.
func (group Group) GetParentGroup() (Group, error) {
	if group.IsNull() {
		return NewGroupNull(), fmt.Errorf("error: attempt to invoke GetParentGroup on a Null group")
	}

	if parentId, err := ncInqGrpParent(group.myId); err != nil {
		return NewGroupNull(), nil //if no parent id is found, return null group
	} else {
		return NewGroup(parentId), nil
	}
}

// Get the group id.
func (group Group) GetId() (ID, error) {
	if group.IsNull() {
		return ID(-1), fmt.Errorf("error: attempt to invoke GetId on a Null group")
	}
	return group.myId, nil
}

// Get the number of Group objects.
func (group Group) GetGroupCount(location GroupLocation /*ChildrenGrps*/) (int, error) {

	if group.IsNull() {
		return -1, fmt.Errorf("error: attempt to invoke GetGroupCount on a Null group")
	}
	myId, _ := group.GetId()
	// initialize group counter
	nGroups := 0

	// record this group
	if location == ParentsAndCurrentGrps || location == AllGrps {
		nGroups ++
	}

	// number of children in current group
	if location == ChildrenGrps || location == AllChildrenGrps || location == AllGrps {
		if numGrps, _, err := NcInqGrps(myId); err != nil {
			return -1, err
		} else {
			nGroups += numGrps
		}
	}

	// search in parent groups
	if location == ParentsGrps || location == ParentsAndCurrentGrps || location == AllGrps {
		groups, err := group.GetGroupsM(ParentsGrps)
		if err != nil {
			return -1, err
		}
		nGroups += len(groups)
	}
	// get the number of all children that are childreof children
	if location == ChildrenOfChildrenGrps || location == AllChildrenGrps || location == AllGrps {
		groups, err := group.GetGroupsM(ChildrenOfChildrenGrps)
		if err != nil {
			return -1, err
		}
		nGroups += len(groups)
	}

	return nGroups, nil
}

// Get the set of child NcxxGroup objects.

func (group Group) GetGroupsM(location GroupLocation) (Multimap, error) {
	ncGroups := NewMultimap()

	if group.IsNull() {
		return ncGroups, fmt.Errorf("error: attempt to invoke GetGroupsM on a Null group")
	}
	myId, _ := group.GetId()

	// record this group
	if location == ParentsAndCurrentGrps || location == AllGrps {
		if name, err := group.GetName(false); err != nil {
			return ncGroups, err
		} else {
			ncGroups.Add(name, group)
		}
	}

	// the child groups of the current group
	if location == ChildrenGrps || location == AllChildrenGrps || location == AllGrps {
		// get the number of groups
		numGrps, ncIds, err := NcInqGrps(myId)
		if err != nil {
			return ncGroups, err
		}
		for i := 0; i < numGrps; i++ {
			tmpGroup := NewGroup(ncIds[i])
			name, err := tmpGroup.GetName(false)
			if err != nil {
				return ncGroups, err
			}
			ncGroups.Add(name, tmpGroup)
			fmt.Println(tmpGroup)
		}
	}

	// search in parent groups.
	if location == ParentsGrps || location == ParentsAndCurrentGrps || location == AllGrps {
		tmpGroup := NewGroupFrom(group)
		isRG, err := tmpGroup.IsRootGroup()
		if err != nil {
			return ncGroups, err
		}
		if !isRG {
			for {
				parentGroup, err := tmpGroup.GetParentGroup()
				if err != nil {
					return ncGroups, err
				}
				if parentGroup.IsNull() {
					break
				}
				name, err := parentGroup.GetName(false)
				if err != nil {
					return ncGroups, err
				}
				ncGroups.Add(name, parentGroup)
				tmpGroup = parentGroup
			}
		}
	}
	if group.myId==ID(65537)&&location==3{
		fmt.Println(ncGroups)
	}
	// search in child groups of the children
	if location == ChildrenOfChildrenGrps || location == AllChildrenGrps || location == AllGrps {
		groupMs, err := group.GetGroupsM(ChildrenGrps)
		if err != nil {
			return ncGroups, err
		}
		for groupKey := range groupMs {
			gps := groupMs.EqualRange(groupKey)
			for _, gp := range gps {
				childGroups, err := gp.(Group).GetGroupsM(AllChildrenGrps)

				if err != nil {
					return ncGroups, err
				}
				keys, fields := childGroups.GetAllPair()
				for i := 0; i < len(keys); i++ {
					ncGroups.Add(keys[i].(string), fields[i].(Group))
				}
			}
		}
	}
	//if group.myId==ID(65537)&&location==3{
	//	fmt.Println(ncGroups)
	//}
	return ncGroups, nil
}

// Get the named child Group object.

func (group Group) GetGroup(name string, location GroupLocation /*ChildrenGrps*/) (Group, error) {
	if group.IsNull() {
		return NewGroupNull(), fmt.Errorf("error: attempt to invoke GetParentGroup on a Null group")
	}
	ncGroups, err := group.GetGroupsM(location)
	if err != nil {
		return NewGroupNull(), err
	}
	ret := ncGroups.EqualRange(name)
	if len(ret) == 0 {
		return NewGroupNull(), nil
	} else {
		gp := ret[0].(Group)
		return gp, nil
	}
}

// Get all NcxxGroup objects with a given name.
func (group Group) GetGroups(name string, location GroupLocation) (Set, error) {
	ncSet := NewSet()
	if group.IsNull() {
		return ncSet, fmt.Errorf("error: attempt to invoke GetGroups on a Null group")
	}
	ncGroups, err := group.GetGroupsM(location)
	if err != nil {
		return ncSet, err
	}
	ret := ncGroups.EqualRange(name)
	for gp := range ret {
		ncSet.Add(gp)
	}
	return ncSet, nil
}

// Add a new child Group object.
func (group Group) AddGroup(name string) (Group, error) {
	if group.IsNull() {
		return NewGroupNull(), fmt.Errorf("error: attempt to invoke addGroup on a Null group")
	}
	newId, err := ncDefGrp(group.myId, name)
	if err != nil {
		return NewGroupNull(), err
	}

	return NewGroup(newId), nil

}

/*! Returns true if this object is null (i.e. it has no contents); otherwise returns false. */
func (group Group) IsNull() bool {
	return group.nullObject
}

// Get the number of Var objects in this group.
// Test
func (group Group) getVarCount(location Location) (int, error) {

	// search in current group.
	tmpGroup := NewGroupFrom(group)

	// search in current group
	nvars := 0
	if (location == ParentsAndCurrent || location == ChildrenAndCurrent || location == Current || location == All) && !tmpGroup.IsNull() {
		id, err := tmpGroup.GetId()
		if err != nil {
			return -1, err
		}
		nvars, err = ncInqNvars(id)
		if err != nil {
			return -1, err
		}
	}

	// search recursively in all parent groups.
	if location == Parents || location == ParentsAndCurrent || location == All {
		tmpGroup, err := group.GetParentGroup()
		if err != nil {
			return -1, err
		}
		for !tmpGroup.IsNull() {
			id, err := tmpGroup.GetId()
			if err != nil {
				return -1, err
			}
			nvarsp, err := ncInqNvars(id)
			if err != nil {
				return -1, err
			}
			nvars += nvarsp
			// continue loop with the parent.
			tmpGroup, err = tmpGroup.GetParentGroup()
			if err != nil {
				return -1, err
			}
		}
	}

	// search recursively in all child groups
	if location == ChildrenAndCurrent || location == Children || location == All {
		groups, err := group.GetGroupsM(ParentsGrps)
		if err != nil {
			return -1, err
		}

		_, gps := groups.GetAllPair()

		for _, gp := range gps {
			nvarTmp, err := gp.(Group).getVarCount(ChildrenAndCurrent)
			if err != nil {
				return -1, err
			}
			nvars += nvarTmp
		}
	}
	return nvars, nil
}

// Add a new netCDF variable.
//func (group Group) AddVarByName(name, typeName, dimName string) (Var, error) {
//	CheckDefineMode(group.myId)
//
//	// get an NcxxType object with the given type name.
//	tmpType, err := group.getType(typeName, ParentsAndCurrent)
//	if err != nil {
//		return NewVarNull(), err
//	}
//	if tmpType.IsNull() {
//		return NewVarNull(), fmt.Errorf("error: attempt to invoke Group.AddVarByName" +
//			" failed: typeName must be defined in either the current group or a parent group",)
//	}
//
//	tmpDim := group.GetDim(dimName, ParentsAndCurrent)
//	if (tmpDim.IsNull()) {
//		throw
//		NcxxNullDim("Attempt to invoke NcxxGroup::addVar failed: "
//		"dimName must be defined in either the current "
//		"group or a parent group", __FILE__, __LINE__);
//	}
//
//	// finally define a new netCDF  variable
//	int
//	varId;
//	int
//	dimId(tmpDim.GetId());
//	ncxxCheck(nc_def_var(myId, name.c_str(), tmpType.GetId(), 1, &dimId, &varId),
//		__FILE__, __LINE__,
//		"NcxxGroup::addVar()", GetName(), name);
//	// return an NcxxVar object for this new variable
//	return NcxxVar(*this, varId);
//}

// Gets the NcxxType object with a given name.
func (group Group) getType(name string, location Location) (Type, error) {
	if group.IsNull() {
		return NewTypeNull(), fmt.Errorf("error: attempt to invoke getType on a Null group")
	}

	switch name {
	case "byte":
		return Byte, nil
	case "ubyte":
		return Ubyte, nil
	case "char":
		return Char, nil
	case "short":
		return Short, nil
	case "ushort":
		return Ushort, nil
	case "int":
		return Int, nil
	case "uint":
		return Uint, nil
	case "int64":
		return Int64, nil
	case "uint64":
		return Uint64, nil
	case "float":
		return Float, nil
	case "double":
		return Double, nil
	case "string":
		return String, nil
	default:
		return NewTypeNull(), fmt.Errorf("error: unknown typeName in Group. getType")
	}

	//// TODO add a user defined type
	//// iterator for the multimap container.
	//multimap < string, NcxxType >::iterator; it;
	//// return argument of equal_range: iterators to lower and upper bounds of the range.
	//pair < multimap < string, NcxxType >::iterator, multimap < string, NcxxType >::iterator > ret;
	//// get the entire collection of types.
	//multimap < string, NcxxType > types(getTypes(location));
	//// define STL set object to hold the result
	//set < NcxxType > tmpType;
	//// get the set of NcxxType objects with a given name
	//ret = types.equal_range(name);
	//if (ret.first == ret.second)
	//return NcxxType();
	//else
	//return ret.first- > second;
}

// /////////////
// Dim-related methods
// /////////////

// Get the number of Dim objects.
func (group Group) GetDimCount(location Location /*Current*/) (int, error) {
	if group.IsNull() {
		return -1, fmt.Errorf("error: attempt to invoke GetDimCount on a Null group")
	}
	myId, _ := group.GetId()

	// intialize counter
	ndims := 0

	// search in current group
	if location == Current || location == ParentsAndCurrent || location == ChildrenAndCurrent || location == All {
		ndimsp, err := NcInqNdims(myId)
		if err != nil {
			return -1, err
		}
		ndims += ndimsp
	}
	// search in parent groups.
	if location == Parents || location == ParentsAndCurrent || location == All {
		groups, err := group.GetGroupsM(ParentsGrps)
		if err != nil {
			return -1, err
		}
		_, gps := groups.GetAllPair()
		for _, gp := range gps {
			ndimTmp, err := gp.(Group).GetDimCount(Current)
			if err != nil {
				return -1, err
			}
			ndims += ndimTmp
		}
	}
	// search in child groups.
	if location == Children || location == ChildrenAndCurrent || location == All {
		groups, err := group.GetGroupsM(AllChildrenGrps)
		if err != nil {
			return -1, err
		}
		_, gps := groups.GetAllPair()
		for _, gp := range gps {
			ndimTmp, err := gp.(Group).GetDimCount(Current)
			if err != nil {
				return -1, err
			}
			ndims += ndimTmp
		}

	}
	return ndims, nil
}

// Get the set of Dim objects.
func (group Group) getDimsM(location Location /*Current*/) (Multimap, error) {
	ncDims := NewMultimap() // create a container to hold the Dim's.

	if group.IsNull() {
		return ncDims, fmt.Errorf("error: attempt to invoke getDimsM on a Null group")
	}
	myId, _ := group.GetId()

	// search in current group
	if location == Current || location == ParentsAndCurrent || location == ChildrenAndCurrent || location == All {
		dimCount, err := group.GetDimCount(Current)
		if err != nil {
			return ncDims, err
		}

		if dimCount > 0 {
			_, dimIds, err := NcInqDimids(myId, false)
			if err != nil {
				return ncDims, err
			}

			// now get the name of each Dim and populate the nDims container.
			for i := 0; i < dimCount; i++ {
				tmpDim := NewDim(group, dimIds[i])
				dimName, err := tmpDim.GetName()
				if err != nil {
					return ncDims, err
				}
				ncDims.Add(dimName, tmpDim)
			}
		}
	}

	// search in parent groups.
	if location == Parents || location == ParentsAndCurrent || location == All {
		groups, err := group.GetGroupsM(ParentsGrps)
		if err != nil {
			return ncDims, err
		}
		_, gps := groups.GetAllPair()
		for _, gp := range gps {
			subNcGroups, err := gp.(Group).getDimsM(Current)
			if err != nil {
				return ncDims, err
			}
			keys, fields := subNcGroups.GetAllPair()
			for i := 0; i < len(keys); i++ {
				ncDims.Add(keys[i].(string), fields[i].(Dim))
			}
		}
	}

	// search in child groups (makes recursive calls).
	if location == Children || location == ChildrenAndCurrent || location == All {
		groups, err := group.GetGroupsM(AllChildrenGrps)
		if err != nil {
			return ncDims, err
		}

		_, gps := groups.GetAllPair()
		for _, gp := range gps {
			subNcGroups, err := gp.(Group).getDimsM(Current)
			if err != nil {
				return ncDims, err
			}
			keys, fields := subNcGroups.GetAllPair()
			for i := 0; i < len(keys); i++ {
				ncDims.Add(keys[i].(string), fields[i].(Dim))
			}
		}
	}

	return ncDims, nil
}

//
//// Get the named NcxxDim object.
//NcxxDim NcxxGroup::getDim(const string& name,NcxxGroup::Location location) const {
//if(isNull()) {
//throw NcxxNullGrp("Attempt to invoke NcxxGroup::getDim on a Null group",
//__FILE__, __LINE__);
//}
//multimap<string,NcxxDim> ncDims(getDims(location));
//pair<multimap<string,NcxxDim>::iterator,multimap<string,NcxxDim>::iterator> ret;
//ret = ncDims.equal_range(name);
//if(ret.first == ret.second) {
//return NcxxDim(); // null group is returned
//} else {
//return ret.first->second;
//}
//}
//
//
//// Get all NcxxDim objects with a given name.
//set<NcxxDim> NcxxGroup::getDims(name string,location Location /*Current*/) const {
//
//if(isNull()) {
//throw NcxxNullGrp("Attempt to invoke NcxxGroup::getDims on a Null group",
//__FILE__, __LINE__);
//}
//// get the set of ncDims in this group and above.
//multimap<string,NcxxDim> ncDims(getDims(location));
//pair<multimap<string,NcxxDim>::iterator,multimap<string,NcxxDim>::iterator> ret;
//multimap<string,NcxxDim>::iterator it;
//ret = ncDims.equal_range(name);
//set<NcxxDim> tmpDim;
//for (it=ret.first; it!=ret.second; ++it) {
//tmpDim.insert(it->second);
//}
//return tmpDim;
//}

// Add a new Dim object.
func (group Group) AddDim(name string, dimSize uint) (Dim, error) {
	CheckDefineMode(group.myId)
	if group.IsNull() {
		return NewDimNull(), fmt.Errorf("error: attempt to invoke addDim on a Null group")
	}
	dimId, err := ncDefDim(group.myId, name, SIZE(dimSize))
	if err != nil {
		return NewDimNull(), err
	}
	// finally return Dim object for this new variable
	return NewDim(group, dimId), nil
}

// Add a new Dim object with unlimited size..
func (group Group) AddDimUl(name string) (Dim, error) {
	CheckDefineMode(group.myId)
	if group.IsNull() {
		return NewDimNull(), fmt.Errorf("error: attempt to invoke addDim on a Null group")
	}
	dimId, err := ncDefDim(group.myId, name, NCUNLIMITED)
	if err != nil {
		return NewDimNull(), err
	}
	// finally return Dim object for this new variable
	return NewDim(group, dimId), nil
}
