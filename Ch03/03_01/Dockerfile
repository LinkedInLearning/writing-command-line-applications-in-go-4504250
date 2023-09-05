FROM golang:1.21-bookworm AS build

WORKDIR /app
COPY * .
ENV CGO_ENABLED=0
RUN go build -o health .

FROM alpine:3.18
COPY --from=build /app/health /usr/local/bin

ENTRYPOINT ["health"]
