FROM golang:alpine AS builder
RUN apk add -U --no-cache ca-certificates
COPY . /build/ecrcredrotation
WORKDIR /build/ecrcredrotation
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o ecrcredrotation ./cmd/ecrcreds/main.go
RUN chmod 0755 ecrcredrotation
RUN ls -ltr

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/ecrcredrotation/ecrcredrotation /go/bin/ecrcredrotation

# Run the hello binary.
ENTRYPOINT ["/go/bin/ecrcredrotation"]