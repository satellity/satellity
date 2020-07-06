FROM golang:alpine as builder

WORKDIR /go/src/app

COPY . .

RUN mv internal/configs/config.example internal/configs/config.yaml \
    && go build

RUN mkdir -p /go/build/app/internal/configs \
    && cp internal/configs/config.yaml      /go/build/app/internal/configs/ \
    && cp satellity                         /go/build/app/

FROM golang:alpine

COPY --from=builder /go/build/app /go/build/app

ENTRYPOINT [ "/go/build/app/satellity" ]
