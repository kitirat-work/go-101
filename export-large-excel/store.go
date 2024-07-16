package exportlargeexcel

import "fmt"

type IStore interface {
	GetData(limit int, offset int) ([]User, error)
}

type Store struct {
}

// GetData implements IStore.
func (s *Store) GetData(limit int, offset int) ([]User, error) {
	// mock 25000 users
	users := make([]User, limit)
	for i := 0; i < limit; i++ {
		users[i] = User{
			ID:    uint(i),
			Name:  fmt.Sprintf("user-%d", i),
			Email: fmt.Sprintf("user-%d@testmail.com", i),
		}
	}

	if offset == 3 {
		return users[:limit/2], nil
	}

	return users, nil
}

func NewStore() IStore {
	return &Store{}
}
