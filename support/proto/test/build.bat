echo off

REM Generate the message for GO
echo #### Generate the message for GO
set source=msg-go
set targets=..\..\..\support\test-server\internal\msgs
set targetc=..\..\..\support\test-client-go\internal\msgs

rm -r %source%\msgs
protoc --go_out=%source% msgid.proto
protoc --go_out=%source% msgtest.proto

rm -r %targets%
mkdir %targets%
copy %source%\msgs\*.* %targets%
copy %source%\msgs-json\*.* %targets%

rm -r %targetc%
mkdir %targetc%
copy %source%\msgs\*.* %targetc%
copy %source%\msgs-json\*.* %targetc%

REM Generate the message for Unity
echo #### Generate the message for Unity
set source=msg-cs
set targetc=..\..\..\support\test-client-cs\internal\msgs

rm -r %source%\msgs
mkdir %source%\msgs
protoc --csharp_out=%source%\msgs msgid.proto
protoc --csharp_out=%source%\msgs msgtest.proto

rm -r %targetc%
mkdir %targetc%
copy %source%\msgs\*.* %targetc%
copy %source%\msgs-json\*.* %targetc%

REM Generate the message for Unity Test
echo #### Generate the message for Unity Test
set source=msg-cs
set targetc=..\..\..\support\client-unity\Packages\com.fouridstudio.mizugo-client-unity\Tests\Runtime\msgs

rm -r %targetc%
mkdir %targetc%
copy %source%\msgs\*.* %targetc%
copy %source%\msgs-json\*.* %targetc%