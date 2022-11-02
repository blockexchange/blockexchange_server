FROM golang:1.19.3 as stage1
COPY . /data
RUN cd /data && \
	go vet && \
	go test ./... && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

FROM alpine:3.16.2
COPY --from=stage1 /data/blockexchange /
EXPOSE 8080

CMD ["/blockexchange"]
