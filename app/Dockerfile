FROM golang:1.25 AS builder
ARG VARIANT="server"

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd/${VARIANT}.go

FROM gcr.io/distroless/static-debian12
COPY --from=builder /go/bin/app /
CMD ["/app"]
