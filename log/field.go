package log

type field struct {
	name  string
	value any
}

type fieldLogger struct {
	entry

	fields []*field
}

func (this *fieldLogger) WithField(key string, value any) entry {
	this.fields = append(this.fields, &field{key, value})
	return this
}
