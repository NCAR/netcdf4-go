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
//    https://github.com/Unidata/netcdf-c
//
//  Modification for Go made by:
//
//    Hao Huang, Nanjing University
//    Email: hhuangwx@gmail.com
//
//////////////////////////////////////////////////////////////////////

//Group represents a netCDF4 group
type Group struct {
	nullObject bool
	id         ID

	// option to use the 'proposed_standard_name' attribute instead
	// of 'standard_name'.
	useProposedStandardName bool
}

/*!
  GroupLocation is an enumeration list contains the options for selecting groups (used for returned set of Group objects).
*/
type GroupLocation int

//Known enums
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

//String conforms to fmt.Stringer interface
func (i GroupLocation) String() string {
	if i < 0 || i >= GroupLocation(len(groupLocationIndex)-1) {
		return "GroupLocation(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return groupLocationName[groupLocationIndex[i]:groupLocationIndex[i+1]]
}

//Location is ....
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

//NewGroupNull returns an empty group
func NewGroupNull() Group {
	return Group{nullObject: true, id: -1, useProposedStandardName: false}
}

//NewGroup returns a new group where its ID is set
func NewGroup(groupId ID) Group {
	return Group{nullObject: false, id: groupId, useProposedStandardName: false}
}

//NewGroupFrom creates a new group from the pass parent??
func NewGroupFrom(rhs Group) (group Group) {
	return Group{nullObject: rhs.nullObject, id: rhs.id, useProposedStandardName: rhs.useProposedStandardName}
}

// /////////////
// Group-related methods
// /////////////

// Name gets the group name
func (g Group) Name(fullName bool) (string, error) {
	if g.IsNull() {
		return "", fmt.Errorf("error: attempt to invoke GetName on a Null group")
	}
	if fullName {
		// return full name of group with foward "/" separarating sub-groups.
		return ncInqGrpnameFull(g.id)
	}
	// return the (local) name of this group.
	return ncInqGrpname(g.id)
}

// IsRootGroup returns true if this is the group root.
func (g Group) IsRootGroup() (bool, error) {
	grpName, err := g.Name(false)
	if err == nil {
		return grpName == "/", nil
	}
	return false, err
}

//GetParentGroup returns the parent group. Get the parent group.
func (g Group) GetParentGroup() (Group, error) {
	if g.IsNull() {
		return NewGroupNull(), fmt.Errorf("error: attempt to invoke GetParentGroup on a Null group")
	}

	if parentID, err := ncInqGrpParent(g.id); err == nil {
		return NewGroup(parentID), nil
	}
	//if no parent id is found, return null group
	return NewGroupNull(), nil

}

// ID returns the group ID
func (g Group) ID() (ID, error) {
	if g.IsNull() {
		return ID(-1), fmt.Errorf("error: attempt to invoke GetId on a Null group")
	}
	return g.id, nil
}

// GetGroupCount returns the number of ??children ??? groups? objects.
func (g Group) GetGroupCount(location GroupLocation) (int, error) {
	if g.IsNull() {
		return -1, fmt.Errorf("error: attempt to invoke GetGroupCount on a Null group")
	}

	n := 0 // initialize group counter

	// record this group
	if location == ParentsAndCurrentGrps || location == AllGrps {
		n++
	}

	// number of children in current group
	if location == ChildrenGrps || location == AllChildrenGrps || location == AllGrps {
		numGrps, _, err := NcInqGrps(g.id)
		if err != nil {
			return -1, err
		}
		n += numGrps
	}

	// search in parent groups
	if location == ParentsGrps || location == ParentsAndCurrentGrps || location == AllGrps {
		groups, err := g.GetGroupsM(ParentsGrps)
		if err != nil {
			return -1, err
		}
		n += len(groups)
	}
	// get the number of all children that are childreof children
	if location == ChildrenOfChildrenGrps || location == AllChildrenGrps || location == AllGrps {
		groups, err := g.GetGroupsM(ChildrenOfChildrenGrps)
		if err != nil {
			return -1, err
		}
		n += len(groups)
	}

	return n, nil
}

// GetGroupsM retrieves the ..... ?? Get the set of child Group objects.
func (g Group) GetGroupsM(location GroupLocation) (MultimapG, error) {
	ncGroups := NewMultimapG()
	if g.IsNull() {
		return ncGroups, fmt.Errorf("error: attempt to invoke GetGroupsM on a Null group")
	}

	// record this group
	if location == ParentsAndCurrentGrps || location == AllGrps {
		if name, err := g.Name(false); err != nil {
			return ncGroups, err
		} else {
			ncGroups.Add(name, g)
		}
	}

	// the child groups of the current group
	if location == ChildrenGrps || location == AllChildrenGrps || location == AllGrps {
		// get the number of groups
		numGrps, ncIds, err := NcInqGrps(g.id)
		if err != nil {
			return ncGroups, err
		}
		for i := 0; i < numGrps; i++ {
			tmpGroup := NewGroup(ncIds[i])
			name, err := tmpGroup.Name(false)
			if err != nil {
				return ncGroups, err
			}
			ncGroups.Add(name, tmpGroup)
		}
	}

	// search in parent groups.
	if location == ParentsGrps || location == ParentsAndCurrentGrps || location == AllGrps {
		tmpGroup := NewGroupFrom(g)
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
				name, err := parentGroup.Name(false)
				if err != nil {
					return ncGroups, err
				}
				ncGroups.Add(name, parentGroup)
				tmpGroup = parentGroup
			}
		}
	}

	// search in child groups of the children
	if location == ChildrenOfChildrenGrps || location == AllChildrenGrps || location == AllGrps {
		groupMs, err := g.GetGroupsM(ChildrenGrps)
		if err != nil {
			return ncGroups, err
		}
		for groupKey := range groupMs {
			gps := groupMs.EqualRange(groupKey)
			for _, gp := range gps {
				childGroups, err := gp.GetGroupsM(AllChildrenGrps)

				if err != nil {
					return ncGroups, err
				}
				keys, fields := childGroups.GetAllPair()
				for i := 0; i < len(keys); i++ {
					ncGroups.Add(keys[i], fields[i])
				}
			}
		}
	}
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
		gp := ret[0]
		return gp, nil
	}
}

