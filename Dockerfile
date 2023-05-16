###
### Build Stage
###

# The base go-image
FROM golang:1.20.4-alpine as build-env

# Create a directory for the app
RUN mkdir /app

# Copy all files from the current directory to the app directory
COPY . /app

# Set working directory
WORKDIR /app

# Run command as described:
# go build will build a 64bit Linux executable binary file named server in the current directory
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cps-backend .

###
### Run stage.
###

FROM alpine:latest

# Copy only required data into this image
COPY --from=build-env /app/cps-backend .

# Copy all the static content necessary for this application to run.
COPY --from=build-env /app/static ./static

EXPOSE 8000

# Run the server executable
CMD [ "/cps-backend" ]

# BUILD
# (Run following from project root directory!)
# docker build -f Dockerfile -t LuchaComics/cps-backend:latest --platform linux/amd64 .

# EXECUTE
# docker tag LuchaComics/cps-backend:latest LuchaComics/cps-backend:latest

# UPLOAD
# docker push LuchaComics/cps-backend:latest
