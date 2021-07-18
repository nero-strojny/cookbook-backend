
# Telling to use Docker's golang ready image
FROM golang
# Generate binary file from our /server
RUN go build
# Expose the port 8080
EXPOSE 8080:8080
# Run the app binary file 
CMD ["./server"]