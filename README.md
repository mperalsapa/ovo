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