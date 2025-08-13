FROM golang:1.23-alpine

RUN apk update && apk add --no-cache build-base musl-dev ca-certificates tzdata

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOFLAGS="-mod=vendor"

WORKDIR /build
COPY go.mod .
COPY go.sum .
COPY . .


ARG COMMIT_HASH
ARG BUILD_DATE

RUN echo $COMMIT_HASH

RUN go mod vendor

# RUN go build -buildvcs=false -trimpath -a -ldflags="-s -w -X main.commitHash=${COMMIT_HASH} -X main.buildDate=${BUILD_DATE} -extldflags" -o main .
RUN go build -trimpath -a -ldflags="-s -w -X main.commitHash=$COMMIT_HASH -X main.buildDate=$BUILD_DATE" -o main .





WORKDIR /dist


RUN cp /build/main .

FROM scratch

COPY --from=0 /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /dist/main /
ENTRYPOINT ["/main"]