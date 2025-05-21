ps -ef | grep '/main' | grep -v grep | awk '{print $2}' | xargs kill -9
# sleep 2
# go run cmd/server/main.go

