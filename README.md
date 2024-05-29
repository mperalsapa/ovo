# OVO
Repository of OVO project.

# Deployment
## Docker
For a demo service deploy we use a docker image built from our git server using CI/CD Pipelines.

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