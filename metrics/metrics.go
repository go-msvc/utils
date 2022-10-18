package metrics

import "sync"

var (
	metrics      = map[string]interface{}{}
	metricsMutex = sync.Mutex{}
)

func All() map[string]interface{} {
	return metrics
}

func Incr(key string) {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()
	if m, ok := metrics[key]; ok {
		if i, ok := m.(int); ok {
			metrics[key] = i + 1
		} else {
			metrics[key] = 1
		}
	} else {
		metrics[key] = 1
	}
}

func Set(key string, value interface{}) {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()
	metrics[key] = value
}
