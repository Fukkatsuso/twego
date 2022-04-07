FROM golang:1.17-bullseye AS base

RUN apt-get update && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone

WORKDIR /go/src/github.com/Fukkatsuso/twego

COPY go.* ./
RUN go mod download

FROM base AS builder

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -x -o /go/bin/twego

FROM debian:bullseye-slim AS runner

RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/* && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone

COPY --from=builder /go/bin/twego /go/bin/twego

ENTRYPOINT [ "/go/bin/twego" ]
CMD [ "--help" ]
