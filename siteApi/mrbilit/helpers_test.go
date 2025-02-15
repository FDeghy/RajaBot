package mrbilit_test

import (
	"RajaBot/siteApi/mrbilit"
	"testing"
	"time"

	ptime "github.com/yaa110/go-persian-calendar"
)

func TestGetTrains(t *testing.T) {
	src := "1"
	dst := "161"
	date := ptime.Date(1403, 11, 29, 0, 0, 0, 0, ptime.Iran()).Time().Format("2006-01-02")

	trains, err := mrbilit.GetTrains(src, dst, date)
	if err != nil {
		t.Error(err)
	}

	for _, tr := range trains {
		ti, _ := time.Parse("2006-01-02T15:04:05", tr.DepartureTime)
		pt := ptime.New(ti)
		p := tr.Prices[0].Classes[0]
		t.Logf("%v %v %v", pt.Format("HH:mm E dd MMM yyyy"), p.Capacity, p.Price)
	}
}
