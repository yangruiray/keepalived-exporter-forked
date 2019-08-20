# Download pkg env
FROM centos:7.6.1810
RUN rpm --import /etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-7 \
  && yum update -y \
  && yum install vim -y \
  && rm -rf /usr/local/src/*

FROM golang:1.12

ARG GOOS=linux
ARG GOARCH=amd64

# Maintainer
MAINTAINER yangrui@kpaas.io

# Initial Workdir
RUN mkdir -p /go/src/keepalived-exporter
RUN mkdir -p /etc/keepalived
WORKDIR /go/src/keepalived-exporter

# Copy file to container
COPY . /go/src/keepalived-exporter

# Volume
#VOLUME ["/etc/keepalived/keepalived.conf", "/etc/keepalived/keepalived.conf"]

# Set env
ENV PORT 9999

# Export
EXPOSE $PORT

# build 
RUN GOOS=linux GOARCH=amd64 go build -o keepalived-exporter main.go && cp keepalived-exporter /usr/local/bin/ && chmod +x /usr/local/bin/keepalived-exporter

# Run container command 
ENTRYPOINT ["keepalived-exporter"]

