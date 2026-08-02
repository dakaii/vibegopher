FROM golang:1.24-bookworm AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /
COPY --from=build /out/server /server

# Cloud Run sets PORT; the app already reads it (default 8081 for local).
ENV PORT=8080
EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/server"]
