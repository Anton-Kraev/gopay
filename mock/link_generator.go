package mock

import (
	"fmt"

	"github.com/Anton-Kraev/gopay"
)

type LinkGenerator struct{}

func NewLinkGenerator() LinkGenerator {
	return LinkGenerator{}
}

func (l LinkGenerator) GenerateLink(id gopay.ID) (gopay.Link, error) {
	return gopay.Link(fmt.Sprintf("localhost:1323/%s", id)), nil
}
