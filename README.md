## 部署说明

```bash
cp config.example.yaml config.yaml
cp Makefile-example Makefile
```

## 本地开发

- air 命令可以热更新。
- 启动命令是 make，需要热更新写代码，使用 air 命令。使用 make 命令，需要改动代码后，手动停止命令并启动 http 服务。

注意：<s>报错了没有提示，看不出来因为报错导致请求失败</s>已更改 .air.toml，失败时会显示报错

相关命令参考 Makefile

```bash
git clone git@codeup.aliyun.com:cblink/flow/china-life/data-distribute/data-distribute-api.git data-distribute-api

make
```

## 目录结构

```
.
├── Dockerfile                                                              # 容器化 Dockerfile
├── Makefile                                                                # 辅助命令，请查看有什么可以执行的命令
├── Makefile-example                                                        # Makefile example 文件
├── README.md                                                               # README 自叙文件
├── cmd                                                                     # 命令行、入口目录
│   ├── main.go                                                             # 入口文件
│   └── server                                                              # server 命令所在目录
│       └── cmd.go                                                          # server 命令，初始化配置、启动引导。请手动控制第三方服务的引导启动顺序，api 接口定义也在这里
├── config.example.yaml                                                     # 配置文件参考
├── config.yaml                                                             # 配置文件
├── go.mod                                                                  # go module 文件
├── go.sum                                                                  # go sum 文件
├── internal                                                                # 内部目录
│   ├── handlers                                                            # 类似于 MVC 的控制器
│   │   ├── common.go                                                       # 公共结构体、函数
│   │   ├── data.go                                                         # 控制器作用，handle 请求处理，data 数据请求相关接口
│   │   ├── qc_task.go                                                      # 控制器作用，handle 请求处理
│   │   └── qc_task_sample.go                                               # 控制器作用，handle 请求处理
│   ├── initialization                                                      # 项目初始化相关文件
│   │   ├── config.go                                                       # 初始化配置，读取 yaml 配置文件
│   │   └── db.go                                                           # 初始化 db 连接
│   └── models                                                              # models 模型，gorm 模型
│       ├── qc_task.go                                                      # 模型 qc_task
│       └── qc_task_sample.go                                               # 模型 qc_task_sample
├── pkg                                                                     # 第三方服务扩展封装
│   └── rabbitmq                                                            # rabbitmq 队列
│       └── client.go                                                       # rabbitmq 队列的 client，可以发送队列、消费消息
└── tmp                                                                     # 临时目录，通过 air 编译、热更新。
    └── main                                                                # 实时编译的入口文件
```
