# syntax=docker/dockerfile:1

FROM golang:1.16-alpine AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /vauth

# STAGE 2: build the container to run
FROM gcr.io/distroless/static AS final
 
USER nonroot:nonroot

WORKDIR /

# copy compiled app
COPY --from=build --chown=nonroot:nonroot /vauth /

# run binary; use vector form
CMD ["/vauth"]