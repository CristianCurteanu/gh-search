## Github Repository Search tool

This web application, written in Go, enables users to search for GitHub repositories and retrieve detailed information about them. 

It leverages the GitHub API to fetch repository data such as description, stars, and forks, as well as providing an overview of the latest commits and contributors. You can easily explore and analyze repositories with a clean, responsive interface that enhances the browsing experience.

### Installation

There are two different options to launch this application:
- Using docker-compose
- Standalone build

#### Using docker compose

There is already a defined docker compose in this repo, and makefile tasks have been created to simplify the run.

Before run, please make sure to up `GITHUB_ID`, `GITHUB_SECRET` and `GITHUB_REDIRECT_URL` to a `.env` (take a look at `.env.example`) file, like this: 

```
GITHUB_ID=xxxxx
GITHUB_SECRET=xxxxxxxxxxxxxxxxxxxxxxxx
GITHUB_REDIRECT_URL=http://localhost:3000/auth/callback/success
```

and run source command:

```shell
$ source .env
```

After, in order to start, following commands should be executed:

```shell
$ make dc-build dc-run
```

This will make the images for server app and redis service build, and run in an inline mode.

If you want to run the containers in a daemon mode, you can use this command:

```shell
$ make dc-rund
```

If you run in daemon mode, and you want to stop all the containers, run the following command:

```shell
$ make dc-stop
```

This will kill and remove all the containers

#### Using standalone build

If you do this way, make sure that you Redis DB is running on your machine. Also, don't forget to add these environment variables before `run`:

```shell
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=<whatever password is set, or you can use false if password not set>
GITHUB_ID=xxxxx
GITHUB_SECRET=xxxxxxxxxxxxxxxxxxxxxxxx
GITHUB_REDIRECT_URL=http://localhost:3000/auth/callback/success
```

Then, you would want to build and run the app manually, you can run following commands:

```shell
$ source .env
$ make build run
```

This will create a `runner` binary inside the folder, and run it.

### Development

The entrypoint in this application is `main.go` file located at the root of the project folder.

There are two other folders, separated by scope:
- `internal` - all internal logic of the app is located, all the packages from here are not supposed to be imported externally
- `pkg` - here are components that are used inside the application, and could be re-used in other projects as well, such like Github API and Cache

The internal folder contains:
- `handlers` where functionality for each logic module of the app is located, and all the related functionality to that module as well, like service, tests, templ pages
- `middlewares` where are located middleware components that could be injected in handlers structs
- `auth` - stores the components for authentication logic, like JWT encoders and session storage logic
- `layouts` - stores all the re-usable Templ templates, like global layout, profile layout and so on. All other module related templates are stored in dedicated `pages` package at the handlers level

### Known issues
- Test coverage is not fully done
- There is a slight quirk on repository details page, when there are scrollable commits, but not scrollable contributors
- Handlers could be less coupled from Services (like in the repository package)