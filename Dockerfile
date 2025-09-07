# syntax=docker/dockerfile:1
FROM golang:1.24 as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /gophkeeper ./cmd/gophkeeper

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /gophkeeper /gophkeeper
EXPOSE 8087
USER nonroot:nonroot
ENTRYPOINT ["/gophkeeper"]

