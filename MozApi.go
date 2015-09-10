// MozApi project MozApi.go
package MozApi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// HMAC-SHA1 hash of your Access ID, the Expires parameter, and your Secret Key.
func getHmac_sha1(accessid string, expires string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(accessid))
	h.Write([]byte("\n"))
	h.Write([]byte(expires))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Does a http post to the moz api url, with your accessid, secret key, url metrics option, and the urls you want to get url metrics
// sample call: GetURLMetrics(xxxx, keyxxxx, URL_METRICS_PAGE_AUTHORITY | URL_METRICS_DOMAIN_AUTHORITY, 300, []string{"google.com", "yahoo.com"})
func GetURLMetrics(accessid string, key string, cols uint64, secondsExpire time.Duration, urls []string) ([]byte, error) {

	currentTime := time.Now()

	//adds expiry
	t := currentTime.Add(time.Second * secondsExpire).Unix()
	//convert to string
	expireTime := strconv.FormatInt(t, 10)

	//build signed auth
	signed := getHmac_sha1(accessid, expireTime, key)

	var mozUrl bytes.Buffer

	//build now the url
	mozUrl.WriteString("https://lsapi.seomoz.com/linkscape/url-metrics/?")
	mozUrl.WriteString("Cols=")
	mozUrl.WriteString(strconv.FormatUint(cols, 10))
	mozUrl.WriteString("&AccessID=")
	mozUrl.WriteString(accessid)
	mozUrl.WriteString("&Expires=")
	mozUrl.WriteString(expireTime)
	mozUrl.WriteString("&Signature=")
	mozUrl.WriteString(url.QueryEscape(signed))
	//fmt.Println(mozUrl.String())

	b, err := json.Marshal(urls)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("post data %v.\n", string(b))

	req, err := http.NewRequest("POST", mozUrl.String(), bytes.NewBuffer(b))
	if err != nil {
		// handle error
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	return body, nil
}

// Since return type of moz api can be:
// {status:401, error_message:"the error"}
// or
// [{pda:82}, {pda:71}]
// so we need a handler to check if return has error, or has the needed moz data values
// if no 'error_message', then proceed to extract the moz return data using the ExtractMozData
func CheckResultError(b []byte) (bool, map[string]interface{}, error) {

	if strings.Contains(string(b), "error_message") {
		var data map[string]interface{}
		if err := json.Unmarshal(b, &data); err != nil {
			//fmt.Printf("Unmarshall err: %v.\n", err)
			return true, nil, err
		}

		return true, data, nil
	}

	return false, nil, nil
}

// just convert the json string from moz api call to array of map
func ExtractMozData(b []byte) ([](map[string]interface{}), error) {
	var docs [](map[string]interface{})
	if err := json.Unmarshal(b, &docs); err != nil {
		fmt.Printf("Unmarshall err: %v.\n", err)
		return nil, err
	}

	return docs, nil
}
