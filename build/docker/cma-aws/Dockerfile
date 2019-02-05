# Initial preparations
FROM golang:1.10.3 as builder

# Let's add the code in
ADD . /go/src/github.com/samsung-cnct/cma-aws

# Build package in container
WORKDIR /go/src/github.com/samsung-cnct/cma-aws
RUN make -f build/Makefile container-build

# Now that we have build a static binary, we can stuff it in a bare-bones alpine image
# But we will need bash and openssh right now
# ca-certificates is needed for aws library
FROM alpine:3.8
RUN apk update && apk add bash openssh ca-certificates
COPY --from=builder /go/src/github.com/samsung-cnct/cma-aws/cma-aws /cma-aws

ENTRYPOINT ["/cma-aws"]
