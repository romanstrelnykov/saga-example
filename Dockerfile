# Build the Go Binary.
FROM golang as build_routes
ENV CGO_ENABLED=0

# Create a location in the container for the source code. Using the
# default GOPATH location.
RUN mkdir -p /service

# Copy the module files first and then download the dependencies. If this
# doesn't change, we won't need to do this again in future builds.
# COPY go.* /service/
# WORKDIR /service
# RUN go mod download

# Copy the source code into the container.
WORKDIR /service
COPY . .

# Build the service binary. We are doing this last since this will be different
# every time we run through this process.
WORKDIR /service/cmd/saga
RUN go build 


# Run the Go Binary in Alpine.
FROM alpine
ARG BUILD_DATE
COPY --from=build_routes /service/cmd/saga/saga /service/saga
WORKDIR /service
CMD ["./saga"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="saga" \
      org.opencontainers.image.authors="Roman Strelnykov <roman.strelnykov@gmail.com>" 

