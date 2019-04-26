############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/main/sonar-qualitygate-validator/
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -d -v
# Build the binary.
RUN go build -o /bin/sonar-qualitygate-validator
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /bin/sonar-qualitygate-validator

############################
# STEP 2 build a small image
############################
FROM alpine
# Copy our static executable.
COPY --from=builder /bin/sonar-qualitygate-validator /bin/sonar-qualitygate-validator
RUN apk add ca-certificates

# Run the binary
CMD ["/bin/sonar-qualitygate-validator"]
