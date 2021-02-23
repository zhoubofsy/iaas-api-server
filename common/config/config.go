package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var configMap = map[string]string{}

// 配置文件格式:
//   # 注释, 支持单行注释, 不支持行尾注释
//   key = value
func InitConfig(path string) bool {
	//打开文件指定目录，返回一个文件f和错误信息
	f, err := os.Open(path)
	defer f.Close()

	//异常处理 以及确保函数结尾关闭文件流
	if err != nil {
		panic(err)
	}

	//创建一个输出流向该文件的缓冲流*Reader
	r := bufio.NewReader(f)
	for {
		//读取，返回[]byte 单行切片给b
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		//去除单行属性两端的空格
		s := strings.TrimSpace(string(b))

		// # 开头的行是注释, 目前只支持单行注释，不支持行尾注释
		if s == "" || s[0] == '#' {
			continue
		}

		//判断等号=在该行的位置
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		//取得等号左边的key值，判断是否为空
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}

		//取得等号右边的value值，判断是否为空
		value := strings.TrimSpace(s[index+1:])

		// 去掉字符串开头的引号
		value = strings.Trim(value, "'\"`")
		if len(value) == 0 {
			continue
		}

		//把配置文件里的属性key=value对，成功载入到内存中c对象里
		configMap[key] = value
		fmt.Printf("add config %s = %s\n", key, value)
	}
	return true
}

type ConfigError struct {
	Err string
}

func (e ConfigError) Error() string {
	return e.Err
}

func GetBool(key string) (bool, error) {
	val, ok := configMap[key]
	if ok == false {
		return false, ConfigError{"config not found: " + key}
	}

	if val == "true" || val == "True" {
		return true, nil
	}

	if val == "false" || val == "False" {
		return false, nil
	}

	return false, ConfigError{"invalid bool value: " + val}
}

func GetInt(key string) (int, error) {
	val, ok := configMap[key]
	if ok == false {
		return 0, ConfigError{"config not found: " + key}
	}

	v, err := strconv.Atoi(val)
	if err != nil {
		return 0, ConfigError{"invalid int value: " + val}
	}

	return v, nil
}

func GetString(key string) (string, error) {
	val, ok := configMap[key]
	if ok == false {
		return "", ConfigError{"config not found: " + key}
	}

	return val, nil
}
