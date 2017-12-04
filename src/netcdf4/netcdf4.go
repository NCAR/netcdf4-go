package netcdf4

// #cgo LDFLAGS: -lnetcdf
// #include <stdlib.h>
// #include <netcdf.h>
import "C"
import (
	"unsafe"
	"fmt"
)

// FileMode represents a file's mode.
type FileMode int
type FileFormat int

// File modes for Open or Create
//const (
//	SHARE FileMode = C.NC_SHARE // share updates, limit cacheing
//)

// File modes for Open
const (
	WRITE   FileMode = iota
	READ
	NEWFILE
	REPLACE
)

// File format for Create
const (
	CLASSIC        FileFormat = iota //!< Classic format, classic data model
	CLASSIC64                        //!< 64-bit offset format, classic data model
	NETCDF4                          //!< (default) netCDF-4/HDF5 format, enhanced data model
	NETCDF4CLASSIC                   //!< netCDF-4/HDF5 format, classic data model
	UNKNOWN
)

const NCUNLIMITED = C.NC_UNLIMITED

// ID represents a ncId or groupid.
type ID C.int

var gNcid ID
// SIZE represents the type for the size_t.

type SIZE C.size_t

// Error represents an error returned by netCDF C library.

type Error C.int

// Error returns a string representation of Error e.

func (e Error) Error() string {
	return C.GoString(C.nc_strerror(C.int(e)))
}

// Create error from the return code

func newError(n C.int) error {
	if n == C.NC_NOERR {
		return nil
	}
	return Error(n)
}

func ncCreate(path string, fMode FileMode, fFormat FileFormat) (ncId ID, err error) {

	var mode C.int
	var format C.int

	switch fMode {
	case NEWFILE:
		mode = C.NC_NOCLOBBER
	case REPLACE:
		mode = C.NC_CLOBBER
	default:
		return ID(-1), fmt.Errorf("wrong fileMode")
	}

	switch fFormat {
	case CLASSIC:
		format = C.NC_CLASSIC_MODEL
	case CLASSIC64:
		format = C.NC_64BIT_OFFSET
	case NETCDF4:
		format = C.NC_NETCDF4
	case NETCDF4CLASSIC:
		format = C.NC_NETCDF4 | C.NC_CLASSIC_MODEL
	case UNKNOWN:
		format = C.NC_NETCDF4
	default:
		return ID(-1), fmt.Errorf("unknown fileFormat")
	}

	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	var id C.int
	err = newError(C.nc_create(cPath, format|mode, &id))
	ncId = ID(id)
	return
}

func ncOpen(path string, fMode FileMode, fFormat FileFormat) (ncId ID, err error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	var id C.int

	var mode C.int
	var format C.int

	switch fMode {
	case WRITE:
		mode = C.NC_WRITE
	case READ:
		mode = C.NC_NOWRITE
	default:
		return ID(-1), fmt.Errorf("wrong fileMode")
	}

	switch fFormat {
	case CLASSIC:
		format = C.NC_CLASSIC_MODEL
	case CLASSIC64:
		format = C.NC_64BIT_OFFSET
	case NETCDF4:
		format = C.NC_NETCDF4
	case NETCDF4CLASSIC:
		format = C.NC_NETCDF4 | C.NC_CLASSIC_MODEL
	case UNKNOWN:
		err = newError(C.nc_open(cPath, mode, &id))
		ncId = ID(id)
		return
	default:
		ncId = ID(-1)
		err = fmt.Errorf("unknown fileFormat")
		return
	}

	err = newError(C.nc_open(cPath, C.int(format)|C.int(mode), &id))
	ncId = ID(id)
	return
}

///* Learn the path used to open/create the file. */
//EXTERNL int nc_inq_path(int ncid, size_t *pathlen, char *path);
//
///* Given an ncid and group name (NULL gets root group), return locid. */
//EXTERNL int nc_inq_ncid(int ncid, const char *name, int *grp_ncid);
//
///* Given a location id, return the number of groups it contains, and
// * an array of their locids. */
func NcInqGrps(ncId ID) (numGrps int, ncIds []ID, err error) {
	var cNumGrps C.int
	tmp := make([]C.int, 1)

	err = newError(C.nc_inq_grps(C.int(ncId), &cNumGrps, &tmp[0]))
	if err != nil {
		return
	}
	numGrps = int(cNumGrps)
	if numGrps==0{
		return numGrps, []ID(nil), nil
	}
	cNcIds := make([]C.int, numGrps)
	err = newError(C.nc_inq_grps(C.int(ncId), &cNumGrps, &cNcIds[0]))
	if err != nil {
		return
	}
	ncIds = make([]ID, numGrps)
	for i := 0; i < numGrps; i++ {
		ncIds[i] = ID(cNcIds[i])
	}
	return
}

