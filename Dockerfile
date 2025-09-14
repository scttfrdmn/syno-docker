# Dockerfile for syno-docker
# Copyright 2025 Scott Friedman

FROM scratch

# Add certificates for SSL/TLS
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY syno-docker /usr/local/bin/syno-docker

# Create non-root user
USER 1000:1000

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/syno-docker"]