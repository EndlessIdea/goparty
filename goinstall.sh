#!/bin/bash
#go环境一键安装

#默认下载1.9.1版本
url="https://dl.google.com/go/go1.9.1.linux-amd64.tar.gz"
filename="go.tar.gz"
if [ $# -gt 0 ]
then
    url=$1
fi

curl -o $filename $url
tar -zxf $filename
rm -rf $filename
rm -rf /usr/local/go_bak
mv /usr/local/go /usr/local/go_bak
mv go /usr/local
echo 'export GOROOT=/usr/local/go' >> ~/.bash_profile
echo 'export GOPATH=/home/funkycoder/gopath' >> ~/.bash_profile
echo 'PATH=$PATH:$GOROOT/bin:$GOPATH/bin' >> ~/.bash_profile

echo 'done'