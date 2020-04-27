FROM golang:1.14 AS build

WORKDIR /go/src/github.com/dstream.cloud/configmap-map-operator 
ADD . /go/src/github.com/dstream.cloud/configmap-map-operator 

RUN go build -o /go/bin/configmap-map-operator ./cmd/manager

FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/configmap-map-operator  /
CMD ["/configmap-map-operator"]
