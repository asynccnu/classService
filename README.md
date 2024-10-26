# classService

## 一、如何运行？

### 1、配置信息

将`configs/config-example.yaml`换成`configs/config.yaml`,并填充配置文件
### 2、构建镜像
在`DockerFile`所在目录下使用命令`docker build -t extra_class:v1`构建镜像
>构建镜像，是需要拉取golang:1.22和debian:stable-slim这两个镜像的，当然，如果你是在自己机子上跑，挂个梯子就可以拉取这两个镜像了，但是如果你是在云服务器上拉取的话，很有可能拉取不了（被墙），这是你可以尝试过构建自己的阿里云镜像仓库，然后现在自己的机子上拉取那两个镜像，然后改下tag，上传至自己的阿里云的镜像仓库，然后你的服务器就可以从你自己的阿里云镜像仓库中拉取这两个镜像了
>
>参考教程如下:
>
>[如何构建自己的阿里云镜像仓库](https://blog.csdn.net/qq_26709459/article/details/128726699)


### 3、运行
在`deploy`下执行`docker-compose up -d`即可

## 二、错误码

| 错误码 | 含义 |
|-------|-----|
| 200|成功|
|450|创建classInfo失败|


## 三、API文档
将文件中`openapi.yaml`导入到`apifox`中即可