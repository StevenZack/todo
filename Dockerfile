FROM golang:1.21.4 AS builder

RUN go version

ARG PROJECT_VERSION

COPY . /go/src/github.com/StevenZack/todo
WORKDIR /go/src/github.com/StevenZack/todo

RUN go mod tidy

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o app .

RUN go test -cover -v ./...

FROM zigzigcheers/upx as minify
WORKDIR /root
COPY --from=builder /go/src/github.com/StevenZack/todo/app .
RUN upx ./app

FROM scratch
WORKDIR /root/
COPY --from=minify /root/app . 
EXPOSE 8080
ENTRYPOINT [ "./app" ]