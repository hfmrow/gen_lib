// gitMainStruct.go

// Source file auto-generated on Fri, 27 Sep 2019 19:41:41 using Gotk3ObjHandler v1.3.8 ©2018-19 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package github

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/hfmrow/gen_lib/tools/gojson"
	"golang.org/x/oauth2"
)

type GitMainStruct struct {
	Login               string
	Oauth2Client        *http.Client
	PersonalAccessToken string
	GoStruct            bool

	// Structures
	Branches *GitBranches
	Releases *GitReleases
	Search   *GitSearch
	Repos    *GitUsersRepos
}

func (g *GitMainStruct) GitMainStructInit(login, token string) {
	g.Login = login
	g.PersonalAccessToken = token

	// Init status codes list
	statusCodesList = generateStatusCode()

	// Init oauth2 client
	if len(token) > 0 {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: g.PersonalAccessToken})
		g.Oauth2Client = oauth2.NewClient(ctx, ts)
	}
	// Init Structures
	g.Branches = new(GitBranches)
	g.Branches.oauth2Client = g.Oauth2Client
	g.Branches.Login = g.Login

	g.Releases = new(GitReleases)
	g.Releases.oauth2Client = g.Oauth2Client
	g.Releases.Login = g.Login

	g.Search = new(GitSearch)
	g.Search.oauth2Client = g.Oauth2Client
	g.Search.Login = g.Login

	g.Repos = new(GitUsersRepos)
	g.Repos.oauth2Client = g.Oauth2Client
	g.Repos.Login = g.Login
}

/***************/
/* Misc funct */
/*************/
var longStatusCodeMessage = true

// Build output display with status code.
func getStatus(context string, sCode int, longDesc ...bool) (out string, err error) {
	var found bool
	if len(longDesc) > 0 {
		longStatusCodeMessage = longDesc[0]
	}
	for _, val := range *statusCodesList { // Status code parsing
		if val.Code == fmt.Sprintf("%d", sCode) && val.Context == context {
			found = true
			if longStatusCodeMessage {
				out += val.MessLong
			} else {
				out += val.MessShort
			}
		}
	}
	if !found { // not found: return generic code
		out += http.StatusText(sCode)
	}
	return
}

// Get:
func Get(client *http.Client, url string) (out []byte, err error) {
	var resp *http.Response
	if client == nil {
		client = new(http.Client)
	}
	if resp, err = client.Get(url); err == nil {
		defer resp.Body.Close()
		out, err = ioutil.ReadAll(resp.Body)
	}
	return
}

// Post:
func Post(client *http.Client, url string, jsonBytes []byte) (out []byte, err error) {
	var resp *http.Response
	r := bytes.NewReader(jsonBytes)

	if client == nil {
		client = new(http.Client)
	}

	if resp, err = client.Post(url, "application/json", r); err == nil {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return out, errors.New(fmt.Sprintf("Unexpected status code, %d\n", resp.StatusCode))
		}
		out, err = ioutil.ReadAll(resp.Body)
	}
	return
}

// GetOauthInfo: get info like scopes for a client's token
func GetOauthInfo(client *http.Client) (out []byte, err error) {
	return Get(client, "https://api.github.com/authorizations") // TODO does not work !!!
}

type objmap map[string]interface{}

func JsonToMap(body io.ReadCloser) (outStr string, err error) {
	var data []byte
	var obj objmap
	if data, err = ioutil.ReadAll(body); err == nil {
		if err = json.Unmarshal(data, &obj); err == nil {
			mapIt(obj)
		}
	}
	return
}

func mapIt(subObj objmap) {
	for key, value := range subObj {
		switch value.(type) {
		case map[string]interface{}:
			mapIt(value.(objmap))
		case map[interface{}]interface{}:
			mapIt(value.(objmap))
		default:
			fmt.Printf("%v = %v\n", key, value)
		}
	}
}

// JsonToStruct: convert json to go structure.
func JsonToStruct(data []byte) (out string, err error) {
	outByte, err := gojson.JsonToStruct(data, "myStruct", "", false, true)
	return string(outByte), err
}

func SaveJsonToStruct(body io.ReadCloser) (err error) {
	var outStr string
	var data []byte
	if data, err = ioutil.ReadAll(body); err == nil {
		if outStr, err = JsonToStruct(data); err == nil {
			err = ioutil.WriteFile("jsonStruct.txt", []byte(outStr), os.ModePerm)
		}
	}
	return
}

// writeRawFile: record to file row information recieved from GET command
// used in developpment to get json data and convert to go struct
func WriteRawFile(resp *http.Response) (err error) {
	var f *os.File
	if f, err = os.Create("jsonRaw.txt"); err == nil {
		defer f.Close()
		w := bufio.NewWriter(f)
		resp.Write(w)
		// resp.Header.Write(w)
		w.Flush()
	}
	return
}

