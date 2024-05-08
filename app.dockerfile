ARG GIT_COMMIT
ARG VERSION
ARG PROJECT

FROM golang:1.22-alpine as app-builder
RUN apk update && apk add --no-cache curl make git

WORKDIR /src

COPY . .
RUN go build ./cmd/app

FROM alpine:latest
RUN apk update && apk add --no-cache curl
WORKDIR /src
COPY --from=app-builder /src/app .

CMD ["./app"]
