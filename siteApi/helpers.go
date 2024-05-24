package siteapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func GetRoutes() (Routes, error) {
	var routes []*Route

	resp, err := http.DefaultClient.Get(BASE_URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	re, _ := regexp.Compile(`<a onclick=\"SetShortcut\('([^']+)','([^']+)'\)\" style=\"cursor:pointer\">`)
	match := re.FindAllStringSubmatch(string(respBody), -1)
	if match == nil {
		return nil, fmt.Errorf("routes not found")
	}

	for _, i := range match {
		r := &Route{ID: i[1], Name: i[2]}
		r.GetStationsID()
		routes = append(routes, r)
	}
	return routes, nil
}

func GetTrains(src, dst, date string) ([]*Train, error) {
	form := url.Values{}
	form.Add("IsGoBack", "false")
	form.Add("mSrcStation", src)
	form.Add("mTrgStation", dst)
	form.Add("GoDate", date)
	form.Add("BackDate", "")

	req, _ := http.NewRequest(
		http.MethodPost,
		BASE_URL+"/Ticket/GetStation",
		strings.NewReader(form.Encode()),
	)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data := []*Train{}
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
