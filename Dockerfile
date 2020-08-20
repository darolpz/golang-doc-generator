FROM golang:alpine

RUN apk update && \
  apk upgrade && \
  apk add git npm

# Set the application directory
WORKDIR $GOPATH/src/github.com/darolpz/golang-doc-generator

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# This container exposes port 4008 to the outside world
EXPOSE 4008

# Run the executable
CMD ["golang-doc-generator"]

