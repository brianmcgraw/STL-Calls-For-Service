# FROM golang:alpine as builder

# WORKDIR /go/src/callsforservice
# RUN apk add --update git build-base
# COPY . .
# RUN go get ./cmd/cfs-api; \
#      go build -o ./app ./cmd/cfs-api


# FROM alpine:latest
# WORKDIR /root
# COPY --from=builder /CallsForService/CFS-API/app .
# CMD ["./app"]

FROM golang:latest as builder

WORKDIR /go/src/
ARG CGO_ENABLED=0
COPY ./ ./
# COPY go.mod go.sum ./
RUN go mod download
RUN go build -o ./CFS-API/app ./CFS-API/cmd;

FROM alpine:latest
WORKDIR /root
ARG BUILD_USER=""
ARG BUILD_COMMIT_SHA=""
ARG BUILD_TIME=""

ENV BUILD_COMMIT_SHA=${BUILD_COMMIT_SHA}
ENV BUILD_USER=${BUILD_USER}
ENV BUILD_TIME=${BUILD_TIME}
COPY --from=builder ./go/src/CFS-API/app .
CMD ["./app"]

