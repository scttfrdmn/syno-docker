# Dockerfile for SynoDeploy
# Copyright 2025 Scott Friedman

FROM scratch

# Add certificates for SSL/TLS
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY synodeploy /usr/local/bin/synodeploy

# Create non-root user
USER 1000:1000

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/synodeploy"]