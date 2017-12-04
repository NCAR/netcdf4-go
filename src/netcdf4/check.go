package netcdf4

// Function checks if the file (group) is in define mode.
// If not, it places it in the define mode.
// While this is automatically done by the underlying C API
// for netCDF-4 files, the netCDF-3 files still need this call.

func CheckDefineMode(ncId ID) error {
	ncRedef(ncId)
	//if status != NC_EINDEFINE{ //TODO redefine
	//	return
	//}
	return nil
}

// Function checks if the file (group) is in data mode.
// If not, it places it in the data mode.
// While this is automatically done by the underlying C API
// for netCDF-4 files, the netCDF-3 files still need this call.

//void CheckDataMode(int ncid, string context /* = "" */)
//{
//int status = nc_enddef(ncid);
//if (status != NC_ENOTINDEFINE) ncxxCheck(status, __FILE__, __LINE__);
//}
