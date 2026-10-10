# dc

轻量级 Docker 命令行工具，用于查看磁盘占用、并发统计资源数量，以及清理 Build Cache。

## 环境要求

- Go 1.25+
- 已安装并可用 [Docker](https://docs.docker.com/get-docker/)（`docker` 在 PATH 中）

## 安装

```bash
go install github.com/vincentwc/dc@latest
```

或在项目目录本地构建：

```bash
go build -o dc .
```

## 用法

```bash
dc <命令>
```

| 命令 | 说明 |
|------|------|
| `df` | 显示 Docker 磁盘占用（镜像、容器、本地卷等） |
| `stats` | 并发统计镜像总数、运行中容器、数据卷总数 |
| `clean-cache` | 清理 Docker Build Cache（需确认，不可恢复） |

### 示例

```bash
# 磁盘占用
dc df

# 资源统计
dc stats

# 清理 build cache（会提示确认）
dc clean-cache
```

## 开发

```bash
go test ./...
```

## 许可证

未指定许可证时，默认可按仓库维护者意愿使用；如需开源协议可自行补充 `LICENSE` 文件。
