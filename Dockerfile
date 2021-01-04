FROM golang:1.15-alpine3.12
RUN go get -u github.com/lucthienbinh/golang_scem
ARG GO111MODULE=on
RUN mkdir /app
## We copy everything in the root directory into our /app directory
ADD . /app
## Add this go mod download command to pull in any dependencies
RUN cd /app
RUN go mod download
## We specify that we now wish to execute any further commands inside our /app/cmd/scem directory
WORKDIR /app/cmd/scem
## we run go build to compile the binary
## executable of our Go program
RUN go build -o main .
## Our start command which kicks off
## our newly created binary executable
CMD ["/app/main"]