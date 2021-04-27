// version_control.go

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

package version_control

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	glhpul "github.com/hfmrow/gen_lib/http/url"
)

// Finds the latest version published in the repository. This assumes that the target
// repository / branch contains a valid "version" file. 'newDebFileUrl' contains the
// new URL pointing to the version of the update file, 'to' contains the current
// software version, 'repos' must be filled with the targeted repository,
// 'appName' contains the name of the current application, it will be written in the
// "version" file, 'doControl' means connect to the repository and check the 'version'
// file, 'writeVersionFile' means that during the development process the 'version' file
// is generated and written to the current local repository and will be uploaded  on the
// next 'git push '. The 'callback' function, gives error information if there is,
// "updtAvailable" is set to "true" if the update is available or "false" otherwise.
// NOTE: The default branch is set to “master”.
func VersionControl(newDebFileUrl, vers, repos, appName string, doControl, writeVersionFile bool,
	callback func(updtAvailable bool, newVersion, newDeb string, err error), branch ...string) error {

	var (
		err     error
		_branch string = "master"
		newVersion,
		newDebFile string
		reUrlError = regexp.MustCompile(`(?mi)invalid|request|not|found`)
		reGetInfos = regexp.MustCompile(`(?m)=(.*)`)
	)
	if len(branch) > 0 {
		_branch = branch[0]
	}
	if writeVersionFile {
		fileMtx := appName + ` current version information.
current=` + vers + `
url-new=` + newDebFileUrl

		err := ioutil.WriteFile("version", []byte(fileMtx), 0644)
		if err != nil {
			return fmt.Errorf("Unable to write file version: %v\n", err)
		}
	}
	if doControl {
		if !strings.HasPrefix(repos, "https://") {
			repos = "https://" + repos
		}
		if !strings.HasSuffix(repos, "/") {
			repos += "/"
		}
		repos += _branch + "/"
		url := strings.ReplaceAll(repos, "github.com", "raw.githubusercontent.com")
		url += "version"
		go func() {
			err = glhpul.GetUrl(url,
				func(data []byte) {
					str := string(data)
					if reUrlError.MatchString(str) {
						err = fmt.Errorf("Unable to [GetUrl] [%s] current version: %s\n", url, str)
						return
					}
					found := reGetInfos.FindAllStringSubmatch(str, -1)
					if len(found) != 2 {
						err = fmt.Errorf("Unable to get current version information #1.")
						return
					}
					if len(found[0]) != 2 {
						err = fmt.Errorf("Unable to get current version information #2.")
						return
					}
					if len(found[1]) != 2 {
						err = fmt.Errorf("Unable to get current version information #3.")
						return
					}
					newVersion = found[0][1]
					newDebFile = found[1][1]
				})
			callback(!strings.Contains(newVersion, vers), newVersion, newDebFile, err)
		}()
	}
	return err
}
