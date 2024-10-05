FROM golang:1.22-alpine as builder

WORKDIR /app

# Copy the application source code and build the binary
COPY . .
RUN go build -o load-balancer

### 
## Step 2: Runtime stage
FROM scratch

# Copy only the binary from the build stage to the final image
COPY --from=builder /app/ /

# Set the entry point for the container
ENTRYPOINT ["/load-balancer"]