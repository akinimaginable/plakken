# Build
FROM golang:1.23 AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY static/ ./static
COPY templates/ ./templates
COPY main.go ./
COPY internal/ ./internal

# Enable static binary, target Linux, remove debug information and strip binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s"

# Copy to our image
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=build /build/plakken ./

ENV PLAKKEN_LISTEN ":3000"

EXPOSE 3000/tcp

ENTRYPOINT ["/app/plakken"]