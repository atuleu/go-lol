package lol

import "fmt"

var getter *RESTStaticGetter
var regionTest *Region

func init() {
	data, err := Asset("data/go-lol_testdata.json")
	if err != nil {
		panic(fmt.Sprintf("Could not load test data: %s", err))
	}

	getter, err = NewRESTStaticGetter(data)
	if err != nil {
		panic(fmt.Sprintf("Could not parse static data: %s", err))
	}

	fmt.Printf("Region code is %s\n", getter.RegionCode())

	regionTest, err = NewRegionByCode(getter.RegionCode())
	if err != nil {
		panic(err)
	}

}
