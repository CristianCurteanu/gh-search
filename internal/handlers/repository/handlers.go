package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/CristianCurteanu/gh-search/internal/middlewares"
	"github.com/CristianCurteanu/gh-search/pkg/githubapi"
)

type RepositoriesHandlers struct {
	githubClient *githubapi.GithubApi
	middlewares.UseMiddleware
}

func NewRepositoriesHandlers(githubClient *githubapi.GithubApi) *RepositoriesHandlers {
	return &RepositoriesHandlers{
		githubClient: githubClient,
	}
}

func (ph *RepositoriesHandlers) Search(w http.ResponseWriter, r *http.Request) {
	ph.Handle(w, r, func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("access_token")
		if err != nil {
			panic(err)
		}
		// queryString := "org%3Amoovweb+gvm&type=repositories"
		// queryString := "user%3ACristianCurteanu+go-boggle-solver"

		queryString := r.URL.Query().Get("query")
		if queryString == "" {
			// TODO: HANDLE WITH PROPER RESULT TEMPLATE
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			w.Write([]byte("{\"error\":\"missing query string\"}"))
			return
		}

		page := r.URL.Query().Get("page")
		if page == "" {
			page = "1"
		}

		params := url.Values{}
		params.Add("q", queryString)
		params.Add("type", "repositories")
		params.Add("page", page)

		queryString = params.Encode()
		fmt.Printf("queryString(): %v\n", queryString)

		githubData, err := ph.githubClient.SearchRepository(token.Value, queryString)
		if err != nil {
			// TODO: HANDLE WITH PROPER RESULT TEMPLATE
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			log.Printf("failed to get repositories, due to err=%q", err)
			w.Write([]byte("{\"error\":\"something went wrong when requesting search\"}"))
			return
		}
		// Set return type JSON
		w.Header().Set("Content-type", "application/json")

		// Prettifying the json
		// var prettyJSON bytes.Buffer
		// // json.indent is a library utility function to prettify JSON indentation
		// parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
		// if parserr != nil {
		// 	log.Panic("JSON parse error")
		// }

		// Return the prettified JSON as a string
		ghDataJson, _ := json.Marshal(githubData)

		fmt.Fprint(w, string(ghDataJson))
	})
}

func (ph *RepositoriesHandlers) GetRepositoryPage(w http.ResponseWriter, r *http.Request) {
	ph.Handle(w, r, func(w http.ResponseWriter, r *http.Request) {})
}