//func NcInqGrps(ncId ID) (numGrps int, ncIds []ID, err error){
//	var cNumGrps C.int
//	//tmp := make([]C.int, 1)
//	var tmp C.int
//	err = newError(C.nc_inq_grps(C.int(ncId), &cNumGrps, &tmp))
//	if err!=nil{
//		return
//	}
//	numGrps = int(cNumGrps)
//
//	ncIds = make([]ID, numGrps)
//
//	//// this method sometime fails maybe due to the gc
//	//for i:=0;i<numGrps;i++{
//	//	tmp2 := (*C.int)(unsafe.Pointer(
//	//		uintptr(unsafe.Pointer(&tmp)) + uintptr(i*int(unsafe.Sizeof(C.int(0))))))
//	//	ncIds[i] = ID(*tmp2)
//	//}
//
//	////method given by  https://github.com/golang/go/wiki/cgo but it may fails
//	////var theCArray *C.YourType = C.getTheArray()
//	////length := C.getTheArrayLength()
//	////slice := (*[1 << 30]C.YourType)(unsafe.Pointer(theCArray))[:length:length]
//	cNcIds := (*[1 << 30]C.int)(unsafe.Pointer(&tmp))[0:numGrps:numGrps]
//	for i:=0;i<numGrps;i++{
//		ncIds[i] = ID(cNcIds[i])
//	}
//	return
//}
//EXTERNL int
//nc_inq_grps(int ncid, int *numgrps, int *ncids);

func ncRedef(ncId ID) (err error) {
	err = newError(C.nc_redef(C.int(ncId)))
	return
}

func ncEnddef(ncId ID) (err error) {
	err = newError(C.nc_enddef(C.int(ncId)))
	return
}
func ncSync(ncId ID) (err error) {
	err = newError(C.nc_sync(C.int(ncId)))
	return
}
func ncAbort(ncId ID) (err error) {
	err = newError(C.nc_abort(C.int(ncId)))
	return
}

func ncClose(ncId ID) (err error) {
	err = newError(C.nc_close(C.int(ncId)))
	return
}

func ncInq(ncId ID) (ndimsp, nvarsp, nattsp, unlimdimidp int, err error) {
	var cNdimsp, cNvarsp, cNattsp, cUnlimdimidp C.int
	err = newError(C.nc_inq(C.int(ncId), &cNdimsp, &cNvarsp, &cNattsp, &cUnlimdimidp))
	ndimsp = int(cNdimsp)
	nvarsp = int(cNvarsp)
	nattsp = int(cNattsp)
	unlimdimidp = int(cUnlimdimidp)
	return
}

func NcInqNdims(ncId ID) (ndimsp int, err error) {
	var cNdimsp C.int
	err = newError(C.nc_inq_ndims(C.int(ncId), &cNdimsp))
	ndimsp = int(cNdimsp)
	return
}

func ncInqNvars(ncId ID) (nvarsp int, err error) {
	var cNvarsp C.int
	err = newError(C.nc_inq_nvars(C.int(ncId), &cNvarsp))
	nvarsp = int(cNvarsp)
	return
}

func ncInqNatts(ncId ID) (nattsp int, err error) {
	var cNattsp C.int
	err = newError(C.nc_inq_natts(C.int(ncId), &cNattsp))
	nattsp = int(cNattsp)
	return
}

func ncInqUnlimdim(ncId ID) (unlimdimidp int, err error) {
	var cUnlimdimidp C.int
	err = newError(C.nc_inq_unlimdim(C.int(ncId), &cUnlimdimidp))
	unlimdimidp = int(cUnlimdimidp)
	return
}

func ncInqUnlimdims(ncId ID) (nunlimdimsp, unlimdimidsp int, err error) {
	var cNunlimdimsp, cUnlimdimidsp C.int
	err = newError(C.nc_inq_unlimdims(C.int(ncId), &cNunlimdimsp, &cUnlimdimidsp))
	nunlimdimsp = int(cNunlimdimsp)
	unlimdimidsp = int(cUnlimdimidsp)
	return
}

///* Get a list of ids for all the variables in a group. */
//EXTERNL int
//nc_inq_varids(int ncid, int *nvars, int *varids);

/* Find all dimids for a location. This finds all dimensions in a
 * group, or any of its parents. */
