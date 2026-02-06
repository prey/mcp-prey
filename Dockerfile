# Build stage
FROM golang:1.25 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/mcp-prey ./cmd/mcp-prey

# Runtime stage
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=build /out/mcp-prey /app/mcp-prey

ENV PREY_API_BASE=https://api.preyproject.com/v1

ENTRYPOINT ["/app/mcp-prey"]
