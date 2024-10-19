# syntax=docker/dockerfile:1
FROM golang:1.22.7 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /sender
FROM scratch AS final
COPY --from=build /sender /sender
ENTRYPOINT ["/sender"]