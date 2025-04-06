package links

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/Anton-Kraev/gopay"
)

var errGenerateLink = errors.New("generate link failed")

type Generator struct {
	baseURL string
}

func NewGenerator(baseURL string) Generator {
	return Generator{baseURL: baseURL}
}

func (g Generator) GenerateLink() (gopay.ID, gopay.Link, error) {
	const op = "links.Generator.GenerateLink"

	id, err := uuid.NewRandom()
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	link := gopay.Link(fmt.Sprintf("%s/api/%s", g.baseURL, id))
	if !link.Validate() {
		return "", "", fmt.Errorf("%s: %w", op, errGenerateLink)
	}

	return gopay.ID(id.String()), link, nil
}
