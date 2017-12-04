package main

// #include <stdlib.h>
// #include <netcdf.h>
import "C"
import (
	"netcdf4"
	"log"
)

func main() {

	file := netcdf4.File()
	//err := file.Open("/Users/hhuang/Downloads/regions.nc", netcdf4.NEWFILE, netcdf4.NETCDF4)
	err := file.Open("/Users/hhuang/Downloads/regions.nc", netcdf4.REPLACE, netcdf4.UNKNOWN)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	grpUSA, err := file.AddGroup("USA")

	if err != nil {
		log.Fatal(err)
	}

	grpColorado, err := grpUSA.AddGroup("Colorado")
	if err != nil {
		log.Fatal(err)
	}
	grpWyoming, err := grpUSA.AddGroup("Wyoming")
	if err != nil {
		log.Fatal(err)
	}
	grpAlaska, err := grpUSA.AddGroup("Alaska")
	if err != nil {
		log.Fatal(err)
	}

	/* define dimensions */
	USATimeDim, err := grpUSA.AddDimUl("time")
	if err != nil {
		log.Fatal(err)
	}
	ColoradoStationsDim, err := grpColorado.AddDim("stations", 5)
	if err != nil {
		log.Fatal(err)
	}
	WyomingStationsDim, err := grpWyoming.AddDim("stations", 4)
	if err != nil {
		log.Fatal(err)
	}
	AlaskaStationsDim, err := grpAlaska.AddDim("stations", 3)
	if err != nil {
		log.Fatal(err)
	}

	_ = USATimeDim
	_ = ColoradoStationsDim
	_ = WyomingStationsDim
	_ = AlaskaStationsDim
	//fmt.Println(grpUSA.GetDimCount(netcdf4.Children))
	grpUSA.GetGroupsM(netcdf4.AllChildrenGrps)
	//fmt.Println(k)

	//TODO convert to a test file
}
