    # stage 1: building application binary file
    FROM golang:1.21.3-alpine3.18 as build
     
    WORKDIR /source
    COPY . .
     
    ARG COMMIT
    RUN CGO_ENABLED=0 go mod tidy
    RUN go mod vendor
    RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.commit=${COMMIT}" -o bin/gf-user main.go
     
    # stage 2: copy only the application binary file and necessary files to the alpine container
    FROM alpine:3.12
    RUN apk --update add ca-certificates
    RUN mkdir -p /bin/application/ssl
     
    COPY --from=build /source/bin/gf-user /bin/application/gf-user
    COPY --from=build /source/database /bin/application/database
     
    EXPOSE 5001
     
    ENTRYPOINT ["./bin/application/gf-user"]