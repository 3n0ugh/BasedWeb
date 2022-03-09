package mock

import (
	"github.com/3n0ugh/BasedWeb/internal/data"
)

func NewModel() data.Model {
	return data.Model{
		Blog: BlogModel{},
	}
}
