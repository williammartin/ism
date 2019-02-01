#!/usr/bin/env sh

arch="amd64"

dep_version="0.5.0"
kubebuilder_version="1.0.8"

# install packages
apk add build-base git

# install helpful golang tools
wget -O /usr/local/bin/dep "https://github.com/golang/dep/releases/download/v${dep_version}/dep-linux-${arch}"
chmod +x /usr/local/bin/dep
go get -u github.com/onsi/ginkgo/ginkgo

# install kubebuilder
wget "https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${kubebuilder_version}/kubebuilder_${kubebuilder_version}_linux_${arch}.tar.gz" -O - | tar -C /usr/local -xzf -
ln -s "/usr/local/kubebuilder_${kubebuilder_version}_linux_${arch}" "/usr/local/kubebuilder"