// Get all Group objects with a given name.

func (group Group) GetGroups(name string, location GroupLocation) (SetG, error) {
	ncSetG := NewSetG()
	if group.IsNull() {
		return ncSetG, fmt.Errorf("error: attempt to invoke GetGroups on a Null group")
	}
	ncGroups, err := group.GetGroupsM(location)
	if err != nil {
		return ncSetG, err
	}
	ret := ncGroups.EqualRange(name)
	fmt.Println(ret)

	for _, gp := range ret {
		ncSetG.Add(gp)
	}
	return ncSetG, nil
}

// Add a new child Group object.
func (group Group) AddGroup(name string) (Group, error) {
	if group.IsNull() {
		return NewGroupNull(), fmt.Errorf("error: attempt to invoke addGroup on a Null group")
	}
	newId, err := ncDefGrp(group.id, name)
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
func (group Group) GetVarCount(location Location /*Current*/) (int, error) {

	// search in current group.
	tmpGroup := NewGroupFrom(group)

	// search in current group
	nvars := 0
	if (location == ParentsAndCurrent || location == ChildrenAndCurrent ||
		location == Current || location == All) && !tmpGroup.IsNull() {
		id, err := tmpGroup.ID()
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
			id, err := tmpGroup.ID()
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
			nvarTmp, err := gp.GetVarCount(ChildrenAndCurrent)
			if err != nil {
				return -1, err
			}
			nvars += nvarTmp
		}
	}
	return nvars, nil
}

// Get the collection of Var objects.

