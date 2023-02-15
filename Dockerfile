FROM golang:1.20.1 as stage1
WORKDIR /data
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go vet && \
	go test ./... && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

FROM alpine:3.17.1
COPY --from=stage1 /data/blockexchange /
EXPOSE 8080

CMD ["/blockexchange"]