/****************************/
/* Build status codes list */
/**************************/
var genCodesList = [][]string{
	{"400", "BAD REQUEST", "The request was invalid or cannot be otherwise served. An accompanying error message will explain further. For security reasons, requests without authentication are considered invalid and will yield this response."},
	{"401", "UNAUTHORIZED", "The authentication credentials are missing, or if supplied are not valid or not sufficient to access the resource."},
	{"403", "FORBIDDEN", "The request has been refused. See the accompanying message for the specific reason (most likely for exceeding rate limit)."},
	{"404", "NOT FOUND", "The URI requested is invalid or the resource requested does not exists."},
	{"406", "NOT ACCEPTABLE", "The request specified an invalid format."},
	{"410", "GONE", "This resource is gone. Used to indicate that an API endpoint has been turned off."},
	{"429", "TOO MANY REQUESTS", "Returned when a request cannot be served due to the application’s rate limit having been exhausted for the resource."},
	{"500", "INTERNAL SERVER ERROR", "Something is horribly wrong."},
	{"502", "BAD GATEWAY", "The service is down or being upgraded. Try again later."},
	{"503", "SERVICE UNAVAILABLE", "The service is up, but overloaded with requests. Try again later."},
	{"504", "GATEWAY TIMEOUT", "Servers are up, but the request couldn’t be serviced due to some failure within our stack. Try again later."}}
var getCodesList = [][]string{
	{"200", "OK", "The request was successful and the response body contains the representation requested."},
	{"302", "FOUND", "A common redirect response; you can GET the representation at the URI in the Location response header."},
	{"304", "NOT MODIFIED", "There is no new data to return."}}
var putCodesList = [][]string{
	{"201", "OK", "The request was successful, we updated the resource and the response body contains the representation."},
	{"202", "ACCEPTED", "The request has been accepted for further processing, which will be completed sometime later."}}
var delCodesList = [][]string{
	{"202", "ACCEPTED", "The request has been accepted for further processing, which will be completed sometime later."},
	{"204", "OK", "The request was successful; the resource was deleted."}}

var statusCodesList *statusCodes

type statusCodes []sCode

type sCode struct {
	Context   string
	Code      string
	MessShort string
	MessLong  string
}

func generateStatusCode() (out *statusCodes) {
	out = new(statusCodes)
	var addList = func(ctx string, list [][]string) (out statusCodes) {
		for _, val := range list {
			var code sCode
			code.Context = ctx
			code.Code = val[0]
			code.MessShort = val[1]
			code.MessLong = val[2]
			out = append(out, code)
		}
		return
	}
	*out = append(*out, addList("gen", genCodesList)...)
	*out = append(*out, addList("get", getCodesList)...)
	*out = append(*out, addList("put", putCodesList)...)
	*out = append(*out, addList("del", delCodesList)...)
	return
}

/**********************************************************/
/* Display list of repository with more than 10000 stars */
/********************************************************/
// owner: is the repository owner
type owner struct {
	Login string
}

// item: is the single repository data structure
type item struct {
	ID              int
	Name            string
	FullName        string `json:"full_name"`
	Owner           owner
	Description     string
	CreatedAt       string `json:"created_at"`
	StargazersCount int    `json:"stargazers_count"`
}

// jSONData contains the GitHub API response
type jSONData struct {
	Count int `json:"total_count"`
	Items []item
}

func search() {
	var err error
	var resp *http.Response
	var body []byte
	if resp, err = http.Get("https://api.github.com/search/repositories?q=stars:>=10000+language:go&sort=stars&order=desc"); err == nil {
		defer resp.Body.Close()
		if body, err = ioutil.ReadAll(resp.Body); err == nil {
			if resp.StatusCode != http.StatusOK {
				log.Fatal("Unexpected status code", resp.StatusCode)
				return
			}
			data := jSONData{}
			if err = json.Unmarshal(body, &data); err == nil {
				printData(data)
			}
		}
	}
	if err != nil {
		log.Fatal(err.Error())
	}
}

func printData(data jSONData) {
	log.Printf("Repositories found: %d", data.Count)
	const format = "%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Repository", "Stars", "Created at", "Description")
	fmt.Fprintf(tw, format, "----------", "-----", "----------", "----------")
	for _, i := range data.Items {
		desc := i.Description
		if len(desc) > 50 {
			desc = string(desc[:50]) + "..."
		}
		t, err := time.Parse(time.RFC3339, i.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(tw, format, i.FullName, i.StargazersCount, t.Year(), desc)
	}
	tw.Flush()
}
