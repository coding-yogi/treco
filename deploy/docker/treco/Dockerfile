ARG BASE_IMAGE="gcr.io/distroless/base:nonroot-amd64"
FROM $BASE_IMAGE
COPY treco .
RUN chmod +x ./startup.sh
CMD ["./startup.sh"]
EXPOSE 8080