FROM gitlab.myteksi.net:4567/dakota/backend/ci/alpine:3.12.3
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
EXPOSE 8080
ENV GOPATH /go
WORKDIR ${GOPATH}/bin/
COPY ../../../treco .
CMD ["./treco"]