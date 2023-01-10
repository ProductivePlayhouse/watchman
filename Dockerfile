FROM golang:1.19-bullseye as backend
WORKDIR /backend
RUN apt-get update && apt-get upgrade -y && apt-get install make gcc g++
COPY . .
RUN go mod download
RUN make build-server

FROM node:18-buster as frontend
WORKDIR /frontend
COPY webui/ .
RUN npm install --legacy-peer-deps
RUN npm run build

FROM debian:stable-slim
ENV AWS_REGION=us-west-2
WORKDIR /watchman
RUN apt-get update && apt-get upgrade -y && apt-get install -y ca-certificates
COPY --from=backend /backend/bin/server /bin/server
COPY --from=frontend /frontend/build/ /watchman/
ENV WEB_ROOT=/watchman/

EXPOSE 8080
EXPOSE 9090
ENTRYPOINT ["/bin/server"]
