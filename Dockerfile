FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

RUN apk add --no-cache git

WORKDIR /src

COPY go.mod go.sum ./
COPY patches/ ./patches/
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOARM=${TARGETVARIANT#v} \
    go build -ldflags="-checklinkname=0 -s -w" -o /vwarp ./cmd/vwarp

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /vwarp /usr/local/bin/vwarp

ENTRYPOINT ["vwarp"]
