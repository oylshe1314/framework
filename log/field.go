package log

type fieldLogger struct {
	entry

	fields map[string]any
}

func (this *fieldLogger) WithField(key string, value any) entry {
	this.fields[key] = value
	return this
}
