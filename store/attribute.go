package store

type Attribute interface {
	Store[string, any]
}

type AttributeProvider interface {
	Attribute() Attribute
}

func NewAttribute() Attribute {
	return NewConcurrent[string, any]()
}
