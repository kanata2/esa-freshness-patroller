FROM golang:1.18-stretch AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build ./cmd/esa-freshness-patroller

FROM gcr.io/distroless/base-debian10
COPY --from=build /app/esa-freshness-patroller /esa-freshness-patroller
USER nonroot:nonroot
ENTRYPOINT ["/esa-freshness-patroller"]
