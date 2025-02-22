package links

import (
	"errors"
	"fmt"

	"github.com/Anton-Kraev/gopay"
)

var errGenerateLink = errors.New("generate link failed")

type Generator struct {
	baseURL string
}

func NewGenerator(baseURL string) Generator {
	return Generator{baseURL: baseURL}
}

func (g Generator) GenerateLink(id gopay.ID) (gopay.Link, error) {
	link := gopay.Link(fmt.Sprintf("%s/api/%s", g.baseURL, id))
	if !link.Validate() {
		return "", errGenerateLink
	}

	return link, nil
}
