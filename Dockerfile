#build stage
FROM golang:alpine AS builder
RUN apk add build-base
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a *.go

FROM mariadb:10.5.9-focal
ARG TARGETPLATFORM

RUN apt-get update

RUN apt install curl -y

# Install kubectl binary
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/${TARGETPLATFORM}/kubectl"
RUN install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# Install MongoDB tools
RUN if [ $TARGETPLATFORM == "linux/amd64" ]; then curl https://fastdl.mongodb.org/tools/db/mongodb-database-tools-ubuntu2004-x86_64-100.3.1.deb --output mongodb-tools.deb; else curl https://fastdl.mongodb.org/tools/db/mongodb-database-tools-ubuntu2004-arm64-100.3.1.deb --output mongodb-tools.deb; fi

RUN apt install ./mongodb-tools.deb
RUN rm ./mongodb-tools.deb

COPY --from=builder /go/src/app/main /main
COPY --from=builder /go/src/app/app/templates /app/templates
COPY --from=builder /go/src/app/static /static

ENTRYPOINT /main
