## usage

- 在 flag_def.go 中定义命令行参数

- 参考 flavor_cli.go，实现自己的客户端测试代码

- 构建
```bash
go build -o flavor.exe flavor_cli.go flag_def.go
```

- 客户端使用方法 (以 flavor_cli 为例)

```bash
./flavor.exe -param '{"apikey":"xxx","tenant_id":"t100","platform_userid":"xxx","flavor_id":"1"}' -timeout 3 -method GetFlavor
```
