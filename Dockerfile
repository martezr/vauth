# syntax=docker/dockerfile:1

FROM golang:1.16-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

COPY frontend ./

RUN go build -o /vauth


FROM scratch

WORKDIR /

COPY --from=build /vauth /vauth

USER nonroot:nonroot

CMD [ "/vauth" ]