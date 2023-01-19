# 使用expvarmon監控客戶端記憶體狀況
expvarmon -ports="http://localhost:8082" -i 1s

# 使用expvarmon監控客戶端key狀況
expvarmon -ports="http://localhost:8082" -i 1s -vars="time:key.time,time(max):key.time(max),time(avg):key.time(avg),count:key.count,count(1m):key.count(1m),count(5m):key.count(5m),count(10m):key.count(10m),count(60m):key.count(60m)"

# 使用expvarmon監控客戶端ping狀況(也包括連線數量)
expvarmon -ports="http://localhost:8082" -i 1s -vars="time:ping.time,time(max):ping.time(max),time(avg):ping.time(avg),count:ping.count,count(1m):ping.count(1m),count(5m):ping.count(5m),count(10m):ping.count(10m),count(60m):ping.count(60m),connect:connect"

# 使用expvarmon監控客戶端連線數量
expvarmon -ports="http://localhost:8082" -i 1s -vars="connect:connect"