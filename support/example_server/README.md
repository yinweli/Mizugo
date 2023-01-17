# 使用expvarmon監控伺服器記憶體狀況
expvarmon -ports="http://localhost:8081" -i 1s

# 使用expvarmon監控伺服器ping狀況
expvarmon -ports="http://localhost:8081" -i 1s -vars="time:ping.time,time(max):ping.time(max),time(avg):ping.time(avg),count:ping.count,count(1m):ping.count(1m),count(5m):ping.count(5m),count(10m):ping.count(10m),count(60m):ping.count(60m)"