# 使用expvarmon查看key訊息狀況
expvarmon -ports="http://localhost:20002" -i 1s -vars="time:key.time,time(max):key.time(max),time(avg):key.time(avg),count:key.count,count(1m):key.count(1m),count(5m):key.count(5m),count(10m):key.count(10m),count(60m):key.count(60m)"

# 使用expvarmon查看json訊息狀況(也包括連線數量)
expvarmon -ports="http://localhost:20002" -i 1s -vars="time:json.time,time(max):json.time(max),time(avg):json.time(avg),count:json.count,count(1m):json.count(1m),count(5m):json.count(5m),count(10m):json.count(10m),count(60m):json.count(60m),connect:connect"

# 使用expvarmon查看proto訊息狀況(也包括連線數量)
expvarmon -ports="http://localhost:20002" -i 1s -vars="time:proto.time,time(max):proto.time(max),time(avg):proto.time(avg),count:proto.count,count(1m):proto.count(1m),count(5m):proto.count(5m),count(10m):proto.count(10m),count(60m):proto.count(60m),connect:connect"

# 使用expvarmon查看plist訊息狀況(也包括連線數量)
expvarmon -ports="http://localhost:20002" -i 1s -vars="time:plist.time,time(max):plist.time(max),time(avg):plist.time(avg),count:plist.count,count(1m):plist.count(1m),count(5m):plist.count(5m),count(10m):plist.count(10m),count(60m):plist.count(60m),connect:connect"

# 使用expvarmon查看記憶體狀況
expvarmon -ports="http://localhost:20002" -i 1s

# 使用圖形化檢視查看記憶體狀況
go tool pprof -http=:8081 http://localhost:20002/debug/pprof/heap