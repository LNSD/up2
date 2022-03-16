FROM golang:1.18-alpine AS builder
RUN apk add --no-cache make

WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod vendor

COPY . .
RUN make build


FROM alpine:3.15
COPY --from=builder /src/out/bin/up2 /usr/bin/up2

EXPOSE 8080
ENTRYPOINT ["up2"]