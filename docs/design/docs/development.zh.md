# 开发指南

## 质量门禁 {#quality-gates}

使用 Quality 工作流选定的 Go 和 linter 版本；最低语言版本由 go.mod 声明。
需要安装 Git、Bash、Make、Python 3.12 和 linter：

```sh
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0
```

从 gosr 仓库根目录执行：

```sh
make setup
make license-headers format
make check
make docs-build
```

`setup` 将 tools/go.mod 指定的 golayout 版本下载到模块缓存并验证校验和。
格式化和文件头处理使用该模块的脚本。完整源码检查包括模块整洁性、文件头、结构布局、
格式化回归测试、vet、覆盖率、随机顺序 race 测试、lint 和脚本测试，单独检出本仓库即可执行。

`make check` 在禁用 GOWORK 的情况下检查根模块，并通过临时模块替换，针对当前根模块源码
检查 tools，不修改仓库中的模块文件。`make verify-tools-published` 检查 tools 模块整洁性，
`make check-published` 不使用本地覆盖，针对声明的依赖版本执行完整检查；这些版本须可解析。

执行模块检查前须清除 GOFLAGS。检查会拒绝继承的模块或构建标志，包括 GOENV 持久配置；
已发布依赖检查还会拒绝 tools/go.mod 中的本地目录替换及对根 gosr 依赖的所有替换。

## 格式与源码声明

`make format` 和 `make check-format` 使用 golayout 的结构布局及 gofumpt 组合规则。
非空可执行代码块和复合声明展开，空 `struct{}`、`interface{}`、`select {}` 及函数可保持单行。
文件职责划分和逻辑换行仍需人工审查。生成的 Go 代码与手写代码遵循同一布局策略。

`codequality/HEADER` 提供 Go 文件头声明。`make license-headers` 更新项目维护的文件，
`make check-license-headers` 只读检查。完整的 `testdata/**/*.go.tmpl` 测试样本参与格式化，
文件头更新排除 testdata。备份、缓存、vendor 和符号链接不在处理范围内。

## 本地工作区 {#local-workspace}

需要在独立模块中使用 gosr 源码时，可在共同的父目录创建工作区。此例假定自己的模块位于 `my-app`：

```sh
go work init ./gosr ./my-app
```

已有工作区使用 `go work use` 添加模块。工具导入适配器时须显式选中 `./gosr/tools`，
选择父模块不会包含嵌套模块。开发时选用本地源码，在发布使用方之前独立验证其声明的模块依赖。

go.work 和 go.work.sum 仅保存在本地，正式模块文件不得包含本地目录替换。
工作区用于验证源码集成，已发布依赖检查用于验证模块能否独立解析。

## 文档

使用 Python 3.12 和 uv 构建手册：

```sh
make docs-setup
make docs-serve
make docs-build
```

文档源码位于 `docs/design/docs`，导航及语言配置位于 `docs/design/mkdocs.yml`。
页面采用 `page.md` 和 `page.zh.md` 配对形式；生成的 HTML 和 Python 环境不纳入版本控制。
