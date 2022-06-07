FROM node:18.3.0-alpine as stage1
COPY public /public
RUN cd /public && \
	npm ci && \
	npm run jshint && \
	npm run bundle

FROM golang:1.18.3 as stage2
COPY . /data
COPY --from=stage1 /public /data/public
RUN cd /data && \
	go vet && \
	go test ./... && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

FROM alpine:3.16.0
COPY --from=stage2 /data/blockexchange /
EXPOSE 8080

CMD ["/blockexchange"]
