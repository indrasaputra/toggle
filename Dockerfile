FROM golang:1.16 AS builder
WORKDIR /app
COPY . .
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.6 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe
RUN WAIT_FOR_VERSION=v2.1.2 && \
    wget -qO/bin/wait-for https://github.com/eficode/wait-for/releases/download/${WAIT_FOR_VERSION}/wait-for && \
    chmod +x /bin/wait-for
RUN make compile

FROM alpine:3.13
WORKDIR /app
COPY --from=builder /bin/grpc_health_probe ./grpc_health_probe
COPY --from=builder /bin/wait-for ./wait-for
COPY --from=builder /app/toggle .
COPY --from=builder /app/bin/start.sh ./start.sh
RUN chmod 755 /app/start.sh /app/wait-for
EXPOSE 8080
EXPOSE 8081
CMD ["./start.sh"]
