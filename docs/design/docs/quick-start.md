# Quick start

Use Go 1.27.0 or later. Import the packages for the capabilities your application
uses from the `github.com/longxi-go/gosr` module. A contract describes the behavior
expected from an implementation; importing it does not create a container,
publisher or configuration engine.

## Express an ordering capability

The implementing object supplies its order; the consumer owns sampling and sorting.

```go
package application

import "github.com/longxi-go/gosr/pkg/ordered"

type Audit struct{}

func (Audit) Order() int {
    return 100
}

var _ ordered.Ordered = Audit{}
```

`PriorityOrdered` sorts before ordinary `Ordered`, independently of numeric order.
`HighestPrecedence` and `LowestPrecedence` bound the native-int range. Unordered
objects use lowest precedence; equal keys follow the consuming API's tie rule.
See [ordering and events](contracts/events-ordering.md) for sampling and concurrency
requirements.

## Choose an implementation

Import the contract package that describes your dependency. An application can
accept an `event.Publisher` from any implementation that provides the specified
matching, delivery and cancellation behavior. An environment consumer can accept
an `environment.SnapshotProvider` without knowing its storage or refresh mechanism.
IoC lookup, application lifecycle and scope management are separate capabilities;
implement only those you need and satisfy their associated requirements.

For example, a producer accepts the publication capability and passes an explicit
event type:

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

`PublishEvent` delivers synchronously; the call returns when delivery completes or
fails. `Notify` passes publication errors to its caller. The publisher implementation
owns listener matching and delivery; the producer depends only on `event.Publisher`.
See the [event contract](contracts/events-ordering.md) for ordering, failure and
concurrency semantics.

The optional pipeline/codegen helpers and the separate `gosr/tools` module are
used by generators that need those mechanisms. They are not required for ordering,
events, environment or IoC interface use.
