# Build
FROM golang:1.22 AS build
LABEL authors="gnousEU"

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY internal/ ./internal
COPY static/ ./static
COPY templates/ ./templates

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" # Enable static binary, target Linux, remove debug information and strip binary

# Copy to our image
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=build /build/plakken ./

ENV PLAKKEN_LISTEN ":3000"

EXPOSE 3000/tcp

ENTRYPOINT ["/app/plakken"]
