FROM golang:1.22.7-alpine3.20 AS builder

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /i_a_m

FROM gcr.io/distroless/static-debian12

COPY --from=builder /i_a_m /

CMD ["/i_a_m"]
