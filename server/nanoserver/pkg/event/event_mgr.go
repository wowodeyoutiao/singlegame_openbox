package event

import (
	"runtime"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

const (
	DefaultQueueCapacity = 2048 // 默认事件队列容量
	DefaultBatchSize     = 10   // 每批次处理的事件数量
)

// Option 配置选项函数
type Option func(*EventManager)

// WithQueueCapacity 设置队列容量
func WithQueueCapacity(capacity int) Option {
	return func(m *EventManager) {
		if capacity > 0 {
			m.queueCapacity = capacity
			m.eventQueue = make(chan *Event, capacity)
		}
	}
}

// WithBatchSize 设置批处理大小
func WithBatchSize(size int) Option {
	return func(m *EventManager) {
		if size > 0 {
			m.batchSize = size
		}
	}
}

type EventHandler func(data interface{})

type Event struct {
	Type EventType
	Data interface{}
}

type EventManager struct {
	handlers   map[EventType][]EventHandler
	mu         sync.RWMutex
	eventQueue chan *Event
	closed     bool

	queueCapacity int
	batchSize     int
}

func NewEventManager(opts ...Option) *EventManager {
	m := &EventManager{
		handlers:   make(map[EventType][]EventHandler),
		eventQueue: make(chan *Event, DefaultQueueCapacity),
		batchSize:  DefaultBatchSize,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m

}

func (em *EventManager) Subscribe(eventType EventType, handler EventHandler) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.handlers[eventType] = append(em.handlers[eventType], handler)
}

func (em *EventManager) Publish(eventType EventType, data interface{}) {
	em.mu.RLock()
	if em.closed {
		em.mu.RUnlock()
		return
	}
	em.mu.RUnlock()

	select {
	case em.eventQueue <- &Event{Type: eventType, Data: data}:
	default:
		log.Warnf("事件队列已满，丢弃事件 type=%d, data=%+v, queueLen=%d", eventType, data, len(em.eventQueue))
	}
}

func (em *EventManager) ProcessEvents() {
	em.mu.RLock()
	if em.closed || em.eventQueue == nil {
		em.mu.RUnlock()
		return
	}
	// 如果队列为空，直接返回
	if len(em.eventQueue) == 0 {
		em.mu.RUnlock()
		return
	}
	em.mu.RUnlock()

	processed := 0
	for processed < em.batchSize {
		select {
		case event := <-em.eventQueue:
			em.processEvent(event)
			processed++
		default:
			return
		}
	}
}

func (em *EventManager) processEvent(event *Event) {
	em.mu.RLock()
	handlers := em.handlers[event.Type]
	em.mu.RUnlock()

	for _, handler := range handlers {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Errorf("事件处理异常 type=%d, data=%+v, error=%v, stack=%s", event.Type, event.Data, r, em.stack())
				}
			}()
			handler(event.Data)
		}()
	}
}

func (em *EventManager) Close() {
	em.mu.Lock()
	defer em.mu.Unlock()

	if !em.closed {
		em.closed = true
		close(em.eventQueue)
	}
}

func (em *EventManager) Reset() {
	em.mu.Lock()
	em.closed = false
	em.mu.Unlock()
	//em.handlers = make(map[EventType][]EventHandler)
	for len(em.eventQueue) > 0 {
		<-em.eventQueue
	}
}

func (em *EventManager) stack() string {
	buf := make([]byte, 10000)
	n := runtime.Stack(buf, false)
	buf = buf[:n]
	s := string(buf)
	const skip = 7
	count := 0
	index := strings.IndexFunc(s, func(c rune) bool {
		if c != '\n' {
			return false
		}
		count++
		return count == skip
	})
	return s[index+1:]
}