func (group Group) GetVarsM(location Location) (MultimapV, error) {
	ncVars := NewMultimapV() // create a container to hold the Var's.
	myId, err := group.ID()
	if err != nil {
		return ncVars, err
	}
	// search in current group.
	tmpGroup := NewGroupFrom(group)

	if (location == ParentsAndCurrent || location == ChildrenAndCurrent ||
		location == Current || location == All) && !tmpGroup.IsNull() {
		// get the number of variables.
		varCount, varIds, err := NcInqVarids(myId)
		if err != nil {
			return ncVars, err
		}
		for i := 0; i < varCount; i++ {
			tmpVar := NewVar(group, varIds[i])
			varName, err := tmpVar.GetName()
			if err != nil {
				return ncVars, err
			}
			ncVars.Add(varName, tmpVar)
		}
	}

	// search recursively in all parent groups.
	if location == Parents || location == ParentsAndCurrent || location == All {
		tmpGroup, err = group.GetParentGroup()
		if err != nil {
			return ncVars, err
		}
		for !tmpGroup.IsNull() {
			// get the number of variables
			// get the number of variables.
			tmpID, _ := tmpGroup.ID()
			varCount, varIds, err := NcInqVarids(tmpID)
			if err != nil {
				return ncVars, err
			}
			for i := 0; i < varCount; i++ {
				tmpVar := NewVar(group, varIds[i])
				varName, err := tmpVar.GetName()
				if err != nil {
					return ncVars, err
				}
				ncVars.Add(varName, tmpVar)
			}

			// continue loop with the parent.
			tmpGroup, err = tmpGroup.GetParentGroup()
			if err != nil {
				return ncVars, err
			}
		}
	}

	// search recusively in all child groups.
	if location == ChildrenAndCurrent || location == Children || location == All {
		groupMs, err := group.GetGroupsM(ChildrenGrps)
		if err != nil {
			return ncVars, err
		}

		for groupKey := range groupMs {
			gps := groupMs.EqualRange(groupKey)
			for _, gp := range gps {
				varsM, err := gp.GetVarsM(ChildrenAndCurrent)

				if err != nil {
					return ncVars, err
				}
				keys, values := varsM.GetAllPair()
				for i := 0; i < len(keys); i++ {
					ncVars.Add(keys[i], values[i])
				}
			}
		}
	}

	return ncVars, nil
}

// Get all Var objects with a given name.
func (group Group) GetVars(name string, location Location /*Current*/) (SetV, error) {
	tmpVar := NewSetV()
	ncVars, err := group.GetVarsM(location)
	if err != nil {
		return tmpVar, err
	}
	ret := ncVars.EqualRange(name)

	for _, v := range ret {
		tmpVar.Add(v)
	}
	return tmpVar, nil
}

// Get the named Var object.
func (group Group) GetVar(name string, location Location /*Current*/) (Var, error) {
	ncVars, err := group.GetVarsM(location)
	if err != nil {
		return NewVarNull(), err
	}
	ret := ncVars.EqualRange(name)

	if len(ret) == 0 {
		return NewVarNull(), nil
	} else {
		v := ret[0]
		return v, nil
	}
}

// Add a new netCDF variable.
func (group Group) AddVarScalar(name string, varType interface{}) (Var, error) {
	return group.AddVar(name, varType, []string{})
}

