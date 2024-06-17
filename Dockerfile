FROM docker.io/golang:1.22.4-bullseye as builder

WORKDIR /build

ADD . .

RUN cd cli && go build -o etz

FROM docker.io/alpine:3.20

WORKDIR /cli

COPY --from=builder /build/cli/etz .

ENTRYPOINT [ "/cli/etz" ]
