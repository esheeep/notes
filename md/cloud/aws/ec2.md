# ec2

## Install go & interactsh
```bash
cd /usr/local/
sudo wget https://dl.google.com/go/go1.13.4.linux-amd64.tar.gz
sudo tar -xzf go1.13.4.linux-amd64.tar.gz
```

Add go to path
```bash
cd /etc/profile.d
sudo nano go.sh
#insert following lines: 
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
```
Install interactsh
```bash
go get -v github.com/projectdiscovery/interactsh/cmd/interactsh-client@latest
```

Set api key
```bash
./interactsh-client -auth
```