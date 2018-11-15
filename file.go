package netcdf4

import (
	"fmt"
)

//File represnets an opened netCDF file
type File struct {
	Group
	pathInUse string
	mode      FileMode
	format    FileFormat
}

//NewFile creates a new file with an empty group set
func NewFile() (f File) {
	return File{
		Group: NewGroupNull(),
	}
}

//Open opens a file
func (f *File) Open(filePath string, fMode FileMode, fFormat FileFormat) (err error) {
	if !f.nullObject {
		if err := f.Close(); err != nil {
			return err
		}
	}

	f.format = fFormat
	f.mode = fMode
	if f.mode == WRITE || f.mode == READ {
		f.id, err = Open(filePath, f.mode, f.format)
	} else if f.mode == NEWFILE || f.mode == REPLACE {
		f.id, err = Create(filePath, f.mode, f.format)
	} else {
		return fmt.Errorf("error wrong filemode in File.Open")
	}

	if err != nil {
		return err
	}
	f.pathInUse = filePath
	gNcid = f.id
	f.nullObject = false
	return
}

//Close closes the opened NetCDF file
func (f *File) Close() error {
	if !f.nullObject {
		gNcid = -1
		err := ncClose(f.id)
		if err != nil {
			return err
		}
	}
	f.nullObject = true
	f.pathInUse = ""
	//f.errStr.clear()
	f.format = NETCDF4
	f.mode = READ
	return nil
}

// Sync forces a Synchronization of an open netcdf dataset to disk
func (f File) Sync() error {
	return ncSync(f.id)
}

//Enddef leaves define mode, used for classic model
func (f File) Enddef() error {
	return ncEnddef(f.id)
}

//GetPathInUse returns the current path in use
func (f File) GetPathInUse() string {
	return f.pathInUse
}
