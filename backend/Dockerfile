FROM golang:alpine AS builder
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app/backend

COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

COPY . .
ARG VERSION=dev
ARG COMMIT=none
ENV CGO_ENABLED=0 GOOS=linux

RUN go build -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" \
  -o /out/app ./cmd/api/main.go

FROM alpine:3.20
RUN adduser -D -H -u 10001 appuser \
  && apk add --no-cache ca-certificates tzdata
WORKDIR /app/backend
COPY --from=builder /out/app /app/backend/app
COPY internal/configs/*.json /app/backend/internal/configs/

ENV PORT=3000
EXPOSE 3000
USER appuser
ENTRYPOINT ["/app/backend/app"]