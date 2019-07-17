FROM golang:1.12 as build

ENV GOPROXY https://goproxy.io
ENV GO111MODULE on

LABEL maintainer="Mulkids <jiangzx0526@gmail.com>"

ARG APP_NAME=go-docker
ARG LOG_DIR=/${APP_NAME}/logs

RUN mkdir -p ${LOG_DIR}

ENV LOG_FILE_LOCATION=${LOG_DIR}/app.log

WORKDIR /go/web/cache

ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /app/mulkids

COPY . .

CMD ["go", "run", "main.go"]

EXPOSE 8080
