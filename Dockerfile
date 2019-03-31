FROM golang:1.12.1 AS build
LABEL authors="Li Yuqing <im.yuqlee@gmail.com>, Guo Huang <guohuang@gmail.com>, Marat Fayzullin <fay@zull.in>"

WORKDIR /api

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o godiscourse

FROM alpine:latest AS runtime
COPY --from=build /api/godiscourse /api/
EXPOSE 8080
ENTRYPOINT ["/api/godiscourse"]