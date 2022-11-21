FROM golang:1.19 AS builder
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY constant/ /go/src/app/constant/
COPY openweather/ /go/src/app/openweather/
COPY proto/ /go/src/app/proto/
COPY utilenv/ /go/src/app/utilenv/
COPY utilgcp/ /go/src/app/utilgcp/
COPY v2/ /go/src/app/v2/
WORKDIR /go/src/app/v2/backgrpc
RUN CGO_ENABLED=0 go build -o app

FROM gcr.io/distroless/static-debian11 AS runtime
COPY --from=builder /go/src/app/v2/backgrpc/app /usr/local/bin/app
ENTRYPOINT ["/usr/local/bin/app"]
