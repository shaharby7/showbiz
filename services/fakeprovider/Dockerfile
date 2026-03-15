FROM golang:1.21-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /fakeprovider ./cmd/fakeprovider

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /fakeprovider /fakeprovider
EXPOSE 8081
ENTRYPOINT ["/fakeprovider"]
