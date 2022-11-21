FROM golang:1.19 AS builder
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY constant/ /go/src/app/constant/
COPY openweather/ /go/src/app/openweather/
COPY proto/ /go/src/app/proto/
COPY util/ /go/src/app/util/
COPY v1/ /go/src/app/v1/
WORKDIR /go/src/app/v1/frontweb
RUN CGO_ENABLED=0 go build -o app

FROM gcr.io/distroless/static-debian11 AS runtime
COPY --from=builder /go/src/app/v1/frontweb/app /usr/local/bin/app
ENTRYPOINT ["/usr/local/bin/app"]
