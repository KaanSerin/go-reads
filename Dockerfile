# Official Go image
FROM golang:1.22.1

# Creating a directory inside the image
WORKDIR /app

# The current directory is /app
COPY go.mod go.sum ./

RUN go mod download

# Copy all files to /app folder
# Create .env using .env.example if it doesn't exist
COPY . ./
RUN cp -n ./.env.example ./.env

RUN rm -rf ./bin
RUN mkdir bin

# Building the application and placing the output at the root of image
RUN go build -o ./bin ./...

EXPOSE 8080

ENTRYPOINT [ "./bin/go_reads" ]
