FROM golang:1.23.3-alpine3.20  AS builder
RUN apk update && \
    apk add git


# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY . .

# Build
RUN rm -f go.mod go.sum 
RUN go mod init main
RUN GOPROXY=direct GOSUMDB=off go mod tidy 
RUN CGO_ENABLED=0 GOOS=linux go build -o /eGKsc

FROM scratch as copytohost
COPY --from=builder /eGKsc eGKsc