//EXTERNL int
//nc_inq_dimids(int ncid, int *ndims, int *dimids, int include_parents);

func NcInqDimids(ncId ID, includeParents bool) (nDims int, dimIds []ID, err error) {
	var cNumDims C.int
	tmp := make([]C.int, 1)
	cIncludeParents := C.int(0)
	if includeParents {
		cIncludeParents = C.int(1)
	}
	err = newError(C.nc_inq_dimids(C.int(ncId), &cNumDims, &tmp[0], cIncludeParents))
	if err != nil {
		return
	}
	nDims = int(cNumDims)
	cDimIds := make([]C.int, nDims)
	err = newError(C.nc_inq_dimids(C.int(ncId), &cNumDims, &cDimIds[0], cIncludeParents))
	if err != nil {
		return
	}
	dimIds = make([]ID, nDims)
	for i := 0; i < nDims; i++ {
		dimIds[i] = ID(cDimIds[i])
	}
	return
}

///* Find all user-defined types for a location. This finds all
// * user-defined types in a group. */
//EXTERNL int
//nc_inq_typeids(int ncid, int *ntypes, int *typeids);
//
///* Are two types equal? */
//EXTERNL int
//nc_inq_type_equal(int ncid1, nc_type typeid1, int ncid2,
//nc_type typeid2, int *equal);

/* Create a group. its ncId is returned as newId. */

func ncDefGrp(parentId ID, name string) (newId ID, err error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var id C.int
	err = newError(C.nc_def_grp(C.int(parentId), cName, &id))
	newId = ID(id)
	return
}

/* Given locid, find name of group. (Root group is named "/".) */

func ncInqGrpname(ncId ID) (name string, err error) {
	cName := C.CString(string(make([]byte, C.NC_MAX_NAME+1)))
	defer C.free(unsafe.Pointer(cName))
	err = newError(C.nc_inq_grpname(C.int(ncId), cName))
	name = C.GoString(cName)
	return
}

/* Given ncId, find full name and len of full name. (Root group is
 * named "/", with length 1.) */

func ncInqGrpnameFull(ncId ID) (name string, err error) {
	lenGrpname, err := ncInqGrpnameLen(ncId)
	if err != nil {
		name = ""
		return
	}
	cfullName := C.CString(string(make([]byte, lenGrpname+1)))
	defer C.free(unsafe.Pointer(cfullName))
	var lenp C.size_t
	err = newError(C.nc_inq_grpname_full(C.int(ncId), &lenp, cfullName))
	name = C.GoString(cfullName)
	return
}

/* Given ncId, find len of full name. */

func ncInqGrpnameLen(ncId ID) (C.size_t, error) {
	var lenp C.size_t
	err := newError(C.nc_inq_grpname_len(C.int(ncId), &lenp))
	return lenp, err
}

/* Given an ncId, find the ncId of its parent group. */

func ncInqGrpParent(ncId ID) (parentId ID, err error) {
	var id C.int
	err = newError(C.nc_inq_grp_parent(C.int(ncId), &id))
	parentId = ID(id)
	return
}

/* Begin _dim */

/* Create a group. its ncId is returned as newId. */

func ncDefDim(ncId ID, name string, dimSize SIZE) (dimId ID, err error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var id C.int
	err = newError(C.nc_def_dim(C.int(ncId), cName, C.size_t(dimSize), &id))
	dimId = ID(id)
	return
}

func ncInqDimid(ncId ID, name string) (dimId ID, err error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var id C.int
	err = newError(C.nc_inq_dimid(C.int(ncId), cName, &id))
	dimId = ID(id)
	return
}

//EXTERNL int
//nc_inq_dim(int ncid, int dimid, char *name, size_t *lenp);
//

func NcInqDimname(ncId ID, dimId ID) (name string, err error) {
	cName := C.CString(string(make([]byte, C.NC_MAX_NAME+1)))
	defer C.free(unsafe.Pointer(cName))
	err = newError(C.nc_inq_dimname(C.int(ncId), C.int(dimId), cName))
	name = C.GoString(cName)
	return
}

func NcInqDimLen(ncId ID, dimId ID) (C.size_t, error) {
	var lenp C.size_t
	err := newError(C.nc_inq_dimlen(C.int(ncId), C.int(dimId), &lenp))
	return lenp, err
}

func NcRenameDim(ncId ID,dimId ID, name string) (err error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	err = newError(C.nc_rename_dim(C.int(ncId), C.int(dimId), cName))
	return
}

/* End _dim */
