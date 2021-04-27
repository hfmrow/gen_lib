// gitReleases.go

// Source file auto-generated on Fri, 27 Sep 2019 19:41:41 using Gotk3ObjHandler v1.3.8 ©2018-19 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

/*************/
/* Releases */
/***********/

// GetReleases: Information about published releases,
// https://developer.github.com/v3/repos/releases/
func (g *GitReleases) List(login, repos string) (err error) {
	var resp *http.Response
	if err = g.init(); err == nil {
		if resp, err = http.Get("https://api.github.com/repos/" + login + "/" + repos + "/releases"); err == nil {
			defer resp.Body.Close()
			err = json.NewDecoder(resp.Body).Decode(g.Data)
		}
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Unexpected status code, %d\n", resp.StatusCode))
	}
	return
}

func (g *GitReleases) Write(filename string) (err error) {
	var jsonData []byte
	var out bytes.Buffer
	if jsonData, err = json.Marshal(g.Data); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(filename, out.Bytes(), os.ModePerm)
		}
	}
	return
}

func (g *GitReleases) Read(filename string) (err error) {
	if err = g.init(); err == nil {
		var textFileBytes []byte
		if textFileBytes, err = ioutil.ReadFile(filename); err == nil {
			err = json.Unmarshal(textFileBytes, g.Data)
		}
	}
	return
}

func (g *GitReleases) init() (err error) {
	g.Data = new(releasesData)
	if g.Data == nil {
		err = errors.New("Cannot initialize the data structure.")
	}
	return
}

// GitReleases: https://developer.github.com/v3/repos/releases/
type GitReleases struct {
	oauth2Client *http.Client
	Data         *releasesData
	Login        string
}
type releasesData []struct {
	Assets []struct {
		BrowserDownloadURL string      `json:"browser_download_url"`
		ContentType        string      `json:"content_type"`
		CreatedAt          string      `json:"created_at"`
		DownloadCount      int         `json:"download_count"`
		ID                 int         `json:"id"`
		Label              interface{} `json:"label"`
		Name               string      `json:"name"`
		NodeID             string      `json:"node_id"`
		Size               int         `json:"size"`
		State              string      `json:"state"`
		UpdatedAt          string      `json:"updated_at"`
		Uploader           struct {
			AvatarURL         string `json:"avatar_url"`
			EventsURL         string `json:"events_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			GravatarID        string `json:"gravatar_id"`
			HTMLURL           string `json:"html_url"`
			ID                int    `json:"id"`
			Login             string `json:"login"`
			NodeID            string `json:"node_id"`
			OrganizationsURL  string `json:"organizations_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			ReposURL          string `json:"repos_url"`
			SiteAdmin         bool   `json:"site_admin"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			Type              string `json:"type"`
			URL               string `json:"url"`
		} `json:"uploader"`
		URL string `json:"url"`
	} `json:"assets"`
	AssetsURL string `json:"assets_url"`
	Author    struct {
		AvatarURL         string `json:"avatar_url"`
		EventsURL         string `json:"events_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		GravatarID        string `json:"gravatar_id"`
		HTMLURL           string `json:"html_url"`
		ID                int    `json:"id"`
		Login             string `json:"login"`
		NodeID            string `json:"node_id"`
		OrganizationsURL  string `json:"organizations_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		ReposURL          string `json:"repos_url"`
		SiteAdmin         bool   `json:"site_admin"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		Type              string `json:"type"`
		URL               string `json:"url"`
	} `json:"author"`
	Body            string `json:"body"`
	CreatedAt       string `json:"created_at"`
	Draft           bool   `json:"draft"`
	HTMLURL         string `json:"html_url"`
	ID              int    `json:"id"`
	Name            string `json:"name"`
	NodeID          string `json:"node_id"`
	Prerelease      bool   `json:"prerelease"`
	PublishedAt     string `json:"published_at"`
	TagName         string `json:"tag_name"`
	TarballURL      string `json:"tarball_url"`
	TargetCommitish string `json:"target_commitish"`
	UploadURL       string `json:"upload_url"`
	URL             string `json:"url"`
	ZipballURL      string `json:"zipball_url"`
}
