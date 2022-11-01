ARG GO_VERSION=1.18

FROM golang:${GO_VERSION}-alpine as builder
RUN apk --no-cache add ca-certificates git

WORKDIR /src

# Fetch dependencies first since they are less susceptible to change on every build
# these will be cached for speeding up the next build
COPY ./go.mod ./go.sum ./
RUN go mod download

# Now import code, for which we can expect most change from.
COPY . .
RUN CGO_ENABLED=0 go build -o /tests.api .

FROM scratch as final

COPY --from=builder /tests.api /tests.api
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["./tests.api"]
