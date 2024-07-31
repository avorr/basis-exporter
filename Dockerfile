FROM docker-pid.art.gos-tech.xyz/golang:1.22 as builder

ARG version=dev

WORKDIR /build

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 go build -ldflags="-X 'main.appVersion=${version}'" -a -o basis-exporter main.go

FROM docker-pid.art.gos-tech.xyz/alpine

COPY --from=builder /build/basis-exporter .

ENTRYPOINT ["./basis-exporter"]
