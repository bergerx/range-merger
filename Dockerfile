# syntax=docker/dockerfile:1

FROM golang:1.18-buster AS build
WORKDIR /app
COPY go.mod ./
COPY *.go ./
RUN go build

FROM gcr.io/distroless/base-debian10
WORKDIR /
EXPOSE 8090
USER nonroot:nonroot
COPY --from=build /app/range-merger /range-merger
ENTRYPOINT ["/range-merger"]
