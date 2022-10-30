FROM golang:1.18 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/geoip-lookup

FROM scratch
COPY --from=build /go/bin/geoip-lookup /usr/bin/geoip-lookup

ENTRYPOINT ["/usr/bin/geoip-lookup"]
