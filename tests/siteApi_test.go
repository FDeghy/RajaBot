package tests

import (
	siteapi "RajaBot/siteApi"
	"testing"
)

func TestGetTrains(T *testing.T) {
	routes, err := siteapi.GetRoutes()
	if err != nil {
		T.Error(err)
	}

	for _, i := range routes {
		src, dst, err := i.GetStationsID()
		if err != nil {
			T.Error(err)
			break
		}
		T.Log(i.ID, i.Name, src, dst)
		trains, err := siteapi.GetTrains(src, dst, "1403/03/15")
		if err != nil {
			T.Error(err)
			break
		}
		for _, j := range trains {
			T.Logf("%+v", j)
		}

		break
	}
}
