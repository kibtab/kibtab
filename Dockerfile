# Kibtab images for the self-hosted stack.
#
# Two targets share this file, so docker compose needs one build context.
#
#   target kibtab  The instance, on a scratch base.
#   target caddy   The edge proxy, with the rate limit module.
#
# Build the instance: docker build --target kibtab -t kibtab .
# Build the proxy:    docker build --target caddy -t kibtab-caddy .
#
# Read docs/self-hosting/deployment.md before you run either build.

FROM golang:1.22-alpine AS kibtab-build
WORKDIR /src
# A build with no C compiler keeps the binary static. No dependency needs one.
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -trimpath -ldflags "-s -w" -o /out/kibtab ./cmd/kibtab

FROM scratch AS kibtab
# The database client needs the certificate authorities for a TLS connection.
COPY --from=kibtab-build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=kibtab-build /out/kibtab /kibtab
COPY LICENSE /LICENSE
COPY docs/THIRD_PARTY_NOTICES.md /THIRD_PARTY_NOTICES.md
# The instance runs as a user with no home directory and no shell.
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/kibtab"]

# Caddy holds the rate limit, so the image carries the module for it.
FROM caddy:builder AS caddy-build
RUN xcaddy build --with github.com/mholt/caddy-ratelimit

FROM caddy:2-alpine AS caddy
COPY --from=caddy-build /usr/bin/caddy /usr/bin/caddy
