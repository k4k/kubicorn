#!/usr/bin/env bash
set -e
cd ~

# ------------------------------------------------------------------------------------------------------------------------
# These values are injected into the script. We are explicitly not using a templating language to inject the values
# as to encourage the user to limit their use of templating logic in these files. By design all injected values should
# be able to be set at runtime, and the shell script real work. If you need conditional logic, write it in bash
# or make another shell script.
#
#
TOKEN="INJECTEDTOKEN"
# ------------------------------------------------------------------------------------------------------------------------

sudo curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
sudo touch /etc/apt/sources.list.d/kubernetes.list
sudo sh -c 'echo "deb http://apt.kubernetes.io/ kubernetes-xenial main" > /etc/apt/sources.list.d/kubernetes.list'

sudo apt-get update -y
sudo apt-get install -y \
    socat \
    ebtables \
    docker.io \
    apt-transport-https \
    kubelet \
    kubeadm

sudo systemctl enable docker
sudo systemctl start docker

sudo -E kubeadm reset
<<<<<<< HEAD
sudo -E kubeadm init

sudo kubectl apply \
  -f http://docs.projectcalico.org/v2.3/getting-started/kubernetes/installation/hosted/kubeadm/1.6/calico.yaml \
  --kubeconfig /etc/kubernetes/admin.conf
=======
sudo -E kubeadm init --apiserver-bind-port 443 --token ${TOKEN}
>>>>>>> 23305e5... Working Kubernetes :)
