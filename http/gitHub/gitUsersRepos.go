// gitUsersRepos.go

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

/****************************/
/* List users Repositories */
/**************************/

// List: public repositories for the specified user.
// user repos: GET /users/:username/repos
// https://developer.github.com/v3/repos/#list-user-repositories
func (g *GitUsersRepos) List(login string) (err error) {
	if err = g.init(); err == nil {
		var resp *http.Response
		if resp, err = http.Get("https://api.github.com/users/" + login + "/repos"); err == nil {
			defer resp.Body.Close()
			err = json.NewDecoder(resp.Body).Decode(g.Data)
		}
	}
	return
}

// NewRepository: create new repository, "POST /user/repos"
// UserReposNew need to be initialised with default values before use.
// https://developer.github.com/v3/repos/#create
func (g *GitUsersRepos) ReposNew(newRepos *UserReposNew) (out []byte, err error) {
	var resp *http.Response
	var jsonData []byte
	// Create Repository
	if jsonData, err = json.Marshal(newRepos); err == nil {
		if resp, err = g.oauth2Client.Post("https://api.github.com/user/repos",
			"application/json", bytes.NewReader(jsonData)); err == nil {
			defer resp.Body.Close()
			out, err = ioutil.ReadAll(resp.Body)
		}
	}
	return
}

// RemoveRepostory:
func (g *GitUsersRepos) ReposRemove(repos string) (out []byte, err error) {
	var resp *http.Response
	var request *http.Request

	if request, err = http.NewRequest("DELETE", "https://api.github.com/repos/"+g.Login+"/"+repos, nil); err == nil {
		if resp, err = g.oauth2Client.Do(request); err == nil {
			defer resp.Body.Close()
			resp.Body.Read(out)
			// out, err = ioutil.ReadAll(resp.Body)
		}
	}
	return
}

func (g *GitUsersRepos) Write(filename string) (err error) {
	var jsonData []byte
	var out bytes.Buffer
	if jsonData, err = json.Marshal(g.Data); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(filename, out.Bytes(), os.ModePerm)
		}
	}
	return
}

func (g *GitUsersRepos) Read(filename string) (err error) {
	if err = g.init(); err == nil {
		var textFileBytes []byte
		if textFileBytes, err = ioutil.ReadFile(filename); err == nil {
			err = json.Unmarshal(textFileBytes, g.Data)
		}
	}
	return
}

func (g *GitUsersRepos) init() (err error) {
	g.Data = new(usersReposData)
	if g.Data == nil {
		err = errors.New("Cannot initialize the data structure.")
	}
	return
}

// GitUserRepos: https://developer.github.com/v3/search/
type GitUsersRepos struct {
	oauth2Client *http.Client
	Data         *usersReposData
	Login        string
}

