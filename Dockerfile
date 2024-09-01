FROM golang:1.23-alpine AS build

RUN adduser --uid 1000 --disabled-password porkbun-ddns-user

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build

FROM scratch
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/porkbun-ddns /porkbun-ddns
USER porkbun-ddns-user
CMD ["/porkbun-ddns"]
