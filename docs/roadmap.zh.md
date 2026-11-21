# Gosr 路线图

[English](roadmap.md) | [简体中文](roadmap.zh.md)

路线图说明已提供的契约领域、辅助库及后续规划的能力。

| 阶段 | 能力与可用结果 | 状态 |
| --- | --- | --- |
| G1 | 可移植 IoC 元数据及独立生命周期契约 | 已实现，[IoC 规范](design/docs/contracts/ioc-contracts.zh.md)定义元数据、查询、所有权和生命周期行为 |
| G2 | 独立排序、类型化同步事件及可选容器通知 | 已实现，[排序与事件规范](design/docs/contracts/events-ordering.zh.md)定义排序与事件投递行为 |
| G3 | 环境值、源、解析器、Profile、快照及类型化绑定契约 | 已实现，[环境规范](design/docs/contracts/environment.zh.md)定义数据和运行时语义 |
| G4 | 可选的有序步骤执行器及职责独立的代码生成辅助功能 | 已实现，[流水线](design/docs/contracts/pipeline.zh.md)及[代码生成](design/docs/contracts/codegen.zh.md)规定可执行行为与所有权 |
| G5 | 可复现的源码与依赖检查、双语规范及独立 root/tools 发布流程 | 已提供工作流，源码检查见[参与开发](design/docs/development.zh.md#quality-gates) |
| G6 | 版本化行为目录及面向独立实现的适配器 | 规划中 |

## 范围与依赖

G1–G3 定义各实现可独立采用的可移植语义。G4 提供职责独立的可选机制，
不定义领域行为或通用编译器框架。

G5 提供源码与依赖检查，以及独立的 root/tools 发布流程。
源码及依赖检查见[参与开发](design/docs/development.zh.md#quality-gates)。

G6 计划将规范要求与面向独立实现的可执行适配器关联。规范定义所需行为，适配器不规定参考运行时。
适配器的覆盖范围按领域划分，容器适配器不覆盖全部事件和环境行为。

新增领域须有实际的共享需求，并单独设计契约。运行时和工具实现的路线图由各项目维护。
