FROM docker.io/library/golang:1.17-alpine
WORKDIR /go/src/github.com/rodneyxr/ffatoolkit
COPY . .
RUN go get -d -v ./... && \
    go install -v ./...
CMD ["ffatoolkit"]