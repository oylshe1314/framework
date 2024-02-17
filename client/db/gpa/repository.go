package gpa

type Entity[ID any] interface {
	Id() ID
}

type Repository[ID any, E Entity[ID]] interface {
	FindById(id ID) (E, error)
	FindAll() ([]E, error)
}

type simpleRepository[ID any, E Entity[ID]] struct {
}

func (this *simpleRepository[ID, E]) FindById(id ID) (e E, err error) {
	return
}

func (this *simpleRepository[ID, E]) FindAll() (es []E, err error) {
	return
}
