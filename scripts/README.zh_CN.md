# 仓库脚本

[English](README.md) | [简体中文](README.zh_CN.md)

这些 Python 脚本用于仓库检查、共享格式化和发布校验。从仓库根目录运行 Make 目标。
使用 Python 3.12，以及同时满足 [go.mod](../go.mod) 和
[tools/go.mod](../tools/go.mod) 要求的 Go 工具链。格式化还需要 Git 和 Bash。
GitHub 检查和发布需要已认证的 GitHub CLI（`gh`）。

## 入口

| 脚本 | 职责 | Make 目标 |
| --- | --- | --- |
| `golayout.py` | 解析 tools 声明的 golayout 版本，调用其共享脚本 | `setup`、`format`、`check-format`、`test-format`、`license-headers`、`check-license-headers` |
| `tools.py` | 按源码或已发布依赖策略执行 tools 模块命令 | `test`、`race`、`vet`、`lint`、`verify-tools-source`、`verify-tools-published` |
| `release.py` | 校验模块标签、身份及发布说明；可选检查依赖的根模块 Release | `validate-release`、`prepare-tools-release` |
| `publish_release.py` | 核对已检查的标签提交，创建或复用内容一致的 Release，处理 Latest 延迟并检查 Pages 前提 | Release 和 Pages 工作流 |
| `go_module.py` | 通过 `go mod edit -json` 读取模块元数据，供其他脚本调用 | 无 CLI 入口 |

常用检查：

```sh
make setup
make check-format
make test-scripts
```

`make format` 重写项目维护的 Go 源码；`make license-headers` 更新其文件头。
仓库环境和质量检查见[参与开发](../docs/design/docs/development.zh.md)。
标签约定、依赖准备和发布顺序见[发布指南](../docs/maintainers/releases.zh_CN.md)。

## 模块与环境策略

`GO` 选择 Go 可执行文件，默认为 `go`。格式化采用根 go.mod 声明的语言版本，
golayout 依赖版本由 tools 声明。

`tools.py` 在 tools 目录执行命令：

- `source` 使用临时模块文件，将根模块依赖替换为当前源码。
  Go 检查通过命令行参数接收 `-modfile`；外部工具和共享格式化脚本使用私有工作区，
  避免将该标志传入它们启动的 Go 子进程。调用结束后清理临时文件。
- `published` 禁用工作区，使用声明的依赖。禁止本地目录替换，也禁止对根模块的任何替换，
  包括版本替换。此模式下，`go mod tidy` 等命令可以更新 tools/go.mod 和 tools/go.sum。

两种模式均拒绝非空的 `GOFLAGS`，包括通过 `GOENV` 持久配置的标志。
该策略控制模块解析，不隔离被执行命令的文件访问，也不禁用依赖下载。

元数据读取会禁用工作区、用户标志、持久 Go 配置及自动工具链选择，
仅解析指定模块文件，不解析依赖图。因此，本地发布校验通过不代表依赖已经发布。

## 测试与外部命令

`make test-scripts` 运行 Python unittest 测试集，需要真实 Go，
不需要 GitHub 认证或已发布的 gosr 根版本。
部分外部操作使用 `tests/testdata` 下的测试替身：

| 测试替身 | 行为 |
| --- | --- |
| `gh` | 检查完整请求参数，返回场景指定的 Release 数据或错误，不向 GitHub 发送请求。 |
| `golayout/go` | 仅拦截 golayout 下载请求；使用真实 Go 检查模块环境，其余命令均转交真实 Go。 |
| `golayout/format.sh` | 检查包装脚本的参数、模块路径和环境，记录语言版本，不执行源码格式化。 |

测试通过子进程环境中的 `GO` 或 `PATH` 选择这些替身，用于验证包装脚本及其失败处理。
共享格式化行为通过 `make test-format` 和 `make check-format` 调用真实 golayout 检查；
已发布依赖的解析须通过 `make check-published` 验证。

发布器测试在外部 API 边界提供受控 HTTP 响应，覆盖 Latest 延迟、已有 Release 内容一致
或冲突、根模块与 tools 的 Latest 策略、创建结果不确定、标签提交不一致等情况，不执行
网络写入。`--check-pages` 和 `--require-latest` 为只读模式，校验选定根源码并查询远程仓库
设置或已发布 Release。发布模式必须提供已通过源码检查的完整提交 SHA。
