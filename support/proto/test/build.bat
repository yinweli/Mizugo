echo off
set test_server=..\..\test-server\internal\messages\
set test_clientgo=..\..\test-client-go\internal\messages\

rm -r .\messages\
protoc --go_out=. msgid.proto
protoc --go_out=. msgtest.proto

rm -r %test_server%
mkdir %test_server%
copy .\messages\*.* %test_server%
copy .\messagesjson\*.* %test_server%

rm -r %test_clientgo%
mkdir %test_clientgo%
copy .\messages\*.* %test_clientgo%
copy .\messagesjson\*.* %test_clientgo%