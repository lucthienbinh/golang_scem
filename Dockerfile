FROM golang:1.15-alpine3.12 AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    RUNENV=docker

# Copy bpmn workflow 
WORKDIR /storage/private/zeebe
COPY ./storage/private/zeebe .

# Copy sample picture and create folder
WORKDIR /storage/private/images
WORKDIR /storage/public/upload/images
COPY ./storage/public/upload/images .

# Move to working directory /build
WORKDIR /build

# Copy the code into the container
COPY . .

## We specify that we now wish to execute any further commands inside our /build/cmd/scem directory
WORKDIR /build/cmd/scem

## we run go build to compile the binary executable of our Go program
RUN go build -o golang_scem_binary .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /

# Copy binary from build to main folder
RUN cp /build/cmd/scem/golang_scem_binary .

EXPOSE 5000
EXPOSE 5001

CMD ["/golang_scem_binary"]

# # Build a small image
# FROM scratch

# COPY --from=builder /dist/main /

# # Command to run
# ENTRYPOINT ["/main"]