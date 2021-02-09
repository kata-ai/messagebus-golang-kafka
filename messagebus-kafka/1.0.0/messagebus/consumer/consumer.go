package consumer

import "github.com/confluentinc/confluent-kafka-go/kafka"

type ConsumerConfiguration struct {
	PollIntervalMs int
	KafkaConfig    *kafka.ConfigMap
}

type ConsumerOption func(c *ConsumerConfiguration)

// Create instance of consumer configuration
// Kafka configuration customization can be done through variadic parameters
// Complete Kafka configuration can be seen at https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
// Example:
// 		NewConsumerConfig("group-1", WithPollIntervalMs(150), WithFetchMinBytes(20))
// Default values:
// 		pollIntervalMs: 100
// 		fetch.min.bytes: 10
// 		fetch.wait.max.ms: 10
// 		max.partition.fetch.bytes: 1048576
// 		session.timeout.ms: 10000
// 		heartbeat.interval.ms: 3000
// 		enable.auto.commit: false
// 		auto.offset.reset: "earliest"
func NewConsumerConfig(groupId string, opts ...ConsumerOption) *ConsumerConfiguration {
	consumerConfig := &ConsumerConfiguration{
		PollIntervalMs: 100,
		KafkaConfig: &kafka.ConfigMap{
			"group.id":                  groupId,
			"fetch.min.bytes":           10,
			"fetch.wait.max.ms":         500,
			"max.partition.fetch.bytes": 1048576,
			"session.timeout.ms":        10000,
			"heartbeat.interval.ms":     3000,
			"enable.auto.commit":        false,
			"auto.offset.reset":         "earliest",
		},
	}
	for _, opt := range opts {
		opt(consumerConfig)
	}
	return consumerConfig
}

// Configure subscription poll interval in millisecond
func WithPollIntervalMs(ms int) ConsumerOption {
	return func(c *ConsumerConfiguration) {
		c.PollIntervalMs = ms
	}
}

// Configure minimum number of bytes the broker responds with
func WithFetchMinBytes(fetchMinBytes int) ConsumerOption {
	return func(c *ConsumerConfiguration) {
		_ = c.KafkaConfig.SetKey("fetch.min.bytes", fetchMinBytes)
	}
}

// Configure maximum time the broker may wait to fill the
// fetch response with fetch.min.bytes of messages
func WithFetchWaitMaxMs(fetchWaitMaxMs int) ConsumerOption {
	return func(c *ConsumerConfiguration) {
		_ = c.KafkaConfig.SetKey("fetch.wait.max.ms", fetchWaitMaxMs)
	}
}

// Configure initial maximum number of bytes per topic+partition
// to request when fetching messages from the broker
func WithMaxPartitionFetchBytes(bytes int) ConsumerOption {
	return func(c *ConsumerConfiguration) {
		_ = c.KafkaConfig.SetKey("max.partition.fetch.bytes", bytes)
	}
}

// Configure client group session and failure detection timeout
func WithSessionTimoutMs(ms int) ConsumerOption {
	return func(c *ConsumerConfiguration) {
		_ = c.KafkaConfig.SetKey("session.timeout.ms", ms)
	}
}

// Configure group session keepalive heartbeat interval
func WithHeartbeatIntervalMs(ms int) ConsumerOption {
	return func(c *ConsumerConfiguration) {
		_ = c.KafkaConfig.SetKey("heartbeat.interval.ms", ms)
	}
}

// Configure cction to take when there is no initial offset
// in offset store or the desired offset is out of range
func WithAutoOffsetReset(autoReset string) ConsumerOption {
	return func(c *ConsumerConfiguration) {
		_ = c.KafkaConfig.SetKey("auto.offset.reset", autoReset)
	}
}

// Configure whether to automatically and periodically commit
// offsets in the background
func WithEnableAutoCommit(autoCommit bool) ConsumerOption {
	return func(c *ConsumerConfiguration) {
		_ = c.KafkaConfig.SetKey("enable.auto.commit", autoCommit)
	}
}
