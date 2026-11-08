# nanoserver

### 环境
Golang版本：go1.21
```
go env -w GOPROXY=https://goproxy.cn,direct


apt install mysql-server
ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'mow123456';
CREATE DATABASE game_mow DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;

apt install redis-server

apt install nginx

```

## 配置

- 数据库配置[database]
SET GLOBAL sql_mode='STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'
- 端口配置[webserver][game-server]

## 运行

go mod tidy

go mod vendor

go run main.go

## GM
添加道具: @add [道具id] [数量]
清空包裹: @clear
满: @god
