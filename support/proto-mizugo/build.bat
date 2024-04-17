echo off

REM Generate the message for GO
echo #### Generate the message for GO
set source=msg-go
set target=..\..\mizugos\msgs

rm -r %source%\msgs
protoc --go_out=%source% protomsg.proto
protoc --go_out=%source% raven.proto

rm -r %target%
mkdir %target%
copy %source%\msgs\*.* %target%
copy %source%\msgs-json\*.* %target%

REM Generate the message for Unity
echo #### Generate the message for Unity
set source=msg-cs
set target=..\..\support\client-unity\Packages\com.fouridstudio.mizugo-client-unity\Runtime\Msgs

rm -r %source%\msgs
mkdir %source%\msgs
protoc --csharp_out=%source%\msgs protomsg.proto
protoc --csharp_out=%source%\msgs raven.proto

rm -r %target%
mkdir %target%
copy %source%\msgs\*.* %target%
copy %source%\msgs-json\*.* %target%