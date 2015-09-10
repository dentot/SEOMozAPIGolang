# SEOMozAPIGolang
Golang for Moz API - currently very limited, only for the url-metrics


Install:

go get github.com/dentot/SEOMozAPIGolang

Usage:

	urls := []string{"http://google.com", "http://yahoo.com"}

	var cols uint64 = MozApi.URL_METRICS_PAGE_AUTHORITY | MozApi.URL_METRICS_DOMAIN_AUTHORITY

	//supply your accessid and key
	b, err := MozApi.GetURLMetrics("mozscape-accessid", "secretkey", cols, 300, urls)
	if err != nil {
		panic(err)
	}

	// check if no error from moz api server
	isError, data, err := MozApi.CheckResultError(b)
	if isError {
		panic(data["error_message"])
	}

	// if no error, then extract the moz data
	res, err := MozApi.ExtractMozData(b)
	fmt.Println("The first item, upa of google.com", res[0]["upa"])
	fmt.Println("The 2nd item, upa of yahoo.com", res[1]["upa"])
