// getReposInfos.go

// Source file auto-generated on Thu, 26 Sep 2019 08:03:26 using Gotk3ObjHandler v1.3.8 ©2018-19 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package goGithub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

/***********************/
/* Token login github */
/*********************/
type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func GithubLogin(token string) (message string) {
	tokenSource := &TokenSource{
		AccessToken: token,
	}
	ctx := context.Background()
	oauthClient := oauth2.NewClient(ctx, tokenSource)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		message = fmt.Sprintf("client.Users.Get() faled with '%s'\n", err)
		return
	}
	d, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		message = fmt.Sprintf("json.MarshlIndent() failed with %s\n", err)
		return
	}
	message = fmt.Sprintf("User:\n%s\n", string(d))
	return
}

/***********************************************************************/
/*	Get ropositories information. https acces via personal token.     */
/*	you need to generate personal access token at:                   */
/*	https://github.com/settings/applications#personal-access-tokens */
/*	or https://github.com/settings/tokens                          */
/******************************************************************/
type GithubUserInfos struct {
	PersonalAccessToken string
	User                string
	Repositories        []repository
	Forks               []repository

	GithubClient    *github.Client
	HttpClient      *http.Client
	RepositoriesRaw []*github.Repository
}

type repository struct {
	Name          string
	NameFull      string
	Owner         string
	DefaultBranch string
	Branches      []branch
	Issues        []issue
	PullRequests  []issue
	Desc          string
	License       string
	IsFork        bool
	IssuesOpen    int
	ForksUrl      string
}

type issue struct {
	Title         string
	CreatedAt     time.Time
	ClosedAt      time.Time
	CommentsCount int
	Url           string
	User          string
	State         string
}

type branch struct {
	Commit    string
	Name      string
	Protected bool
}

// GithubUserReposList: Get information on user's repositories.
func (gui *GithubUserInfos) GithubUserReposList() (err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gui.PersonalAccessToken},
	)
	gui.HttpClient = oauth2.NewClient(ctx, ts)
	gui.GithubClient = github.NewClient(gui.HttpClient)

	// list all repositories for the authenticated user
	gui.RepositoriesRaw, _, err = gui.GithubClient.Repositories.List(ctx, "", nil)
	for _, repo := range gui.RepositoriesRaw {

		// Get repository infos
		r := repository{
			Name:          repo.GetName(),
			License:       repo.License.GetName(),
			NameFull:      repo.GetFullName(),
			Owner:         repo.Owner.GetLogin(),
			DefaultBranch: repo.GetDefaultBranch(),
			Desc:          repo.GetDescription(),
			IssuesOpen:    repo.GetOpenIssuesCount()}

		// Get branches list
		if branches, _, err := gui.GithubClient.Repositories.ListBranches(ctx, repo.GetOwner().GetLogin(), repo.GetName(), nil); err == nil {
			// if branches, _, err := gui.GithubClient.Repositories.ListBranches(ctx, "gotk3", "gotk3", nil); err == nil {

			// Store branch
			for _, b := range branches {
				r.Branches = append(r.Branches, branch{
					Name:      b.GetName(),
					Protected: b.GetProtected()})
			}

			// Get issues list
			if issues, _, err := gui.GithubClient.Issues.ListByRepo(ctx, repo.GetOwner().GetLogin(), repo.GetName(), nil); err == nil {
				// if issues, _, err := gui.GithubClient.Issues.ListByRepo(ctx, "gotk3", "gotk3", nil); err == nil {

				// Store issue
				for _, i := range issues {

					// Only opened ones
					if i.ClosedAt != nil {
						oneIssue := issue{
							Title:         i.GetTitle(),
							CreatedAt:     i.GetCreatedAt(),
							ClosedAt:      i.GetClosedAt(),
							CommentsCount: i.GetComments(),
							Url:           i.GetHTMLURL(),
							User:          i.GetUser().GetLogin(),
							State:         i.GetState(),
						}

						// Dispaching issue/pull request
						if i.IsPullRequest() {
							r.PullRequests = append(r.PullRequests, oneIssue)
						} else {
							r.Issues = append(r.Issues, oneIssue)
						}
					}
				}
			}

			// Add Repository or Fork to main struct
			if repo.GetFork() {
				r.ForksUrl = repo.GetForksURL()
				gui.Forks = append(gui.Forks, r)
			} else {
				gui.Repositories = append(gui.Repositories, r)
			}
		}
	}
	if err != nil {
		log.Fatal(err.Error())
	}
	return
}
