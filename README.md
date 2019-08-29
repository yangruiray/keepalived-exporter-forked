# keepalived-exporter

监控机器上的keepalived状态，主要包含keepalived进程状态和VIP网络是否正常

暴露数据格式为prometheus metrics

# usage
编译
go build go/src/github.com/keepalived-exporter/cmd/keepalived-exporter/main.go

使用
./main -port
