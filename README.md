# classService

## 一、如何运行？

### 1、配置信息

将`configs/config-example.yaml`换成`configs/config.yaml`,并填充配置文件
### 2、构建镜像
在`DockerFile`所在目录下使用命令`docker build -t extra_class:v1`构建镜像
### 3、运行
在`deploy`下执行`docker-compose up -d`即可

## 二、错误码

| 错误码 | 含义 |
|-------|-----|
| 200|成功|
|450|创建classInfo失败|


## 三、API文档
将文件中`openapi.yaml`导入到`apifox`中即可