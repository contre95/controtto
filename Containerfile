FROM golang:latest AS builder
WORKDIR /app
ADD . .
# COPY src .
# COPY cmd .
# COPY public .
# COPY go.mod .
# COPY go.sum .
RUN go mod tidy
RUN go build -o /app/controtto */**.go
ENV CONTROTTO_DB_PATH=/data/pnl.db
ENTRYPOINT ["/app/controtto"]
