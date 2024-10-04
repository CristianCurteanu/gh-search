package githubapi

type RepositorySearchResult struct {
	Total int          `json:"total_count"`
	Items []*RepoItems `json:"items"`
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

type RepoItems struct {
	Id          int          `json:"id"`
	NodeID      string       `json:"node_id"`
	Name        string       `json:"name"`
	Url         string       `json:"url"`
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
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}
