# ec2

## Launching AMI
Script on launch
```bash
#!/bin/bash
# Update the package manager
sudo yum update -y

# Install Git, Vim, and Go
sudo yum install git vim -y
sudo amazon-linux-extras install golang1.11 -y

# Set Go environment variables
echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
echo "export GOPATH=\$HOME/go" >> ~/.bashrc
source ~/.bashrc

# Install the interactsh-client
go install -v github.com/projectdiscovery/interactsh/cmd/interactsh-client@latest

```