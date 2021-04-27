// gitBranches.go

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
	"strings"
)

/*************/
/* Branches */
/***********/

// GetBranches: GET /repos/:owner/:repo/branches
func (g *GitBranches) List(login, repos string) (err error) {
	var resp *http.Response
	if err = g.init(); err == nil {
		if resp, err = http.Get("https://api.github.com/repos/" + login + "/" + repos + "/branches"); err == nil {
			defer resp.Body.Close()
			err = json.NewDecoder(resp.Body).Decode(g.Data)
		}
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Unexpected status code, %d\n", resp.StatusCode))
	}
	return
}

// NewFromMaster: create new branch
func (g *GitBranches) NewFromMaster(user, repos, newBranch string) (out []byte, err error) {
	var resp *http.Response

	// Get information on "master" branch
	if resp, err = g.oauth2Client.Get("https://api.github.com/repos/" + user + "/" + repos + "/git/refs/heads"); err == nil {
		defer resp.Body.Close()
		branchHead := new(BranchHead)
		if err = json.NewDecoder(resp.Body).Decode(branchHead); err == nil {
			branchNew := new(BranchNew)
			for _, branch := range *branchHead {
				// Search for "master" branch, usualy [0] but ...
				if strings.Contains(branch.Ref, "refs/heads/master") {
					branchNew.Sha = branch.Object.Sha
					splitted := strings.Split(branch.Ref, "/")
					splitted = splitted[:len(splitted)-1]
					splitted = append(splitted, newBranch)
					branchNew.Ref = strings.Join(splitted, "/")
					break
				}
			}
			// Create branch
			var jsonData []byte
			if jsonData, err = json.Marshal(branchNew); err == nil {
				if resp, err = g.oauth2Client.Post("https://api.github.com/repos/"+user+"/"+repos+"/git/refs",
					"application/json", bytes.NewReader(jsonData)); err == nil {
					defer resp.Body.Close()
					out, err = ioutil.ReadAll(resp.Body)
				}
			}
		}
	}
	if resp.StatusCode != http.StatusCreated {
		err = errors.New(fmt.Sprintf("Unexpected status code, %d\n", resp.StatusCode))
	}
	return
}

func (g *GitBranches) Write(filename string) (err error) {
	var jsonData []byte
	var out bytes.Buffer
	if jsonData, err = json.Marshal(g.Data); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(filename, out.Bytes(), os.ModePerm)
		}
	}
	return
}

func (g *GitBranches) Read(filename string) (err error) {
	if err = g.init(); err == nil {
		var textFileBytes []byte
		if textFileBytes, err = ioutil.ReadFile(filename); err == nil {
			err = json.Unmarshal(textFileBytes, g.Data)
		}
	}
	return
}

func (g *GitBranches) init() (err error) {
	g.Data = new(branchesData)
	if g.Data == nil {
		err = errors.New("Cannot initialize the data structure.")
	}
	return
}

// https://developer.github.com/v3/repos/branches/
type GitBranches struct {
	oauth2Client *http.Client
	Data         *branchesData
	Login        string
}
type branchesData []struct {
	Commit struct {
		Sha string `json:"sha"`
		URL string `json:"url"`
	} `json:"commit"`
	Name       string `json:"name"`
	Protected  bool   `json:"protected"`
	Protection struct {
		Enabled              bool `json:"enabled"`
		RequiredStatusChecks struct {
			Contexts         []string `json:"contexts"`
			EnforcementLevel string   `json:"enforcement_level"`
		} `json:"required_status_checks"`
	} `json:"protection"`
	ProtectionURL string `json:"protection_url"`
}

// BranchHead: Used to retrieve "head branch" when using
// https://api.github.com/repos/USER/REPO/git/refs/heads
type BranchHead []struct {
	NodeID string `json:"node_id"`
	Object struct {
		Sha  string `json:"sha"`
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"object"`
	Ref string `json:"ref"`
	URL string `json:"url"`
}
type BranchNew struct {
	Ref string `json:"ref"`
	Sha string `json:"sha"`
}
