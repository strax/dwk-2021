FROM golang:1.16 AS build

WORKDIR /src
ADD go.mod go.sum /src/
RUN go mod download
ADD . /src
RUN go build -o /build/dwk-broadcaster

FROM gcr.io/distroless/base
COPY --from=build /build/dwk-broadcaster /dwk-broadcaster
CMD ["/dwk-broadcaster"]
