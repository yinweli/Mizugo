protoc --go_out=. messageid.proto
protoc --go_out=. msgkey.proto
protoc --go_out=. msgping.proto

copy .\messages\*.* ..\..\example_server\internal\messages\
copy .\messages\*.* ..\..\example_clientgo\internal\messages\