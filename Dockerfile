FROM node:23.6.1-alpine as node-app
WORKDIR /public
COPY /public/package-lock.json /public/package.json ./
RUN npm ci
COPY public/ .
RUN npm run jshint && \
	npm run bundle

FROM golang:1.23.5 as go-app
WORKDIR /data
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=node-app /public /data/public
RUN go vet && \
	go test ./... && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

FROM alpine:3.21.3
COPY --from=go-app /data/blockexchange /
EXPOSE 8080

CMD ["/blockexchange"]
