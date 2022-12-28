FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# COPY *.go ./
COPY . ./

RUN ls

RUN go build -o app .

EXPOSE 2565

ENV DATABASE_URL=postgres://vvisrafd:Ul_i2EARyUxFjSSCvyo7TCaZP-EsFJxt@tiny.db.elephantsql.com/vvisrafd
ENV PORT=:2565

# CMD ["/app"]
