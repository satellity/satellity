# Build app within the base builder image
FROM golang:1.11
LABEL maintainer="Guo Huang <guohuang@gmail.com>"

WORKDIR /src/godiscourse/api

# Copy solution
COPY . .

# Publish application for release
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["godiscourse"]