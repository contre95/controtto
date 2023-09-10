FROM golang:alpine AS builder

# Required for alpine + sqlite3 driver
ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev

# Copy files
WORKDIR /app
ADD . .
# Buiild
RUN go mod tidy
RUN go build -ldflags='-s -w -extldflags "-static"' -o /app/controtto */**.go
# RUN CGO_ENABLED=1 GOOS=linux go build -o /app/controtto -installsuffix 'static' -a -ldflags '-s -w' cmd/main.go

FROM scratch
ENV CONTROTTO_DB_PATH=/data/pnl.db
LABEL maintainer="contre95"
# USER nonroot:nonroot
# COPY --from=builder --chown=nonroot:nonroot /app/controtto /app/controtto
WORKDIR /app
COPY ./views /app/views
COPY ./public /app/public
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/controtto /app/controtto
ENTRYPOINT ["/app/controtto"]
