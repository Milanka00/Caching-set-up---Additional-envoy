# FROM envoyproxy/envoy-alpine:v1.20-latest
# COPY envoy.yaml /etc/envoy/envoy.yaml
# CMD ["envoy", "-c", "/etc/envoy/envoy.yaml", "--log-level", "info"]

FROM envoyproxy/envoy-alpine:v1.20-latest

COPY envoy.yaml /etc/envoy/envoy.yaml

ENTRYPOINT ["/usr/local/bin/envoy"]

CMD ["--service-node", "proxy", "--service-cluster", "proxy", "--config-path", "/etc/envoy/envoy.yaml", "--concurrency", "1", "--mode", "serve", "--log-level", "info"]

