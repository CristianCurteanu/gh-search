package githubapi

import "time"

type RepositorySearchResult struct {
	Total int           `json:"total_count"`
	Items []*Repository `json:"items"`
}

type ProfileData struct {
	Id        int    `json:"id"`
	Username  string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	Company   string `json:"company"`
	Repos     int    `json:"public_repos"`
	Gists     int    `json:"public_gists"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
}

type Repository struct {
	Id          int          `json:"id"`
	NodeID      string       `json:"node_id"`
	Name        string       `json:"name"`
	Url         string       `json:"url"`
	HtmlUrl     string       `json:"html_url"`
	Owner       *ProfileData `json:"owner"`
	FullName    string       `json:"full_name"`
	IsPrivate   bool         `json:"private"`
	IsFork      bool         `json:"fork"`
	Description string       `json:"description"`
	Size        int          `json:"size"`
	Stars       int          `json:"stargazers_count"`
	Watchers    int          `json:"watchers_count"`
	Forks       int          `json:"forks_count"`
	Language    string       `json:"language"`
	Issues      int          `json:"open_issues_count"`
	UpdatedAt   *time.Time   `json:"updated_at"`
	PushedAt    *time.Time   `json:"pushed_at"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type Contributors []Contributor

type Contributor struct {
	Id        int    `json:"id"`
	Username  string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HtmlUrl   string `json:"html_url"`
}

type Commits []Commit

type Commit struct {
	Url      string         `json:"url1"`
	HTMLUrl  string         `json:"html_url"`
	Sha      string         `json:"sha"`
	Author   *Contributor   `json:"author"`
	Commiter *Contributor   `json:"committer"`
	Commit   *CommitDetails `json:"commit"`
}

type CommitDetails struct {
	Author    *CommitAuthor `json:"author"`
	Committer *CommitAuthor `json:"commiter"`
	Message   string        `json:"message"`
	Comments  int           `json:"comment_count"`
}

type CommitAuthor struct {
	Name  string     `json:"name"`
	Email string     `json:"email"`
	Date  *time.Time `json:"date"`
}
