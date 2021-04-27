// url_curl.go

/*
	©2021 hfmrow, https://github.com/hfmrow

	package url_curl v1.0
	cointain curl url operations.

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php

	Other library used:

	-https://github.com/andelf/go-curl
	 Copyright 2014 Shuyu Wang <andelf@gmail.com>
	 Apache License 2.0

	-Make Sure You Have libcurl (and its develop headers, static/dynamic libs) installed!
	 ubuntu: 'sudo apt install libcurl4-gnutls-dev'
*/

package url_curl

import (
	curl "github.com/andelf/go-curl"
)

// GetUrl:
func GetUrl(url string, callback func(data []byte)) error {

	curl.GlobalInit(curl.GLOBAL_ALL)
	easy := curl.EasyInit()
	defer easy.Cleanup()
	easy.Setopt(curl.OPT_URL, url)
	// easy.Setopt(curl.OPT_VERBOSE, true)
	easy.Setopt(curl.OPT_WRITEFUNCTION,
		func(ptr []byte, userdata interface{}) bool {
			callback(ptr)
			return true
		})
	return easy.Perform()
}
