package tracing

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
)

func TestSanitizeAttributesReplacesInvalidUTF8Strings(t *testing.T) {
	attrs := sanitizeAttributes([]attribute.KeyValue{
		{Key: attribute.Key("db.\xffstatement"), Value: attribute.StringValue("select \xff")},
		attribute.StringSlice("db.memcached.item", []string{"valid", "bad \xfe"}),
		attribute.Int("http.status_code", 200),
	})

	assert.True(t, utf8.ValidString(string(attrs[0].Key)))
	assert.Equal(t, attribute.Key("db.?statement"), attrs[0].Key)
	assert.True(t, utf8.ValidString(attrs[0].Value.AsString()))
	assert.Equal(t, "select ?", attrs[0].Value.AsString())
	assert.True(t, utf8.ValidString(attrs[1].Value.AsStringSlice()[1]))
	assert.Equal(t, []string{"valid", "bad ?"}, attrs[1].Value.AsStringSlice())
	assert.Equal(t, int64(200), attrs[2].Value.AsInt64())
}
