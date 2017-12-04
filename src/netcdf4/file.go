package netcdf4

import (
	"fmt"
)

type file struct {
	Group
	pathInUse string
	mode      FileMode
	format    FileFormat
}

func File() (f file) {
	f.Group = NewGroupNull()
	f.Close()
	return
}

func (f *file) Open(filePath string, fMode FileMode, fFormat FileFormat) (err error) {

	if !f.nullObject {
		err = f.Close()
		if err != nil {
			return err
		}
	}

	f.format = fFormat
	f.mode = fMode
	if f.mode == WRITE || f.mode == READ {
		f.myId, err = ncOpen(filePath, f.mode, f.format)
	} else if f.mode == NEWFILE || f.mode == REPLACE {
		f.myId, err = ncCreate(filePath, f.mode, f.format)
	} else {
		err = fmt.Errorf("error wrong filemode in file.Open")
	}

	if err != nil {
		return
	}
	f.pathInUse = filePath
	gNcid = f.myId
	f.nullObject = false
	return
}

func (f *file) Close() error {
	if !f.nullObject {
		gNcid = -1
		err := ncClose(f.myId)
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

/////////////////////////////////////////////
// Synchronize an open netcdf dataset to disk

func (f file) Sync() error {
	return ncSync(f.myId)
}

//////////////////////////////////////////////
// Leave define mode, used for classic model

func (f file) Enddef() error {
	return ncEnddef(f.myId)
}

//////////////////////////////////////////////
//! get the path in use

func (f file) GetPathInUse() string {
	return f.pathInUse
}
