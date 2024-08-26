ARG CONTAINER_REPO=registry1.dso.mil

FROM golang:1.20-bookworm as backend
WORKDIR /backend
RUN apt-get update && apt-get upgrade -y && apt-get install make gcc g++
COPY . .
RUN go mod download
RUN make build-server

FROM node:21-bookworm as frontend
WORKDIR /frontend
COPY webui/ .
RUN npm install --legacy-peer-deps
RUN npm run build

FROM registry.access.redhat.com/ubi9/ubi-minimal:9.4
ENV AWS_REGION=us-west-2
WORKDIR /watchman
USER root
RUN microdnf update -y && \
    microdnf install -y ca-certificates
USER 1001
COPY --from=backend /backend/bin/server /bin/server
COPY --from=frontend /frontend/build/ /watchman/
ENV WEB_ROOT=/watchman/

EXPOSE 8080
EXPOSE 9090
ENTRYPOINT ["/bin/server"]
