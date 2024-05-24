package siteapi

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

// Get and Set
func (r *Route) GetStationsID() (string, string, error) {
	var src, dst string

	jar, _ := cookiejar.New(nil)
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.Header.Del("Content-Type")
			return nil
		},
		Jar: jar,
	}
	form := url.Values{}
	form.Add("mID", r.ID)
	form.Add("mCaption", r.Name)
	req, _ := http.NewRequest(
		http.MethodPost,
		BASE_URL+"/Home/Shortcut",
		strings.NewReader(form.Encode()),
	)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	re, _ := regexp.Compile(`sourceValue = (\d+);`)
	match := re.FindStringSubmatch(string(respBody))
	if match == nil {
		return "", "", fmt.Errorf("src not found")
	}
	src = match[1]

	re, _ = regexp.Compile(`destinationValue = (\d+);`)
	match = re.FindStringSubmatch(string(respBody))
	if match == nil {
		return "", "", fmt.Errorf("dst not found")
	}
	dst = match[1]

	// set
	r.Src, r.Dst = src, dst

	// return
	return src, dst, nil
}

func (rts Routes) FindRoute(id string) *Route {
	for _, i := range rts {
		if i.ID == id {
			return i
		}
	}
	return nil
}
