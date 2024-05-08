FROM golang:1.22-alpine AS build
RUN apk update && apk add --no-cache curl make git
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /src
COPY . .

RUN go mod download
RUN go build -gcflags="all=-N -l" -o app ./cmd/app

FROM alpine:latest
COPY --from=build /go/bin/dlv /usr/local/bin/dlv
WORKDIR /src
COPY --from=build /src/app .
COPY --from=build /usr/local/go/lib/time/zoneinfo.zip /
ENV ZONEINFO=/zoneinfo.zip

ENTRYPOINT ["dlv"]
