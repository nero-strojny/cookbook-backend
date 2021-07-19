
# Telling to use Docker's golang ready image
FROM golang
# Create app folder 
RUN mkdir /server
# Copy our file in the host container to our container
ADD . /server
# Set /server to the go folder as workdir
WORKDIR /server
# Generate binary file from our /server
RUN go build
# Expose the port 8080
EXPOSE 8080:8080
# Run the app binary file 
CMD ["./server"]