# Foodcourt
> Food order app for Limelight Lviv office

## Documentation
Api docs are available at [Postman Documenter](https://documenter.getpostman.com/view/1108844/voracity/77mZMFv). In the future will be moved to Swagger

## Development
### Installation
#### Backend (API)
1. Install [Go compiler](https://golang.org/) (v1.9+)
2. Install [Docker](https://www.docker.com/)
3. Install [Glide](https://glide.sh/) package manager
4. Clone this project
5. Open terminal and go to folder with this project
6. Bootstrap docker environment with `docker-compose up -d` (once)
7. Install all dependencies using Glide: `glide install`

**Attention** - Ensure that you've defined the `GOPATH` [global variable](https://github.com/golang/go/wiki/GOPATH).

#### Frontend (UI)
1. Install [Node.js](https://nodejs.org/en/) (LTS recommended)
2. Install [Yarn](https://yarnpkg.com/en/)
3. Open the terminal and navigate to the `src/ui` folder
4. Install all dependencies using `yarn install`
5. Build Angular project using `yarn build` command

### Running
1. Open terminal and navigate to the project
2. Start docker containers with `docker-compose start`
3. Start web-server using `go run fc-server.go` (or compile it using `go build`)

**Tip** - To shutdown Docker containers, use `docker-compose stop`

### Dev-environment and demo data
By default, web-server starts at `http://localhost:8000` and configuration is provided at `.env` file.

Default credentials are:
* **Email** - `admin@llnw.com`
* **Password** - `q1w2e3`

Default demo data is available at `demo` folder.

To import demo data, run `import.sh` script at `demo` folder:
```bash
$ cd voracity-ua
$ cd demo
$ bash ./import.sh
* Importing dishes list ...
* Deleting existing photos ...
* Extracting photos ...
Archive:  photos.zip
   creating: ../src/public/photos/
  inflating: ../src/public/photos/06d5d616179f99d45fb014141e953437.jpg
  ...
  ...
===== IMPORT FINISH ======
```

**Tip** - You need to import demo data only *after* the UI part was build. During UI build, `public` directory will be overwriten and all photos will be removed.
