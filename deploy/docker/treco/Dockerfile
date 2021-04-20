ARG BASE_IMAGE="gcr.io/distroless/base:nonroot-amd64"
FROM $BASE_IMAGE

COPY ../../../treco/build .
CMD ["./treco serve"]