## usage

- 在 flag_def.go 中定义命令行参数

- 参考 flavor_cli.go，实现自己的客户端测试代码

- 客户端使用方法 (以 flavor_cli 为例)

```bash
./cli.exe -param '{"apikey":"xxx","tenant_id":"t100","platform_userid":"xxx","flavor_id":"1"}' -timeout 3 -method GetFlavor

./cli.exe -param '{"apikey":"xxx","tenant_id":"t100","platform_userid":"xxx","page_number":1,"page_size":1}' -timeout 3 -method ListFlavors
```
