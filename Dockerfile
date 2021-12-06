FROM golang:alpine
RUN apk add bash ca-certificates git gcc g++ libc-dev
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod download
RUN go build -o main .
EXPOSE 8081


