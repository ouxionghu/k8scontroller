# Kubernetes Controller

这是一个用于开发 Kubernetes Controller 的项目模板。该项目使用 Go 语言和 controller-runtime 框架来构建自定义的 Kubernetes Controller。

## 项目结构

```
.
├── cmd/
│   └── controller/        # 主程序入口
├── pkg/
│   ├── apis/             # API 定义
│   │   └── example/
│   │       └── v1/       # API 版本
│   └── controller/       # Controller 实现
└── go.mod               # Go 模块文件
```

## 开发环境要求

- Go 1.21 或更高版本
- Kubernetes 集群（用于测试）
- kubectl 配置正确

## 构建和运行

1. 下载依赖：
```bash
go mod download
```

2. 构建项目：
```bash
go build -o controller cmd/controller/main.go
```

3. 运行 controller：
```bash
./controller
```

## 配置

Controller 支持以下命令行参数：

- `--metrics-bind-address`: 指标服务绑定地址（默认：:8080）
- `--health-probe-bind-address`: 健康检查服务绑定地址（默认：:8081）
- `--leader-elect`: 是否启用 leader 选举（默认：false）

## 开发指南

1. 在 `pkg/apis/example/v1` 目录下定义你的自定义资源（CRD）
2. 在 `pkg/controller` 目录下实现你的 controller 逻辑
3. 使用 `make generate` 生成代码（需要安装 code-generator）
4. 构建并部署到 Kubernetes 集群

## 许可证

MIT 