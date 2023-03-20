FROM ubuntu:latest

# update container certificates
RUN apt-get -y update && \
    apt-get install -y ca-certificates && \
    update-ca-certificates

# set the Current Working Directory inside the container
WORKDIR ~/app

# copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY ./build/ .

# this container exposes port 8100 to the outside world
EXPOSE 8100

# run the server executable
CMD ["./bobc"]
