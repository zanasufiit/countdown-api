FROM golang:1.13.5-buster AS builder

WORKDIR /build

COPY . .

RUN go build -a -o countdown-api .

# ---

FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY --from=builder /build/countdown-api /app

ENTRYPOINT ["./countdown-api"]