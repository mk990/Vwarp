FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-checklinkname=0 -s -w" -o /vwarp ./cmd/vwarp

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /vwarp /usr/local/bin/vwarp

ENTRYPOINT ["vwarp"]
