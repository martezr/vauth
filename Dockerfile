# syntax=docker/dockerfile:1

FROM golang:1.16-alpine AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN ls -al

RUN go build -o vauth


FROM scratch

WORKDIR /

COPY --from=build /app/vauth /vauth

USER nonroot:nonroot

CMD [ "/vauth" ]