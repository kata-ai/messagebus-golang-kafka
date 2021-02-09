package serialization

import (
	"errors"
	"fmt"
)

type SubjectStrategy int

const (
	TOPIC_NAME_STRATEGY = iota
	TOPIC_RECORD_NAME_STRATEGY
	RECORD_NAME_STRATEGY
)

func (s SubjectStrategy) String() string {
	return [...]string{"TOPIC_NAME_STRATEGY", "TOPIC_RECORD_NAME_STRATEGY", "RECORD_NAME_STRATEGY"}[s]
}

func prepareSubjectName(topic string, schemaStr string, strategy SubjectStrategy, isKey bool) (string, error) {
	schema := Schema{schema: schemaStr}
	switch strategy {
	case TOPIC_NAME_STRATEGY:
		if isKey {
			return fmt.Sprintf("%s-key", topic), nil
		}
		return fmt.Sprintf("%s-value", topic), nil
	case TOPIC_RECORD_NAME_STRATEGY:
		return fmt.Sprintf("%s-%s", topic, schema.FullName()), nil
	case RECORD_NAME_STRATEGY:
		return schema.FullName(), nil
	default:
		return "", errors.New("unknown subject strategy")
	}
}
