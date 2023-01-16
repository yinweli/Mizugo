rm -r .\messages\

protoc --go_out=. msgid.proto
protoc --go_out=. msgverify.proto

rm -r ..\..\example_server\internal\messages\
rm -r ..\..\example_clientgo\internal\messages\

mkdir ..\..\example_server\internal\messages\
mkdir ..\..\example_clientgo\internal\messages\

copy .\messages\*.* ..\..\example_server\internal\messages\
copy .\messages\*.* ..\..\example_clientgo\internal\messages\