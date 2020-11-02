############################
# STEP 1 build executable binary
############################
FROM golang:latest as builder

ARG service

COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -o ./bin/svc ./cmd/$service/...

############################
# STEP 2 build a small image
############################

FROM scratch

WORKDIR /app

# Copy our static executable
COPY --from=builder /app/bin/svc ./svc

# Run the svc binary.
CMD ["./svc"]
