# Build the manager binary
FROM golang:1.26 as builder

WORKDIR /workspace

# Copy the source.
COPY . .

# Build.
RUN make build CGO_ENABLED=0

# Use distroless as minimal base image to package the manager binary.
# Refer to https://github.com/GoogleContainerTools/distroless for more details.
FROM gcr.io/distroless/static:nonroot
LABEL source_repository="https://github.com/sapcc/digicert-issuer"
LABEL org.opencontainers.image.source="https://github.com/sapcc/digicert-issuer"
WORKDIR /
COPY --from=builder /workspace/bin/digicert-issuer .
USER nonroot:nonroot
RUN ["/digicert-issuer", "--version"]
ENTRYPOINT ["/digicert-issuer"]
