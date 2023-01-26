set mizugo=..\..\..\mizugos\msgs\

rm -r .\msgs\
protoc --go_out=. protomsg.proto
protoc --go_out=. plistmsg.proto

mkdir %mizugo%
copy .\msgs\*.* %mizugo%