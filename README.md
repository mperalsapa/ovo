# OVO
Repository of the OVO project.

This software, aims to provide an easy to use and easy to host media content library.
It expects to have a list of directories containing video files, named by the content, and provides a web application to interact with.

The main features are:
- Syncplay
    Sync multiple video players playback.
- Multifolder libraries
    Joins multiple directories as a single one
- Easy to deploy
    Runtime can be provided over a binary file or a docker container. More on [Deployment section](README.md#deployment)

# Dev
To run this project as a developer follow this steps:
- Install GO (current version for this project is 1.22.1).
    - Install AIR for a live binary compilation on changes
    - Install Templ CLI for template generation
- Install Node.JS (or equivalent) in order to fetch packages from NPM
- Install SASS is in system's PATH
- Install FFmpeg is in system's PATH
- Prepare any demo video for testing purposes
- Start the IDE and run the following commands
    - In one terminal, run the backend command
        ```bash
        air
        ```
        This will watch files that end in `.go` and when they change the app will recompile itself.
        The alternative would be the `go run` command.
        
    - In a second terminal, run the frontend related command
        ```bash 
        npm run watch:css
        ```
        This will watch any `.scss` file and will recompile using SASS & Tailwind.

By default, a `.env` file is required. When that `.env` file is not present, the app will try to get environment variables. If environment variables are still missing, the app will shutdown itself. For example, a TMDB API key is required.

Lastly, the default database is SQLite. If you plan to use another database supported by GORM (MySQL, PosgreSQL, etc...), you need to specify the connection parameters in the `.env` file or as environment variables.

## TODO
Explain frontend submodule
```bash
git clone <URL_REPOSITORIO_BACKEND>
cd go-backend
git submodule init
git submodule update
```
And frontend build
```bash
cd frontend
npm install
npm run dev
```



# Deployment
## Binary
An way of running the server, could be using the binary file. This eliminates docker from the equation, but requires to have FFmpeg in system's path.

A binary file can be found here: LINK (Cooming soon)

## Docker
For sure the easiest way is Docker. Although it can consume a bit more resources, the difference could be negligible given that the software stack is minimal and the memory impact is really low.

### Build
The build process is automated by the [Dockerfile](Dockerfile). To build the image, run the following command:
```bash
docker build -t ovo-server .
```
### Run
In order to run the service, using a simple docker command could do the trick, but given that it is recommended to use a database instead of SQLite we provide a simple [docker compose](docker-compose.yaml).

Simply start the stack with the following command
```bash
docker compose up -d
```