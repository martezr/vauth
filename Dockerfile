# syntax=docker/dockerfile:1

FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o vauth

RUN touch /app/vauth.db

# STAGE 2: build the container to run
FROM gcr.io/distroless/static AS final
 
USER nonroot:nonroot

WORKDIR /app

# copy compiled app
COPY --from=build --chown=nonroot:nonroot /app /app

# run binary; use vector form
CMD ["/app/vauth","server"]