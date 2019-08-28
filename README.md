# keepalived-exporter

监控机器上的keepalived状态，主要包含keepalived进程状态和VIP网络是否正常

暴露数据格式为prometheus metrics

# 使用方式
编译
cd go/src/github.com/keepalived-exporter/cmd/keepalived-exporter && go build main.go

运行
./main --port <9999>


