######################################
# Step 1: Build the application in
# a Go development image
######################################

FROM golang:1.11.5-alpine3.9 AS builder

# Set up work directory
WORKDIR /app

# Copy local files in Docker image 
COPY . .

# Get dependencies
RUN go get -v -d ./...

# Compile application
RUN go build -a -installsuffix cgo -o app .

# Check required shared libraries
RUN ldd app

######################################
# Step 2: Build app in tiny image
######################################

# Start with empty image
FROM scratch

# Define working directory
WORKDIR /app

# Add CA certificates (required for X509 email)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy in /lib the required shared libraries
COPY --from=builder \
	/lib/ld-musl-x86_64.so.1 \
	/lib/libc.musl-x86_64.so.1
	/lib/

# Copy our application & driver in /app
COPY --from=builder /app/app /app/
COPY --from=builder /app/static /app/static
COPY --from=builder /app/templates /app/templates

# Prepare environment  variable: DEBUG
ARG DEBUG
ENV DEBUG ${DEBUG}

# Expose port 80
EXPOSE 80

ENTRYPOINT ["/app/app"]
