FROM golang:1.24 AS build

WORKDIR /go/src/docker-mtrpz
COPY . .

RUN go test ./...
ENV CGO_ENABLED=0
RUN go install ./cmd/...

# ==== Final image ====
FROM alpine:latest
WORKDIR /opt/docker-mtrpz
COPY entry.sh /opt/docker-mtrpz/
RUN chmod +x /opt/docker-mtrpz/entry.sh
COPY --from=build /go/bin/* /opt/docker-mtrpz
RUN ls /opt/docker-mtrpz
ENTRYPOINT ["/opt/docker-mtrpz/entry.sh"]
CMD ["server"]