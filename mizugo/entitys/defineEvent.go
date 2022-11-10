package entitys

const eventBufferSize = 1000 // 事件緩衝區大小

// eventAwake awake事件
type eventAwake struct {
    module IModule // 模組物件
}

// eventStart start事件
type eventStart struct {
    module IModule // 模組物件
}

// eventDispose dispose事件
type eventDispose struct {
    module IModule // 模組物件
}
