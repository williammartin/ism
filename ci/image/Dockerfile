FROM golang:alpine3.8

COPY install-deps.sh install-deps.sh

RUN ./install-deps.sh

ENV PATH="${PATH}:/usr/local/kubebuilder/bin"
