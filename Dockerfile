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

ENV GITLAB_URL=http://10.200.172.71
ENV PORT=4008
ENV SLACK_URL=https://devsargentinos.slack.com/api/files.upload
ENV SLACK_APP_TOKEN=xoxb-809980440567-1321592942178-XVnMfDxZiXe9N3TIkgCVRisM
ENV HOST_URL=http://localhost:4008

ENV DARO=UPRKH2ZQC
# Run the executable
CMD ["golang-doc-generator"]

