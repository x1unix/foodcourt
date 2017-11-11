# Foodcourt
> Food order app for Limelight Lviv office

## Documentation
Api docs are available at [Postman Documenter](https://documenter.getpostman.com/view/1108844/voracity/77mZMFv). In the future will be moved to Swagger

## Development
### Installation
1. Install Go (v1.9+)
2. Install [Docker](https://www.docker.com/)
3. Clone this project
4. Open terminal and go to folder with this project
5. Bootstrap docker environment with `docker-compose up -d` (once)
6. Install all dependencies from `dependencies.txt` using `go get ...`

### Running
1. Open terminal and navigate to the project
2. Start docker containers with `docker-compose start`
3. Start web-server using `go run fc-server.go` (or compile it using `go build`)

**Tip** - To shutdown Docker containers, use `docker-compose stop`
