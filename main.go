package main

import (
	"fmt"
	"log"
	"net/http"

	authentication "github.com/CristianCurteanu/gh-search/internal/auth"
	"github.com/CristianCurteanu/gh-search/internal/handlers/auth"
	"github.com/CristianCurteanu/gh-search/internal/handlers/profile"
	"github.com/CristianCurteanu/gh-search/internal/handlers/repository"
	"github.com/CristianCurteanu/gh-search/internal/middlewares"
	"github.com/CristianCurteanu/gh-search/pkg/githubapi"
	"github.com/redis/go-redis/v9"
)

const (
	githubClientID     = "Ov23liXmz8CFjEWHmlo8"
	githubClientSecret = "0d7c1c44965462927f7ba2a223877c373db4800f"
)

func main() {
	mux := http.NewServeMux()
	githubClient := githubapi.NewGithubClient(githubClientID, githubClientSecret)

	redisSessionStorage := authentication.NewRedisSessionStorage(redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}))
	signer := authentication.NewJWTAuth(githubClientSecret)

	authHandlers := auth.NewAuthHandlers(
		auth.AuthHandlersConfig{
			ClientId:     githubClientID,
			ClientSecret: githubClientSecret,
			RedirectUrl:  "http://localhost:3000/auth/callback/success",
		},
		githubClient,
		redisSessionStorage,
		signer,
	)
	requestLog := middlewares.NewRequestLog()
	authHandlers.Use(requestLog)

	mux.HandleFunc("/", authHandlers.RootHandler)
	mux.HandleFunc("/login/github/", authHandlers.GithubLoginHandler)
	mux.HandleFunc("/auth/callback/success", authHandlers.GithubCallbackHandler())

	sessionMiddleware := middlewares.NewCookieSessionHandler(
		redisSessionStorage,
		signer,
	)

	profileHandlers := profile.NewProfileHandlers(githubClient)
	profileHandlers.Use(requestLog)
	profileHandlers.Use(sessionMiddleware)
	// Route where the authenticated user is redirected to
	mux.HandleFunc("/profile", profileHandlers.GetProfilePage)
	mux.HandleFunc("/logout", profileHandlers.Logout)

	repositories := repository.NewRepositoriesHandlers(githubClient)
	profileHandlers.Use(requestLog)
	repositories.Use(sessionMiddleware)
	mux.HandleFunc("/search", repositories.Search)

	mux.HandleFunc("/repository", repositories.GetRepositoryPage)

	fmt.Println("[ UP ON PORT 3000 ]")
	log.Panic(
		http.ListenAndServe(":3000", mux),
	)
}
