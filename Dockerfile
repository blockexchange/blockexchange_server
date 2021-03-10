FROM golang:1.16.0 as builder
COPY . /data
RUN cd /data && go build .

FROM alpine:3.13.2
COPY --from=builder /data/blockexchange /blockexchange
EXPOSE 8080

CMD ["/blockexchange"]
