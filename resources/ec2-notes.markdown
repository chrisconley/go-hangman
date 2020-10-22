# Rsync up to ec2 instance
ip=<instance-ip>
ssh -t ubuntu@$ip "screen -dR setup"
ssh-add ~/.ssh/<key-file>.pem
rsync -avt --exclude=".git" --exclude="__pycache__" . ubuntu@$ip:/home/ubuntu/go-hangman
ssh -t ubuntu@$ip "screen -dR setup"

# On EC2, setting up go
Using this as a guide: https://www.linode.com/docs/development/go/install-go-on-ubuntu/

```
wget "https://golang.org/dl/go1.15.3.linux-amd64.tar.gz"
tar -xvf go1.15.3.linux-amd64.tar.gz go/
sudo chown -R root:root ./go
sudo mv go /usr/local
vim ~/.profile

export GOPATH=$HOME/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

source ~/.profile
cd go-hangman/
```


```
aws ec2 run-instances \
  --image-id <ami-id> \
  --instance-type m5.4xlarge \
  --subnet-id <subnet-id> \
  --key-name <key-pair> \
  --iam-instance-profile Name=<instance-profile>
```


# Sync results back to local

```
 rsync -avt --exclude=".git" --exclude="__pycache__" ubuntu@$ip:/home/ubuntu/mastermind-results.json ./
 ```



 4x7 timing
 17 min: 71 lines (14s/line)
 19 min: 83 lines (14.73s/line, 10s/line in last two minutes)
 41 min: 217 lines (11.33s/line, 9.85s/line in last 22 minutes)
