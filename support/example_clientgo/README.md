# 使用expvarmon監控客戶端記憶體狀況
expvarmon -ports="http://localhost:8081" -i 1s

# 使用expvarmon監控客戶端echo狀況
expvarmon -ports="http://localhost:8081" -i 1s -vars="time:echo.time,time(max):echo.time(max),time(avg):echo.time(avg),count:echo.count,count(1m):echo.count(1m),count(5m):echo.count(5m),count(10m):echo.count(10m),count(60m):echo.count(60m)"

# 使用expvarmon監控客戶端ping狀況
expvarmon -ports="http://localhost:8081" -i 1s -vars="time:ping.time,time(max):ping.time(max),time(avg):ping.time(avg),count:ping.count,count(1m):ping.count(1m),count(5m):ping.count(5m),count(10m):ping.count(10m),count(60m):ping.count(60m)"