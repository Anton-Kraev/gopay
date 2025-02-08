package mock

import (
	"io"
	"os"

	"github.com/Anton-Kraev/gopay"
)

type FileStorage struct{}

func NewFileStorage() FileStorage {
	return FileStorage{}
}

func (f FileStorage) GetData(_ gopay.ID) ([]byte, error) {
	file, err := os.Open("C:/Users/anton/OneDrive/Рабочий стол/AntonKraev-CV.pdf")
	if err != nil {
		return nil, err
	}

	return io.ReadAll(file)
}
