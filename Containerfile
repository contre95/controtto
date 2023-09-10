FROM golang:latest AS builder
WORKDIR /app
ADD . .
RUN go mod tidy
RUN go build -o /app/controtto */**.go

FROM scratch
ENV CONTROTTO_DB_PATH=/data/pnl.db
WORKDIR /app
COPY --from=builder /app/controtto .
ENTRYPOINT ["/app/controtto"]
