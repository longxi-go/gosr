# 快速开始

使用 Go 1.27.0 或更高版本，从 `github.com/longxi-go/gosr` 模块导入应用所需能力对应的包。
契约规定实现方应提供的行为；导入契约不会创建容器、发布器或配置引擎。

## 声明排序能力

对象提供顺序值，由使用方负责采样和排序：

```go
package application

import "github.com/longxi-go/gosr/pkg/ordered"

type Audit struct{}

func (Audit) Order() int {
    return 100
}

var _ ordered.Ordered = Audit{}
```

`PriorityOrdered` 组始终排在普通 `Ordered` 之前，组间优先级不受数字大小影响。
`HighestPrecedence` 和 `LowestPrecedence` 表示当前平台 int 的两个边界。
未声明排序的对象采用最低优先级；排序值相同时，按使用方的具体契约处理。
采样和并发要求见[排序与事件](contracts/events-ordering.zh.md)。

## 选择实现

导入描述所需能力的契约包。应用可以接受任何满足规定的匹配、投递和取消行为的
`event.Publisher` 实现。环境使用方可以接受 `environment.SnapshotProvider`，
无需了解其存储或刷新机制。IoC 查询、应用生命周期和作用域管理是独立能力，
按需实现并满足各自的要求。

例如，事件生产方接收发布能力，并传入明确的事件类型：

```go
package application

import (
    "context"

    "github.com/longxi-go/gosr/pkg/event"
)

type InvoicePosted struct {
    event.BaseEvent
    InvoiceID string
}

func Notify(ctx context.Context, publisher event.Publisher, invoiceID string) error {
    return publisher.PublishEvent(ctx, InvoicePosted{
        InvoiceID: invoiceID,
    })
}
```

`PublishEvent` 同步投递，调用在投递完成或失败后返回。`Notify` 将发布错误交给调用方处理。
发布器实现负责监听器匹配和投递，生产方仅依赖 `event.Publisher`。
排序、失败和并发语义见[事件契约](contracts/events-ordering.zh.md)。

需要相关机制的生成器可使用 pipeline/codegen 辅助功能及独立的 `gosr/tools` 模块。
使用排序、事件、环境或 IoC 接口不需要这些工具包。
