ARG BASE_IMAGE="gcr.io/distroless/base:nonroot-amd64"
FROM $BASE_IMAGE
COPY treco .
CMD source "/vault/secrets/dev-treco-db01-rw"; ./treco serve
EXPOSE 8080