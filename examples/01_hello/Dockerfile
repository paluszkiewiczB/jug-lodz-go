FROM golang:1.22 as builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o app main.go

FROM scratch
COPY --from=builder /build/app app
CMD ["./app"]