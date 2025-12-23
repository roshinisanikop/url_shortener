FROM golang:1.21-alpine AS build

WORKDIR /src

# Only copy go.mod (go.sum is optional, especially if you have no external deps)
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app ./...

FROM gcr.io/distroless/static
COPY --from=build /app /app
EXPOSE 8080
ENTRYPOINT ["/app"]