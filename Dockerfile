# syntax=docker/dockerfile:1

FROM golang:1.22.0 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /kreasi-nusantara .

FROM gcr.io/distroless/base-debian11 AS build-release

WORKDIR /

COPY --from=build /kreasi-nusantara .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT [ "/kreasi-nusantara" ]