// Add a new netCDF variable.
func (group Group) AddVar(name string, varType, dims interface{}) (Var, error) {
	CheckDefineMode(group.id)
	var typeId NcType
	var dimIDs []ID
	errType := fmt.Errorf("io error:attempt to invoke Group.addVar failed: varType " +
		"should be defined as either Type or string in either the current group or a parent group")
	errDim := fmt.Errorf("io error: attempt to invoke Group.addVar failed: " +
		"dims must be defined as Dim or string in either the current group or a parent group")
	switch vType := varType.(type) {
	case string:
		{
			tmpType, err := group.GetType(vType, ParentsAndCurrent)
			if err != nil {
				return NewVarNull(), err
			}
			if tmpType.IsNull() {
				return NewVarNull(), errType
			}
			typeId = tmpType.GetId()
		}
	case Type:
		if vType.IsNull() {
			return NewVarNull(), errType
		}
		typeId = vType.GetId()
	default:
		return NewVarNull(), errType

	}

	switch dimTmp := dims.(type) {
	case string:
		{
			tmpDim, err := group.GetDim(dimTmp, ParentsAndCurrent)
			if err != nil {
				return NewVarNull(), err
			}
			if tmpDim.IsNull() {
				return NewVarNull(), errDim
			}
			dimIDs = append(dimIDs, tmpDim.ID())
		}
	case []string:
		{
			for _, dimName := range dimTmp {
				tmpDim, err := group.GetDim(dimName, ParentsAndCurrent)
				if err != nil {
					return NewVarNull(), err
				}
				if tmpDim.IsNull() {
					return NewVarNull(), errDim
				}
				dimIDs = append(dimIDs, tmpDim.ID())
				fmt.Println(tmpDim)
			}
		}
	case Dim:
		{
			if dimTmp.IsNull() {
				return NewVarNull(), errDim
			}
			isValid, err := dimTmp.IsValidDim(group)
			if err != nil {
				return NewVarNull(), errDim
			}
			if !isValid {
				return NewVarNull(), fmt.Errorf("io error: Dim is not the valid dimension for this group")
			}
			dimIDs = append(dimIDs, dimTmp.ID())
		}
	case []Dim:
		{
			for _, tmpDim := range dimTmp {
				if tmpDim.IsNull() {
					return NewVarNull(), errDim
				}
				isValid, err := tmpDim.IsValidDim(group)
				if err != nil {
					return NewVarNull(), errDim
				}
				if !isValid {
					return NewVarNull(), fmt.Errorf("io error: Dim is not the valid dimension for this group")
				}
				dimIDs = append(dimIDs, tmpDim.ID())
			}
		}
	default:
		return NewVarNull(), errDim
	}
	// finally define a new netCDF  variable varId;
	varId, err := NcDefVar(group.id, name, typeId, dimIDs)
	// return an Var object for this new variable
	return NewVar(group, varId), err
}

// Gets the Type object with a given name.
func (group Group) GetType(name string, location Location) (Type, error) {
	if group.IsNull() {
		return NewTypeNull(), fmt.Errorf("error: attempt to invoke GetType on a Null group")
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
		return NewTypeNull(), fmt.Errorf("error: unknown typeName in Group. GetType")
	}

	//// TODO add a user defined type
	//// iterator for the multimap container.
	//multimap < string, Type >::iterator; it;
	//// return argument of equal_range: iterators to lower and upper bounds of the range.
	//pair < multimap < string, Type >::iterator, multimap < string, Type >::iterator > ret;
	//// get the entire collection of types.
	//multimap < string, Type > types(GetTypes(location));
	//// define STL set object to hold the result
	//set < Type > tmpType;
	//// get the set of Type objects with a given name
	//ret = types.equal_range(name);
	//if (ret.first == ret.second)
	//return Type();
	//else
	//return ret.first- > second;
}

// Adds a new netCDF Enum type.
//NcxxEnumType NcxxGroup::addEnumType(const string& name,NcxxEnumType::ncEnumType baseType) const {
//ncxxCheckDefineMode(myId);
//nc_type typeId;
//ncxxCheck(nc_def_enum(myId, baseType, name.c_str(), &typeId),
//__FILE__, __LINE__,
//"NcxxGroup::addEnumType()", getName(), name);
//NcxxEnumType ncTypeTmp(*this,name);
//return ncTypeTmp;
//}
//
//
//// Adds a new netCDF Vlen type.
//NcxxVlenType NcxxGroup::addVlenType(const string& name,NcxxType& baseType) const {
//ncxxCheckDefineMode(myId);
//nc_type typeId;
//ncxxCheck(nc_def_vlen(myId,  const_cast<char*>(name.c_str()),baseType.getId(),&typeId),
//__FILE__, __LINE__,
//"NcxxGroup::addVlenType()", getName(), name);
//NcxxVlenType ncTypeTmp(*this,name);
//return ncTypeTmp;
//}
//
//
//// Adds a new netCDF Opaque type.
//NcxxOpaqueType NcxxGroup::addOpaqueType(const string& name, size_t size) const {
//ncxxCheckDefineMode(myId);
//nc_type typeId;
//ncxxCheck(nc_def_opaque(myId, size,const_cast<char*>(name.c_str()), &typeId),
//__FILE__, __LINE__,
//"NcxxGroup::addOpaqueType()", getName(), name);
//NcxxOpaqueType ncTypeTmp(*this,name);
//return ncTypeTmp;
//}
//
//// Adds a new netCDF UserDefined type.
//NcxxCompoundType NcxxGroup::addCompoundType(const string& name, size_t size) const {
//ncxxCheckDefineMode(myId);
//nc_type typeId;
//ncxxCheck(nc_def_compound(myId, size,const_cast<char*>(name.c_str()),&typeId),
//__FILE__, __LINE__,
//"NcxxGroup::addCompoundType()", getName(), name);
//NcxxCompoundType ncTypeTmp(*this,name);
//return ncTypeTmp;
//}

