package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var Host = "http://10.255.44.33"
var ChallengeURL string = Host + "/cgi-bin/get_challenge"
var UserInfoURL string = Host + "/cgi-bin/rad_user_info"
var LoginURL string = Host + "/cgi-bin/srun_portal"

var defaultHeaders = map[string]string{
	"Accept":           "text/javascript, application/javascript, application/ecmascript, application/x-ecmascript, */*; q=0.01",
	"Accept-Encoding":  "gzip, deflate",
	"Accept-Language":  "zh-CN,zh;q=0.9,en;q=0.8",
	"Cache-Control":    "no-cache",
	"Connection":       "keep-alive",
	"Host":             "10.255.44.33",
	"Pragma":           "no-cache",
	"Referer":          "http://10.255.44.33/srun_portal_pc?ac_id=1&theme=pro",
	"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0",
	"X-Requested-With": "XMLHttpRequest",
}

type QueryParm map[string]string
type ResponseBody map[string]interface{}

func encodeQueryParm(query QueryParm) string {
	v := url.Values{}
	for key, value := range query {
		v.Add(key, value)
	}
	return v.Encode()
}

func sendGetRequest(url string, query QueryParm) (ResponseBody, error) {
	url += "?" + encodeQueryParm(query)

	req, _ := http.NewRequest("GET", url, nil)
	for key, value := range defaultHeaders {
		req.Header.Add(key, value)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return parseResponse(string(bodyBytes)), nil
}

func RequestUserInfo() (ResponseBody, error) {
	return sendGetRequest(UserInfoURL, QueryParm{
		"callback": generateCallback(),
	})
}

func RequestChallenge(username string) (ResponseBody, error) {
	return sendGetRequest(ChallengeURL, QueryParm{
		"callback": generateCallback(),
		"username": username,
		// "ip":       ip,
		"_": CurrentMilisecond(),
	})
}

func RequestLogin(username string, password string, ip string, challenge string) (ResponseBody, error) {
	_info := fmt.Sprintf(`{"username":"%s","password":"%s","ip":"%s","acid":"1","enc_ver":"srun_bx1"}`, username, password, ip)
	info := encodeUserInfo(_info, challenge)

	hmd5 := encodeMD5(password, challenge)
	chksum := encodeChksum(
		strings.Join([]string{username, hmd5[5:], "1", ip, "200", "1", info}, challenge),
		challenge,
	)

	queryParm := QueryParm{
		"callback":     generateCallback(),
		"action":       "login",
		"username":     username,
		"password":     hmd5,
		"os":           "Windows 10",
		"name":         "Windows",
		"double_stack": "0",
		"chksum":       chksum,
		"info":         info,
		"ac_id":        "1",
		"ip":           ip,
		"n":            "200",
		"type":         "1",
		"_":            CurrentMilisecond(),
	}

	return sendGetRequest(LoginURL, queryParm)
}

func generateCallback() string {
	// 经过测试，其实 callback 字符串可以随意填写，但是为了保险起见，还是仿照原来的格式生成
	return "jQuery_11240619677586845294_" + CurrentMilisecond()
}

func parseResponse(str string) ResponseBody {
	r := regexp.MustCompile(`^jQuery_\d+_\d+\((.+?)\)$`)
	match := r.FindStringSubmatch(str)
	if len(match) <= 1 {
		log.Println("Error: parse response failed")
		return nil
	}
	var jsonMap ResponseBody
	json.Unmarshal([]byte(match[1]), &jsonMap)
	return jsonMap
}
