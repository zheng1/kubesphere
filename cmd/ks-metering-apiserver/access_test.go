package main

import (
	"testing"
	"time"
)

func Test_getMeter(t *testing.T) {
	resp := getMeter(
		[]string{
			"pod/default/ng-cb5d9758-w99xs",
			"pod/default/punk-puma-redis-ha-server-0",
		},
		[]string{
			"pod_cpu_usage",
			"pod_memory_usage",
			"pod_memory_usage_wo_cache",
			"pod_net_bytes_transmitted",
			"pod_net_bytes_received",
		},
		time.Now(),
		time.Now().AddDate(0, 0, 1))
	t.Log(resp)
}
