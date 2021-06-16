FROM golang:1.16.3 as stage1
COPY . /data
RUN cd /data && \
	go vet && \
	go test ./... && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

FROM alpine:3.14.0
COPY --from=stage1 /data/auth_proxy /
EXPOSE 8080

CMD ["/auth_proxy"]
