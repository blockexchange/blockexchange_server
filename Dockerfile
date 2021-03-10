FROM golang:1.16.0-alpine as builder
COPY . /data
RUN cd /data && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

FROM scratch
COPY --from=builder /data/blockexchange /
EXPOSE 8080

CMD ["/blockexchange"]
