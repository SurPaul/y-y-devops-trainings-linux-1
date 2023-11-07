FROM golang:1.21-alpine AS build-stage

WORKDIR /data

COPY go.mod go.sum ./
COPY *.go ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app

FROM gcr.io/distroless/static-debian12 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app /app

EXPOSE 8080

CMD [ "/app" ]