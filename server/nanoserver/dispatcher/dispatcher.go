package dispatcher

import (
	"sync"

	"github.com/asaskevich/EventBus"
)

var (
	eventBus EventBus.Bus
	once     sync.Once
)

// GetInstance 获取事件总线单例实例
func GetInstance() EventBus.Bus {
	once.Do(func() {
		eventBus = EventBus.New()
	})
	return eventBus
}

// MustGetInstance 获取事件总线单例实例，如果未初始化则 panic
func MustGetInstance() EventBus.Bus {
	if eventBus == nil {
		panic("event bus not initialized, call GetInstance() first")
	}
	return eventBus
}
