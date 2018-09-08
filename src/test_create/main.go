package main

// #include <stdlib.h>
// #include <netcdf.h>
import "C"
import (
	"netcdf4"
	"log"
	"fmt"
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
	//USATimeDim, err := grpUSA.AddDimUl("time") //no length included
	USATimeDim, err := grpUSA.AddDim("time", 2) //no length included
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
	_, err = grpAlaska.AddDim("time", 5)
	if err != nil {
		log.Fatal(err)
	}

	_ = USATimeDim
	_ = ColoradoStationsDim
	_ = WyomingStationsDim
	_ = AlaskaStationsDim

	//[]string{"time","stations2"}
	tempVar, err := grpWyoming.AddVar("average_temperature", netcdf4.Float, []netcdf4.Dim{USATimeDim, WyomingStationsDim})
	//tempVar2, err := grpWyoming.AddVar("average_temperature1", netcdf4.Float, []netcdf4.Dim{})
	fmt.Println(tempVar.DataLength())
	//fmt.Println(tempVar2.DataLength())

	//idW,_:= grpWyoming.GetId()
	//idC,_:= grpColorado.GetId()
	//idA,_:= grpAlaska.GetId()
	//idU,_:= grpUSA.GetId()
	//fmt.Println(grpAlaska.GetDim("time",netcdf4.ParentsAndCurrent))
	//fmt.Println(grpAlaska.GetDims("time",netcdf4.ParentsAndCurrent))
	//fmt.Println(grpWyoming.GetDim("time",netcdf4.ParentsAndCurrent))
	//fmt.Println(netcdf4.NcInqDimids(idW,true))
	//fmt.Println(netcdf4.NcInqDimids(idC,true))
	//fmt.Println(netcdf4.NcInqDimids(idA,true))
	//fmt.Println(netcdf4.NcInqDimids(idU,true))
	//tempVar, err := grpWyoming.AddVarScalar("average_temperature1", netcdf4.Double, )
	//tempVar.PutValAll(2.)
	//fmt.Println(tempVar,err)
	//fmt.Println(grpUSA.GetDimCount(netcdf4.All))
	//fmt.Println(grpUSA.GetDimsM(netcdf4.All))
	//fmt.Println(grpUSA.GetDims("stations", netcdf4.All))
	//fmt.Println(grpUSA.GetGroups("Alaska", netcdf4.AllGrps))
	//fmt.Println(grpUSA.GetGroup("Alaska", netcdf4.AllGrps))
	//grpUSA.GetGroupsM(netcdf4.AllChildrenGrps)
	//fmt.Println(k)
	//TODO convert to a test file
}

//
//func main()  {
//	fmt.Println(reflect.TypeOf([]int{1.,})==reflect.SliceOf(reflect.TypeOf(int(0))))
//
//}
