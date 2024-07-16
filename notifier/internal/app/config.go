package app

import (
	"os"
	"strings"
)

type Config struct {
	bootstrapServers []string
	topics           []string
	groupId          string
}

func NewConfig() Config {
	cfg := Config{}
	cfg.parse()
	return cfg
}

func (c *Config) parse() {
	parsedTopics := strings.Split(getEnvHelper("KAFKA_TOPICS", "loms.order-events"), ";")
	parsedBrokers := strings.Split(getEnvHelper("KAFKA_BOOTSTRAP", "localhost:9092"), ";")

	c.bootstrapServers = make([]string, 0, len(parsedBrokers))
	c.topics = make([]string, 0, len(parsedTopics))

	for _, broker := range parsedTopics {
		if broker != "" {
			c.topics = append(c.topics, broker)
		}
	}

	for _, server := range parsedBrokers {
		if server != "" {
			c.bootstrapServers = append(c.bootstrapServers, server)
		}
	}
	c.groupId = getEnvHelper("KAFKA_GROUP_ID", "notify")
}

func getEnvHelper(envVar, defaultVal string) string {
	val := os.Getenv(envVar)
	if val == "" {
		val = defaultVal
	}
	return val
}
