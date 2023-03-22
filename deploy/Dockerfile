# Build the manager binary
FROM golang:alpine as builder
# Copy in the go src
WORKDIR /root
COPY . .
# Build
WORKDIR /root/
RUN apk --no-cache add tzdata
RUN apk --no-cache add ca-certificates
WORKDIR /root/cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /root/cmd/banklocator
# Copy the basedetector into a thin image

FROM scratch
WORKDIR /root/
ENV TZ Asia/Shanghai
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /root/cmd/banklocator .
COPY --from=builder /root/cmd/templates ./templates
COPY --from=builder /root/cmd/static ./static
CMD ["./banklocator"]
