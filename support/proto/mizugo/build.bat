set mizugo=..\..\..\mizugos\msgs\

rm -r .\msgs\
protoc --go_out=. protomsg.proto
protoc --go_out=. stackmsg.proto

mkdir %mizugo%
copy .\msgs\*.* %mizugo%