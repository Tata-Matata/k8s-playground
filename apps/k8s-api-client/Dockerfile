# STAGE 1 BUILD

# Linux container with Go installed
FROM golang:1.23 AS build

WORKDIR /src

# copies just go.mod first. This is mainly for Docker layer caching.
# downloads Go dependencies. Because go.mod is copied separately, Docker can often reuse this layer if dependencies did not change.
COPY go.mod ./
RUN go mod download

# copy rest of app source and build
COPY . .

#CGO_ENABLED=0 makes it a static binary, useful for tiny runtime images.
#GOOS=linux and GOARCH=amd64 = built for Linux on amd64.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/k8s-api-client ./...

# STAGE 2 Copy into runtime image

# distroless contains almost nothing except what is needed to run the binary 
# nonroot means it runs as a non-root user by default.
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=build /out/k8s-api-client /app/k8s-api-client

# metadata documenting that the app listens on port 8080. It does not publish the port by itself.
EXPOSE 8080

# the binary starts when container rolls out
ENTRYPOINT ["/app/k8s-api-client"]
