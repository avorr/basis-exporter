FROM golang:1.21 as builder

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum
COPY client client
RUN go mod download

COPY main.go main.go

RUN CGO_ENABLED=0 go build -a -o main main.go

FROM alpine
WORKDIR /
COPY --from=builder /build/main .
COPY --from=builder /build/client ./client

ENTRYPOINT ["/main"]
