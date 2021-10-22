FROM envoyproxy/envoy-alpine:v1.20.0
COPY infrastructure/envoy.yaml /etc/envoy/envoy.yaml
COPY bin/image.bin /etc/envoy/toggle.pb
RUN chmod go+r /etc/envoy/envoy.yaml