package features

import (
	"sync/atomic"
)

var JsonCounter atomic.Int64  // Json計數器
var ProtoCounter atomic.Int64 // Proto計數器
var RavenCounter atomic.Int64 // Raven計數器
