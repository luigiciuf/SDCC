#!/bin/bash

# INSTALLARE GO
echo "Installing Go..."
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version

# INSTALLARE DOCKER
echo "Installing Docker..."
sudo yum install -y docker
sudo service docker start
sudo usermod -aG docker ec2-user

# INSTALLARE DOCKER COMPOSE
echo "Installing Docker Compose..."
sudo curl -L "https://github.com/docker/compose/releases/download/v2.16.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

echo "Installation complete. Please log out and log back in to apply Docker permissions."
