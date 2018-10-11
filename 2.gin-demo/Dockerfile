FROM golang:latest
#FROM scratch

WORKDIR $GOPATH/src/demo/2.gin-demo
COPY . $GOPATH/src/demo/2.gin-demo
RUN go build -o blog .

EXPOSE 8000
ENTRYPOINT ["./blog"]
