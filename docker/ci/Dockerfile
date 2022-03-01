    FROM golang:1.16-alpine3.12 as build
     
    WORKDIR /source
    COPY . .
     
    ARG COMMIT
    RUN CGO_ENABLED=0 go mod tidy
    RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.commit=${COMMIT}" -o bin/gf-user main.go
     
    FROM alpine:3.12
     
    COPY --from=build /source/bin/gf-user /bin/application/gf-user
    COPY --from=build /source/database /bin/application/database
     
    EXPOSE 5002
     
    ENTRYPOINT ["./bin/application/gf-user"]