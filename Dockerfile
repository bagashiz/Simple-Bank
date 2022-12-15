# build stage
FROM golang:1.18-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY db/migration ./db/migration

EXPOSE 8080
EXPOSE 9090
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
