FROM golang:1.12-alpine AS builder
RUN apk update && apk add git ca-certificates net-tools tzdata ngrep
WORKDIR /go/src/github.com/Arnold1/ai-app
ENV GO111MODULE=on
COPY go.mod .
RUN go mod download
COPY . .
WORKDIR /go/src/github.com/Arnold1/ai-app/aiappd
RUN pwd

ARG BUILD_TIME
ARG VERSION
ARG BUILD_SHA

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X ${BUILD_TIME} -X ${VERSION} -X ${BUILD_SHA}" -o bin/aiappd

FROM alpine:3.7
RUN apk update && apk upgrade \
    && apk add bash \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/*
COPY --from=builder /go/src/github.com/Arnold1/ai-app/aiappd/bin/aiappd /apps/aiappd/bin/aiappd
COPY --from=builder /go/src/github.com/Arnold1/ai-app/aiappd/view /apps/aiappd/bin/view
WORKDIR /apps/aiappd/bin
ENTRYPOINT ["./aiappd"]