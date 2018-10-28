FROM golang:1.11.1-alpine3.8 as builder

# Install SSL ca certificates
RUN apk add --update --no-cache ca-certificates git

# Create appuser
RUN adduser -D -g '' appuser

RUN mkdir /krane
WORKDIR /krane
COPY go.mod /krane
COPY go.sum /krane

# Get dependancies
RUN go mod download

COPY . /krane

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/krane-operator ./cmd/krane-operator

# Build a small image
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /go/bin/krane-operator /go/bin/krane-operator

USER appuser

ENTRYPOINT ["/go/bin/krane-operator"]
