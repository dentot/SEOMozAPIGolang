package MozApi

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorAccessId(t *testing.T) {
	assert := assert.New(t)

	urls := []string{"http://google.com", "http://yahoo.com"}

	b, err := GetURLMetrics(
		"wrongaccessid",
		"wrongkey",
		URL_METRICS_PAGE_AUTHORITY|URL_METRICS_DOMAIN_AUTHORITY,
		300,
		urls)

	//fmt.Println()
	assert.Nil(err, "Should not have error")

	// result should have error_message from moz server
	isError, data, err := CheckResultError(b)

	assert.True(isError, "Should have an error message since wrong access id")
	assert.Contains(data["error_message"], "authentication failed", "Should contain authentication failed")

	fmt.Println(string(b))
}

func TestBatchGetURLMetrics(t *testing.T) {
	assert := assert.New(t)

	urls := []string{"http://google.com", "http://yahoo.com"}

	// provide your accessid and secret key
	b, err := GetURLMetrics(
		"mozscape-accessid",
		"secretkey",
		URL_METRICS_PAGE_AUTHORITY|URL_METRICS_DOMAIN_AUTHORITY,
		300,
		urls)

	assert.Nil(err, "Should not have error")

	isError, _, _ := CheckResultError(b)
	assert.False(isError, "Moz server should not send an error json")

	data, err := ExtractMozData(b)

	assert.Nil(err, "Should not have error extracting data")
	fmt.Println("The data", data)

	assert.Equal(len(data), 2, "Returned data should have exactly 2 items.")
	assert.NotNil(data[0]["upa"], "First item(google.com) must not be null")
	assert.NotNil(data[1]["upa"], "2nd item(yahoo.com) must not be null")

}
