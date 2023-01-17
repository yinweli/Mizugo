set mizugo=..\..\..\mizugos\msgs\

rm -r .\msgs\
protoc --go_out=. complexmsg.proto
protoc --go_out=. protomsg.proto

rm -r %mizugo%
mkdir %mizugo%
copy .\msgs\*.* %mizugo%