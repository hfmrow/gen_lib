// tools.go

// Source file auto-generated on Wed, 25 Sep 2019 16:10:45 using Gotk3ObjHandler v1.3.8 ©2018-19 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package http

import (
	"net/http"
	"time"
)

// func main() {
// 	fmt.Println(IsExistUrl("https://raw.githubusercontent.com/hfmrow/fake-repo/main/version"))
// }

// IsExistUrl: Check if an rl is available (not 404)
func IsExistUrl(url string) (isUrlOk bool) {
	var err error
	var resp *http.Response

	timeout := time.Duration(10 * time.Second)
	client := http.Client{Timeout: timeout}

	if resp, _ = client.Get(url); err == nil {
		isUrlOk = resp.StatusCode != 404
	}
	return
}
