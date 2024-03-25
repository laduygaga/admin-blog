# Stage 1: Build stage
FROM golang:1.21 as build

WORKDIR /golang
COPY . .

RUN CGO_ENABLED=0 go build -o ota-blog-admin cmd/web/main.go

# Stage 2: Final image
FROM alpine:latest

WORKDIR /golang

# Copy only necessary files from the build stage
COPY --from=build /golang/ota-blog-admin .
COPY --from=build /golang/templates templates
COPY --from=build /golang/static static
# COPY --from=build /golang/.env .

EXPOSE 9090
CMD ["./ota-blog-admin"]

