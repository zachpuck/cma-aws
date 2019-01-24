# Initial preparations
FROM golang:1.10.3 as builder
RUN apk update && apk add git make curl

# Let's add the code in
ADD . /go/src/gitlab.com/mvenezia/cma-aws

# Build package in container
WORKDIR /go/src/gitlab.com/mvenezia/cma-aws
RUN make -f build/Makefile container-build

# Now that we have build a static binary, we can stuff it in a bare-bones alpine image
# But we will need bash and openssh right now
# ca-certificates is needed for aws library
FROM alpine:3.8
RUN apk update && apk add bash openssh ca-certificates
COPY --from=builder /go/src/gitlab.com/mvenezia/cma-aws/cma-aws /cma-aws

ENTRYPOINT ["/cma-aws"]
