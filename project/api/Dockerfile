FROM golang:1.16 AS build

WORKDIR /src
ADD go.mod go.sum /src/
RUN go mod download
ADD . /src
RUN go build -o /build/dwk-api

FROM gcr.io/distroless/base
COPY --from=build /build/dwk-api /dwk-api
CMD ["/dwk-api"]
