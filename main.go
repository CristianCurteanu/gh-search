package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	authentication "github.com/CristianCurteanu/gh-search/internal/auth"
	"github.com/CristianCurteanu/gh-search/internal/handlers/auth"
	"github.com/CristianCurteanu/gh-search/internal/handlers/profile"
	"github.com/CristianCurteanu/gh-search/internal/handlers/repository"
	"github.com/CristianCurteanu/gh-search/internal/middlewares"
	"github.com/CristianCurteanu/gh-search/pkg/cache"
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

	redisSessionStorage := authentication.NewRedisSessionStorage(
		newRedisClient(authentication.RedisDBSessionStorage),
	)
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

	dataCache := cache.NewRedisCache(newRedisClient(2))
	profileHandlers := profile.NewProfileHandlers(
		githubClient,
		dataCache,
	)
	profileHandlers.Use(requestLog)
	profileHandlers.Use(sessionMiddleware)

	mux.HandleFunc("/profile", profileHandlers.GetProfilePage)
	mux.HandleFunc("/logout", profileHandlers.Logout)

	repositories := repository.NewRepositoriesHandlers(githubClient, dataCache)
	repositories.Use(requestLog)
	repositories.Use(sessionMiddleware)

	mux.HandleFunc("/search", repositories.Search)
	mux.HandleFunc("/repository", repositories.GetRepositoryPage)

	fmt.Println("[ UP ON PORT 3000 ]")
	log.Panic(
		http.ListenAndServe(":3000", mux),
	)
}

func newRedisClient(db int) *redis.Client {
	redisHost, found := os.LookupEnv("REDIS_HOST")
	if !found {
		panic("REDIS_HOST not defined")
	}

	redisPort, found := os.LookupEnv("REDIS_PORT")
	if !found {
		panic("REDIS_PORT not defined")
	}

	redisPassword, found := os.LookupEnv("REDIS_PASSWORD")
	if !found {
		panic("REDIS_PASSWORD not defined")
	}

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	return redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       db,
	})
}
