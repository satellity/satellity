FROM golang:1.12.1 AS build
LABEL authors="Li Yuqing <im.yuqlee@gmail.com>, Guo Huang <guohuang@gmail.com>, Marat Fayzullin <fay@zull.in>"

WORKDIR /godiscourse

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN make production

FROM alpine:latest AS runtime
COPY --from=build /godiscourse/bin/godiscourse /api/
EXPOSE 8080
ENTRYPOINT ["/api/godiscourse"]