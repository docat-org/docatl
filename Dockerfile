FROM golang:1.21-alpine AS build

RUN apk update \
    && apk add -U --no-cache ca-certificates \
    && update-ca-certificates

WORKDIR /src
COPY . /src/
RUN CGO_ENABLED=0 go build -o /bin/docatl

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/docatl /bin/docatl
WORKDIR /docs
ENTRYPOINT [ "/bin/docatl" ]
