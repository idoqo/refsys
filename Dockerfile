FROM golang:alpine AS builder
MAINTAINER Michael Okoko <me@mchl.xyz>

COPY go.mod go.sum /go/src/gitlab.com/idoko/refsys/
WORKDIR /go/src/gitlab.com/idoko/refsys
RUN go mod download
COPY . /go/src/gitlab.com/idoko/refsys
RUN CGO_ENABLED=0 GOOS=linux go build -a --installsuffix cgo -o build/refsys ./cmd/refsys

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/gitlab.com/idoko/refsys/build/refsys /usr/bin/refsys

ENTRYPOINT ["/usr/bin/refsys"]