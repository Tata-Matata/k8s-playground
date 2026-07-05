# STAGE 1 BUILD
FROM golang:1.23 AS build

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/k8s-api-client ./...

# STAGE 2 Copy into runtime image
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=build /out/k8s-api-client /app/k8s-api-client

EXPOSE 8080

ENTRYPOINT ["/app/k8s-api-client"]
