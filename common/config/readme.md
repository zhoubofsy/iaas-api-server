## usage

- 初始化

在 main() 函数中调用 config.InitConfig() 初始化:
```go
import "iaas-api-server/common/config"

config.InitConfig("xx.conf")
```

- 配置文件格式
  
```bash
# 配置文件 xx.conf
# # 开头的行为注释

# bool 类型, 取值范围: false, False, true, True
boo = false

# 字符串类型, 两边引号可选, 单引号、双引号、反引号均可
s = "hello world"

# 整数类型
i = 32
```

- 用法

```go
import "iaas-api-server/common/config"

b, err := config.GetBool("boo")
s, err := config.GetString("s")
i, err := config.GetInt("i")
```
