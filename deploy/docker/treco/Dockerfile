ARG BASE_IMAGE="gcr.io/distroless/base:nonroot-amd64"
FROM $BASE_IMAGE
COPY treco .
CMD ./treco serve -c "/vault/secrets/treco_db_creds.env"
EXPOSE 8080