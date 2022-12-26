FROM golang:1.19-alpine AS build
WORKDIR /build
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/
COPY go.mod .
COPY go.sum .
RUN go build -o ranna cmd/ranna/main.go

FROM alpine:latest AS final
COPY --from=build /build/ranna /bin/ranna
COPY spec/ spec/
RUN chmod +x /bin/ranna
ENTRYPOINT ["/bin/ranna"]