type usersReposData []struct {
	ArchiveURL       string      `json:"archive_url"`
	Archived         bool        `json:"archived"`
	AssigneesURL     string      `json:"assignees_url"`
	BlobsURL         string      `json:"blobs_url"`
	BranchesURL      string      `json:"branches_url"`
	CloneURL         string      `json:"clone_url"`
	CollaboratorsURL string      `json:"collaborators_url"`
	CommentsURL      string      `json:"comments_url"`
	CommitsURL       string      `json:"commits_url"`
	CompareURL       string      `json:"compare_url"`
	ContentsURL      string      `json:"contents_url"`
	ContributorsURL  string      `json:"contributors_url"`
	CreatedAt        string      `json:"created_at"`
	DefaultBranch    string      `json:"default_branch"`
	DeploymentsURL   string      `json:"deployments_url"`
	Description      string      `json:"description"`
	Disabled         bool        `json:"disabled"`
	DownloadsURL     string      `json:"downloads_url"`
	EventsURL        string      `json:"events_url"`
	Fork             bool        `json:"fork"`
	ForksCount       int         `json:"forks_count"`
	ForksURL         string      `json:"forks_url"`
	FullName         string      `json:"full_name"`
	GitCommitsURL    string      `json:"git_commits_url"`
	GitRefsURL       string      `json:"git_refs_url"`
	GitTagsURL       string      `json:"git_tags_url"`
	GitURL           string      `json:"git_url"`
	HasDownloads     bool        `json:"has_downloads"`
	HasIssues        bool        `json:"has_issues"`
	HasPages         bool        `json:"has_pages"`
	HasProjects      bool        `json:"has_projects"`
	HasWiki          bool        `json:"has_wiki"`
	Homepage         string      `json:"homepage"`
	HooksURL         string      `json:"hooks_url"`
	HTMLURL          string      `json:"html_url"`
	ID               int         `json:"id"`
	IsTemplate       bool        `json:"is_template"`
	IssueCommentURL  string      `json:"issue_comment_url"`
	IssueEventsURL   string      `json:"issue_events_url"`
	IssuesURL        string      `json:"issues_url"`
	KeysURL          string      `json:"keys_url"`
	LabelsURL        string      `json:"labels_url"`
	Language         interface{} `json:"language"`
	LanguagesURL     string      `json:"languages_url"`
	License          struct {
		Key    string `json:"key"`
		Name   string `json:"name"`
		NodeID string `json:"node_id"`
		SpdxID string `json:"spdx_id"`
		URL    string `json:"url"`
	} `json:"license"`
	MergesURL        string `json:"merges_url"`
	MilestonesURL    string `json:"milestones_url"`
	MirrorURL        string `json:"mirror_url"`
	Name             string `json:"name"`
	NetworkCount     int    `json:"network_count"`
	NodeID           string `json:"node_id"`
	NotificationsURL string `json:"notifications_url"`
	OpenIssuesCount  int    `json:"open_issues_count"`
	Owner            struct {
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
	} `json:"owner"`
	Permissions struct {
		Admin bool `json:"admin"`
		Pull  bool `json:"pull"`
		Push  bool `json:"push"`
	} `json:"permissions"`
	Private            bool        `json:"private"`
	PullsURL           string      `json:"pulls_url"`
	PushedAt           string      `json:"pushed_at"`
	ReleasesURL        string      `json:"releases_url"`
	Size               int         `json:"size"`
	SSHURL             string      `json:"ssh_url"`
	StargazersCount    int         `json:"stargazers_count"`
	StargazersURL      string      `json:"stargazers_url"`
	StatusesURL        string      `json:"statuses_url"`
	SubscribersCount   int         `json:"subscribers_count"`
	SubscribersURL     string      `json:"subscribers_url"`
	SubscriptionURL    string      `json:"subscription_url"`
	SvnURL             string      `json:"svn_url"`
	TagsURL            string      `json:"tags_url"`
	TeamsURL           string      `json:"teams_url"`
	TemplateRepository interface{} `json:"template_repository"`
	Topics             []string    `json:"topics"`
	TreesURL           string      `json:"trees_url"`
	UpdatedAt          string      `json:"updated_at"`
	URL                string      `json:"url"`
	WatchersCount      int         `json:"watchers_count"`
}

type UserReposNew struct {
	AllowMergeCommit bool   `json:"allow_merge_commit"`
	AllowRebaseMerge bool   `json:"allow_rebase_merge"`
	AllowSquashMerge bool   `json:"allow_squash_merge"`
	AutoInit         bool   `json:"auto_init"`
	Description      string `json:"description"`
	HasIssues        bool   `json:"has_issues"`
	HasProjects      bool   `json:"has_projects"`
	HasWiki          bool   `json:"has_wiki"`
	Homepage         string `json:"homepage"`
	IsTemplate       bool   `json:"is_template"`
	Name             string `json:"name"`
	Private          bool   `json:"private"`
}

func (r *UserReposNew) InitDefault(reposName, description string) {
	r.AllowMergeCommit = true
	r.AllowRebaseMerge = true
	r.AllowSquashMerge = true
	r.Description = description
	r.HasIssues = true
	r.HasProjects = true
	r.HasProjects = true
	r.HasWiki = true
	r.Homepage = "https://github.com"
	r.Name = reposName
}
