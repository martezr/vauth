# syntax=docker/dockerfile:1

FROM golang:1.16-alpine AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o vauth

# STAGE 2: build the container to run
FROM gcr.io/distroless/static AS final
 
USER nonroot:nonroot
 
# copy compiled app
COPY --from=build --chown=nonroot:nonroot /app/vauth /vauth
 
# run binary; use vector form
ENTRYPOINT ["/vauth"]