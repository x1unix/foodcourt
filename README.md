# Foodcourt
> Food order app for Limelight Lviv office

## Documentation
Api docs are available at [Postman Documenter](https://documenter.getpostman.com/view/1108844/voracity/77mZMFv). In the future will be moved to Swagger

## Development

### Files

* `fc-server.go` - Web-server
* `fc-run.go` - Command runner

### Installation

1. Install [Go compiler](https://golang.org/) (v1.9+)
2. Install [Docker](https://www.docker.com/)
3. Install [dep](https://github.com/golang/dep) dependency manager
4. Define `GOPATH` variable (if not defined)
5. Clone this project to `$GOPATH/src/foodcourt`
6. Install all dependencies using Glide: `dep ensure`

**Attention** - Ensure that you've defined the `GOPATH` [global variable](https://github.com/golang/go/wiki/GOPATH).


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