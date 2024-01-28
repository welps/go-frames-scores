FROM golang:1.21-alpine AS build

ARG TARGETARCH

RUN apk update && apk add make git gcc musl-dev
# Add a non-root user that we will copy into the container that runs the bin
RUN adduser -u 10001 app -D

ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-static-$TARGETARCH /tini
RUN chmod +x /tini

WORKDIR /app

COPY . .

RUN make build

FROM scratch AS bin

COPY --from=build /tini /tini
ENTRYPOINT ["/tini", "--"]

WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/bin/go-frames-scores /app/
COPY --from=build /etc/passwd /etc/passwd

USER app
CMD ["/app/go-frames-scores"]

