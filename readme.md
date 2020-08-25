# Golang-doc-generator

  Golang-doc-generator is a PDF generator written in Go. It genere a release notes pdf from git commit messages.
  Important: It do needs commits writen using [convenctional commits](https://www.conventionalcommits.org/en/v1.0.0/) in order to work



### Installation

golang-doc-generator requires [Golang](https://golang.org/) to run.

Install golang-doc-generator in your GO PATH and install the dependencies using go get.

```sh
 cd $GOPATH/src/github.com/darolpz/golang-doc-generator
 go get -u
```

Then create a .env file to set environment variables. There is a .env.example in the root of this project.

```sh
GITLAB_URL=A VALID GITLAB ROUTE
PORT=PORT WHERE THE SERVICE WILL RUN

SLACK_URL=A URL TO A SLACK WORKSPACE WHERE YOU WANT TO NOTIFY
SLACK_APP_TOKEN=TOKEN OF YOUR APP
HOST_URL= THE URL OF YOUR SERVICE

CHANNEL1=CHANNEL WHERE YOU WANT TO SEND A NOTIFICATION
CHANNEL2=CHANNEL WHERE YOU WANT TO SEND A NOTIFICATION
CHANNEL3=CHANNEL WHERE YOU WANT TO SEND A NOTIFICATION
```

Start the server using

```sh
  go run main.go
```

OR 

```sh
go build . 
sh golang-doc-generator
```

OR 

```sh
go install .
golang-doc-generator
```



### Docker
golang-doc-generator is very easy to install and deploy in a Docker container.

By default, the Docker will expose port 4008, so change this within the Dockerfile if necessary. When ready, simply use the Dockerfile to build the image.

```sh
cd $GOPATH/src/github.com/darolpz/golang-doc-generator
docker build . -t golang-doc-generator .
```
This will create the golang-doc-generator image and pull in the necessary dependencies.

Once done, run the Docker image and map the port to whatever you wish on your host. In this example, we simply map port 4008 of the host to port 4008 of the Docker (or whatever port was exposed in the Dockerfile):

```sh
docker run -d -p 4008:4008 --name=golang-doc-generator golang-doc-generator
```

Verify the deployment by navigating to your server address in your preferred browser.

```sh
127.0.0.1:4008/ping
```
