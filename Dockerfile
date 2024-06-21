FROM docker.io/golang:1.22.4-bullseye as builder

WORKDIR /build

ENV CGO_ENABLED=0
ENV GOOS=linux 
ENV GOARCH=amd64

ADD . .

RUN cd cli && go build -o etz

FROM docker.io/alpine:3.20

WORKDIR /

COPY --from=builder /build/cli/etz .

ENTRYPOINT [ "./etz" ]
