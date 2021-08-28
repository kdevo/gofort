FROM golang:1.17-alpine AS build

RUN apk add --no-cache make
RUN apk add --no-cache fortune

WORKDIR /go/src/gofort/

COPY . .
RUN BUILD_DIR=/bin/gofort make build

FROM scratch
COPY --from=build /bin/gofort/gofort* ./gofort
ENTRYPOINT ["./gofort"]
