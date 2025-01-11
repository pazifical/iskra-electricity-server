FROM golang:1.23 AS builder

WORKDIR /workdir

COPY go.mod .

COPY iskra iskra
COPY internal internal
COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o iskra-electricity-server .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /workdir/iskra-electricity-server /
CMD ["/iskra-electricity-server"]