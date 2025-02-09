#!/bin/bash

URL=https://github.com/kubernetes-sigs/kubebuilder/releases/download/v4.5.0/kubebuilder_$(go env GOOS)_$(go env GOARCH)
wget -O bin/kubebuilder $URL
chmod a+x bin/kubebuilder