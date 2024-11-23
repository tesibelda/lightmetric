// lightmetric helps with easy manipulation of simple telegraf metrics
//  without telegraf libraries dependencies
//
// License: The MIT License (MIT)

package metric

import (
	"hash/fnv"
	"sort"
	"time"
)

// Tag represents a single tag key and value.
type Tag struct {
	Key   string
	Value string
}

// Field represents a single field key and value.
type Field struct {
	Key   string
	Value interface{}
}

type Metric struct {
	name   string
	tags   []*Tag
	fields []*Field
	tm     time.Time
}

func New(
	name string,
	tags map[string]string,
	fields map[string]interface{},
	tm time.Time,
) Metric {
	m := Metric{
		name:   name,
		tags:   nil,
		fields: nil,
		tm:     tm,
	}

	if len(tags) > 0 {
		m.tags = make([]*Tag, 0, len(tags))
		for k, v := range tags {
			m.tags = append(m.tags,
				&Tag{Key: k, Value: v})
		}
		sort.Slice(m.tags, func(i, j int) bool { return m.tags[i].Key < m.tags[j].Key })
	}

	if len(fields) > 0 {
		var vi interface{}
		m.fields = make([]*Field, 0, len(fields))
		for k, v := range fields {
			vi = convertField(v)
			if v == nil {
				continue
			}
			m.AddField(k, vi)
		}
	}

	return m
}

// FromMetric returns a deep copy of the metric with any tracking information
// removed.
func FromMetric(other Metric) Metric {
	m := Metric{
		name:   other.Name(),
		tags:   make([]*Tag, len(other.TagList())),
		fields: make([]*Field, len(other.FieldList())),
		tm:     other.Time(),
	}

	for i, tag := range other.TagList() {
		m.tags[i] = &Tag{Key: tag.Key, Value: tag.Value}
	}

	for i, field := range other.FieldList() {
		m.fields[i] = &Field{Key: field.Key, Value: field.Value}
	}
	return m
}

func (m *Metric) Name() string {
	return m.name
}

func (m *Metric) Tags() map[string]string {
	tags := make(map[string]string, len(m.tags))
	for _, tag := range m.tags {
		tags[tag.Key] = tag.Value
	}
	return tags
}

func (m *Metric) TagList() []*Tag {
	return m.tags
}

func (m *Metric) Fields() map[string]interface{} {
	fields := make(map[string]interface{}, len(m.fields))
	for _, field := range m.fields {
		fields[field.Key] = field.Value
	}

	return fields
}

func (m *Metric) FieldList() []*Field {
	return m.fields
}

func (m *Metric) Time() time.Time {
	return m.tm
}

func (m *Metric) SetName(name string) {
	m.name = name
}

func (m *Metric) AddPrefix(prefix string) {
	m.name = prefix + m.name
}

func (m *Metric) AddSuffix(suffix string) {
	m.name += suffix
}

func (m *Metric) AddTag(key, value string) {
	for i, tag := range m.tags {
		if key > tag.Key {
			continue
		}

		if key == tag.Key {
			tag.Value = value
			return
		}

		m.tags = append(m.tags, nil)
		copy(m.tags[i+1:], m.tags[i:])
		m.tags[i] = &Tag{Key: key, Value: value}
		return
	}

	m.tags = append(m.tags, &Tag{Key: key, Value: value})
}

func (m *Metric) HasTag(key string) bool {
	for _, tag := range m.tags {
		if tag.Key == key {
			return true
		}
	}
	return false
}

func (m *Metric) GetTag(key string) (string, bool) {
	for _, tag := range m.tags {
		if tag.Key == key {
			return tag.Value, true
		}
	}
	return "", false
}

func (m *Metric) Tag(key string) string {
	v, _ := m.GetTag(key)
	return v
}

func (m *Metric) RemoveTag(key string) {
	for i, tag := range m.tags {
		if tag.Key == key {
			copy(m.tags[i:], m.tags[i+1:])
			m.tags[len(m.tags)-1] = nil
			m.tags = m.tags[:len(m.tags)-1]
			return
		}
	}
}

func (m *Metric) AddField(key string, value interface{}) {
	for i, field := range m.fields {
		if key == field.Key {
			m.fields[i] = &Field{Key: key, Value: convertField(value)}
			return
		}
	}
	m.fields = append(m.fields, &Field{Key: key, Value: convertField(value)})
}

func (m *Metric) HasField(key string) bool {
	for _, field := range m.fields {
		if field.Key == key {
			return true
		}
	}
	return false
}

func (m *Metric) GetField(key string) (interface{}, bool) {
	for _, field := range m.fields {
		if field.Key == key {
			return field.Value, true
		}
	}
	return nil, false
}

func (m *Metric) Field(key string) interface{} {
	if v, found := m.GetField(key); found {
		return v
	}
	return nil
}

func (m *Metric) RemoveField(key string) {
	for i, field := range m.fields {
		if field.Key == key {
			copy(m.fields[i:], m.fields[i+1:])
			m.fields[len(m.fields)-1] = nil
			m.fields = m.fields[:len(m.fields)-1]
			return
		}
	}
}

func (m *Metric) SetTime(t time.Time) {
	m.tm = t
}

func (m *Metric) Copy() Metric {
	m2 := Metric{
		name:   m.name,
		tags:   make([]*Tag, len(m.tags)),
		fields: make([]*Field, len(m.fields)),
		tm:     m.tm,
	}

	for i, tag := range m.tags {
		m2.tags[i] = &Tag{Key: tag.Key, Value: tag.Value}
	}

	for i, field := range m.fields {
		m2.fields[i] = &Field{Key: field.Key, Value: field.Value}
	}
	return m2
}

func (m *Metric) HashID() uint64 {
	h := fnv.New64a()
	h.Write([]byte(m.name))
	h.Write([]byte("\n"))
	for _, tag := range m.tags {
		h.Write([]byte(tag.Key))
		h.Write([]byte("\n"))
		h.Write([]byte(tag.Value))
		h.Write([]byte("\n"))
	}
	return h.Sum64()
}

// Convert field to a supported type or nil if inconvertible.
func convertField(v interface{}) interface{} {
	switch v := v.(type) {
	case float64:
		return v
	case int64:
		return v
	case string:
		return v
	case bool:
		return v
	case int:
		return int64(v)
	case uint:
		return uint64(v)
	case uint64:
		return v
	case []byte:
		return string(v)
	case int32:
		return int64(v)
	case int16:
		return int64(v)
	case int8:
		return int64(v)
	case uint32:
		return uint64(v)
	case uint16:
		return uint64(v)
	case uint8:
		return uint64(v)
	case float32:
		return float64(v)
	default:
		return nil
	}
}
