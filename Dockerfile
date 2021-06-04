FROM golang:1.16 AS builder
WORKDIR /app
COPY . .
RUN make compile

FROM alpine:3.13
WORKDIR /app
COPY --from=builder /app/toggle .
EXPOSE 8080
EXPOSE 8081
CMD ["/app/toggle"]
