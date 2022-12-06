FROM golang:1.19.3 as stage1
WORKDIR /data
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go vet && \
	go test ./... && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

FROM scratch
COPY --from=stage1 /data/blockexchange /
EXPOSE 8080

CMD ["/blockexchange"]
