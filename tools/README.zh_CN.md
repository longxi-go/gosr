# 生成工具适配器

[English](README.md) | [简体中文](README.zh_CN.md)

`github.com/longxi-go/gosr/tools` 为 Go 工具提供可选的包加载和注解适配器。
该模块具有独立的模块路径和版本，将 x/tools、controller-tools 依赖与根模块中
仅依赖标准库的契约和辅助功能分开。

| 包 | 提供的能力 | 调用方控制 |
| --- | --- | --- |
| `pkg/packageload` | 基于 go/packages 加载包，复制配置与 overlay，提供顺序确定的诊断 | 入口、加载模式、标志、环境、模块策略和隔离 |
| `pkg/annotation` | 带源码位置的 controller-tools 标记解析 | 模式、位置、重复及领域校验规则 |

`packageload.Load` 接收显式指定的入口，nil 配置采用 go/packages 默认值。
发现阶段可以接受不完整类型；需要确认包是否存在错误时，须调用 `Errors`。
传入 context 或 overlay 不会隔离文件系统访问或模块下载。

注解解析器处理传入的行文档注释组。前缀必须以 `+` 开头且不含空白；单独使用 `+` 表示选择所有标记。
每个定义的 `Sample` 必须携带类型，无参标记使用 `struct{}{}`。
优先使用类型自身的文档注释，缺少时回退到外层声明；分组声明中的每个选定类型都适用此规则。
适配器不定义领域注解语言或编译流水线。

整数转换沿用 controller-tools 的行为。对于取值范围受限的字段，使用宽整数或字符串样本，
在转换为较窄的整数类型之前检查取值范围。

API 保证及集成边界见[代码生成机制](../docs/design/docs/contracts/codegen.zh.md)。
各包 `example_test.go` 包含可执行样例。

仓库环境配置和源码检查见[参与开发](../docs/design/docs/development.zh.md)。

模块发布流程见[维护指南](../docs/maintainers/releases.zh_CN.md)。