// /////////////
// Dim-related methods
// /////////////

// Get the number of Dim objects.
func (group Group) GetDimCount(location Location /*Current*/) (int, error) {
	if group.IsNull() {
		return -1, fmt.Errorf("error: attempt to invoke GetDimCount on a Null group")
	}
	myId, _ := group.ID()

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
			ndimTmp, err := gp.GetDimCount(Current)
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
			ndimTmp, err := gp.GetDimCount(Current)
			if err != nil {
				return -1, err
			}
			ndims += ndimTmp
		}

	}
	return ndims, nil
}

// Get the set of Dim objects.
func (group Group) GetDimsM(location Location /*Current*/) (MultimapD, error) {
	ncDims := NewMultimapD() // create a container to hold the Dim's.

	if group.IsNull() {
		return ncDims, fmt.Errorf("error: attempt to invoke GetDimsM on a Null group")
	}
	myId, _ := group.ID()

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
				dimName, err := tmpDim.Name()
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
			subNcGroups, err := gp.GetDimsM(Current)
			if err != nil {
				return ncDims, err
			}
			keys, fields := subNcGroups.GetAllPair()
			for i := 0; i < len(keys); i++ {
				ncDims.Add(keys[i], fields[i])
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
			subNcGroups, err := gp.GetDimsM(Current)
			if err != nil {
				return ncDims, err
			}
			keys, fields := subNcGroups.GetAllPair()
			for i := 0; i < len(keys); i++ {
				ncDims.Add(keys[i], fields[i])
			}
		}
	}

	return ncDims, nil
}

//// Get the named Dim object.

func (group Group) GetDim(name string, location Location /*Current*/) (Dim, error) {
	if group.IsNull() {
		return NewDimNull(), fmt.Errorf("error: attempt to invoke GetDim on a Null group")
	}
	ncDims, err := group.GetDimsM(location)
	if err != nil {
		return NewDimNull(), err
	}
	ret := ncDims.EqualRange(name)
	if len(ret) == 0 {
		return NewDimNull(), nil
	} else {
		gp := ret[0] //if there are multiple, get the current first
		return gp, nil
	}
}

// Get all Dim objects with a given name.

func (group Group) GetDims(name string, location Location) (SetD, error) {
	ncSetD := NewSetD()
	if group.IsNull() {
		return ncSetD, fmt.Errorf("error: attempt to invoke GetGroups on a Null group")
	}
	ncDims, err := group.GetDimsM(location)
	if err != nil {
		return ncSetD, err
	}
	ret := ncDims.EqualRange(name)
	for _, dimS := range ret {
		ncSetD.Add(dimS)
	}
	return ncSetD, nil
}

// Add a new Dim object.

func (group Group) AddDim(name string, dimSize uint) (Dim, error) {
	CheckDefineMode(group.id)
	if group.IsNull() {
		return NewDimNull(), fmt.Errorf("error: attempt to invoke addDim on a Null group")
	}
	dimId, err := ncDefDim(group.id, name, SIZE(dimSize))
	if err != nil {
		return NewDimNull(), err
	}
	// finally return Dim object for this new variable
	return NewDim(group, dimId), nil
}

// Add a new Dim object with unlimited size..

func (group Group) AddDimUl(name string) (Dim, error) {
	CheckDefineMode(group.id)
	if group.IsNull() {
		return NewDimNull(), fmt.Errorf("error: attempt to invoke addDim on a Null group")
	}
	dimId, err := ncDefDim(group.id, name, NCUNLIMITED)
	if err != nil {
		return NewDimNull(), err
	}
	// finally return Dim object for this new variable
	return NewDim(group, dimId), nil
}
