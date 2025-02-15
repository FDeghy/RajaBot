package mrbilit

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

func GetTrains(src, dst, date string) ([]*Trains, error) {
	form := url.Values{}
	form.Add("from", src)
	form.Add("to", dst)
	form.Add("date", date)
	form.Add("genderCode", "3")
	form.Add("adultCount", "1")
	form.Add("childCount", "0")
	form.Add("infantCount", "0")
	form.Add("exclusive", "false")
	form.Add("availableStatus", "Both")

	req, _ := http.NewRequest(
		http.MethodGet,
		BASE_URL+"/api/GetAvailable/v2?"+form.Encode(),
		nil,
	)
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
	req.Header.Add("x-playerid", uuid.New().String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, ErrBadHttpCode
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := &Response{}
	err = json.Unmarshal(respBody, data)
	if err != nil {
		return nil, ErrJsonDecode
	}

	return data.Trains, nil
}
