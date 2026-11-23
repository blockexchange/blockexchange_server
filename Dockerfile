FROM node:26.13.1-alpine as node-app
WORKDIR /public
COPY /public/package-lock.json /public/package.json ./
RUN npm ci
COPY public/ .
RUN npm run lint && npm run build

FROM golang:1.27 as go-app
WORKDIR /data
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=node-app /public/dist /data/public/dist
RUN go vet && \
	go test ./... && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

FROM alpine:3.24.4
# git is needed to clone the mod dependencies for the wasm packs
RUN apk add --no-cache git ca-certificates
COPY --from=go-app /data/blockexchange /
EXPOSE 8080

CMD ["/blockexchange"]
