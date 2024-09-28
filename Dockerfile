FROM golang:1.22.7-alpine3.20 AS builder

RUN apk add --no-cache make

COPY . .

RUN make release

FROM scratch

COPY --from=builder release/iam .

CMD ["iam"]
