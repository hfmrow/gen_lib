// gitSearch.go

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
	"io/ioutil"
	"net/http"
	"os"
)

/***********/
/* Search */
/*********/

// For: execute a search command. https://developer.github.com/v3/search/
// - for a complete list of available qualifiers.
// https://help.github.com/en/articles/searching-on-github
// - how to use operators to match specific quantities, dates, or to exclude.
// https://help.github.com/articles/understanding-the-search-syntax/
func (g *GitSearch) For(target, query string) (err error) {
	var resp *http.Response
	if err = g.init(); err == nil {
		if resp, err = http.Get("https://api.github.com/search/" + target + "?q=" + query); err == nil {
			defer resp.Body.Close()
			JsonToMap(resp.Body)
			g.Status, err = getStatus("get", resp.StatusCode)
			err = json.NewDecoder(resp.Body).Decode(g.Data)
		}
	}
	return
}

func (g *GitSearch) Write(filename string) (err error) {
	var jsonData []byte
	var out bytes.Buffer
	if jsonData, err = json.Marshal(g.Data); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(filename, out.Bytes(), os.ModePerm)
		}
	}
	return
}

func (g *GitSearch) Read(filename string) (err error) {
	if err = g.init(); err == nil {
		var textFileBytes []byte
		if textFileBytes, err = ioutil.ReadFile(filename); err == nil {
			err = json.Unmarshal(textFileBytes, g.Data)
		}
	}
	return
}

func (g *GitSearch) init() (err error) {
	g.Data = new(searchData)
	if g.Data == nil {
		err = errors.New("Cannot initialize the data structure.")
	}
	return
}

// GitSearch: https://developer.github.com/v3/search/
type GitSearch struct {
	oauth2Client *http.Client
	Data         *searchData
	Status       string
	Login        string
}
type searchData struct {
	IncompleteResults bool `json:"incomplete_results"`
	Items             []struct {
		ArchiveURL       string `json:"archive_url"`
		Archived         bool   `json:"archived"`
		AssigneesURL     string `json:"assignees_url"`
		BlobsURL         string `json:"blobs_url"`
		BranchesURL      string `json:"branches_url"`
		CloneURL         string `json:"clone_url"`
		CollaboratorsURL string `json:"collaborators_url"`
		CommentsURL      string `json:"comments_url"`
		CommitsURL       string `json:"commits_url"`
		CompareURL       string `json:"compare_url"`
		ContentsURL      string `json:"contents_url"`
		ContributorsURL  string `json:"contributors_url"`
		CreatedAt        string `json:"created_at"`
		DefaultBranch    string `json:"default_branch"`
		DeploymentsURL   string `json:"deployments_url"`
		Description      string `json:"description"`
		Disabled         bool   `json:"disabled"`
		DownloadsURL     string `json:"downloads_url"`
		EventsURL        string `json:"events_url"`
		Fork             bool   `json:"fork"`
		Forks            int64  `json:"forks"`
		ForksCount       int64  `json:"forks_count"`
		ForksURL         string `json:"forks_url"`
		FullName         string `json:"full_name"`
		GitCommitsURL    string `json:"git_commits_url"`
		GitRefsURL       string `json:"git_refs_url"`
		GitTagsURL       string `json:"git_tags_url"`
		GitURL           string `json:"git_url"`
		HasDownloads     bool   `json:"has_downloads"`
		HasIssues        bool   `json:"has_issues"`
		HasPages         bool   `json:"has_pages"`
		HasProjects      bool   `json:"has_projects"`
		HasWiki          bool   `json:"has_wiki"`
		Homepage         string `json:"homepage"`
		HooksURL         string `json:"hooks_url"`
		HTMLURL          string `json:"html_url"`
		ID               int64  `json:"id"`
		IssueCommentURL  string `json:"issue_comment_url"`
		IssueEventsURL   string `json:"issue_events_url"`
		IssuesURL        string `json:"issues_url"`
		KeysURL          string `json:"keys_url"`
		LabelsURL        string `json:"labels_url"`
		Language         string `json:"language"`
		LanguagesURL     string `json:"languages_url"`
		License          struct {
			Key    string `json:"key"`
			Name   string `json:"name"`
			NodeID string `json:"node_id"`
			SpdxID string `json:"spdx_id"`
			URL    string `json:"url"`
		} `json:"license"`
		MergesURL        string      `json:"merges_url"`
		MilestonesURL    string      `json:"milestones_url"`
		MirrorURL        interface{} `json:"mirror_url"`
		Name             string      `json:"name"`
		NodeID           string      `json:"node_id"`
		NotificationsURL string      `json:"notifications_url"`
		OpenIssues       int64       `json:"open_issues"`
		OpenIssuesCount  int64       `json:"open_issues_count"`
		Owner            struct {
			AvatarURL         string `json:"avatar_url"`
			EventsURL         string `json:"events_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			GravatarID        string `json:"gravatar_id"`
			HTMLURL           string `json:"html_url"`
			ID                int64  `json:"id"`
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
		} `json:"owner"`
		Private         bool    `json:"private"`
		PullsURL        string  `json:"pulls_url"`
		PushedAt        string  `json:"pushed_at"`
		ReleasesURL     string  `json:"releases_url"`
		Score           float64 `json:"score"`
		Size            int64   `json:"size"`
		SSHURL          string  `json:"ssh_url"`
		StargazersCount int64   `json:"stargazers_count"`
		StargazersURL   string  `json:"stargazers_url"`
		StatusesURL     string  `json:"statuses_url"`
		SubscribersURL  string  `json:"subscribers_url"`
		SubscriptionURL string  `json:"subscription_url"`
		SvnURL          string  `json:"svn_url"`
		TagsURL         string  `json:"tags_url"`
		TeamsURL        string  `json:"teams_url"`
		TreesURL        string  `json:"trees_url"`
		UpdatedAt       string  `json:"updated_at"`
		URL             string  `json:"url"`
		Watchers        int64   `json:"watchers"`
		WatchersCount   int64   `json:"watchers_count"`
	} `json:"items"`
	TotalCount int64 `json:"total_count"`
}

var qualifiers = map[string]string{
	"Repositories":        "repositories",
	"Code":                "code",
	"Commits":             "commits",
	"Issues":              "issues",
	"Users":               "users",
	"Topics":              "topics",
	"Text match metadata": "text-match"}
