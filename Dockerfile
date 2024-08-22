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

FROM ${CONTAINER_REPO}/ironbank/redhat/ubi/ubi8:latest
ENV AWS_REGION=us-west-2
WORKDIR /watchman
USER root
RUN dnf update -y && dnf upgrade -y && dnf -y install ca-certificates
USER 1001
COPY --from=backend /backend/bin/server /bin/server
COPY --from=frontend /frontend/build/ /watchman/
ENV WEB_ROOT=/watchman/

EXPOSE 8080
EXPOSE 9090
ENTRYPOINT ["/bin/server"]
