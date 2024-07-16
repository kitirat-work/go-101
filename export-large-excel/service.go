package exportlargeexcel

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/xuri/excelize/v2"
)

type IService interface {
	Export(ctx context.Context, w http.ResponseWriter) error
}

type Service struct {
	store IStore
}

func (s *Service) Export(ctx context.Context, w http.ResponseWriter) error {
	file := excelize.NewFile()
	streamWriterSheet, err := file.NewStreamWriter("Sheet1")
	if err != nil {
		return err
	}
	headerRow := []interface{}{"Id", "User", "Email"}
	streamWriterSheet.SetRow("A1", headerRow)

	var wg sync.WaitGroup // Declare a WaitGroup
	wg.Add(1)             // Increment the WaitGroup counter
	const limitSize = 10000
	localCtx, cancel := context.WithCancel(ctx)
	dataCh := make(chan []User, 1)
	rowIdx := 2
	go func(ctx context.Context, cancel context.CancelFunc) {
		defer wg.Done() // Decrement the counter when the goroutine completes
		defer close(dataCh)
		for {
			select {
			case <-ctx.Done():
				return
			case users := <-dataCh:
				for _, user := range users {
					userRow := []interface{}{
						user.ID,
						user.Name,
						user.Email,
					}
					err := streamWriterSheet.SetRow(fmt.Sprintf("A%d", rowIdx), userRow)
					if err != nil {
						cancel()
					}
					rowIdx++
				}
				if len(users) < limitSize {
					return
				}
			}
		}

	}(localCtx, cancel)

	for i := 0; ; i++ {
		users, err := s.store.GetData(limitSize, i)
		if err != nil {
			cancel()
			return err
		}

		dataCh <- users

		if len(users) < limitSize {
			break
		}
	}
	wg.Wait() // Wait for the goroutine to finish
	err = streamWriterSheet.Flush()
	if err != nil {
		return err
	}
	err = file.Write(w)
	return err
}

func NewService(iStore IStore) IService {
	return &Service{
		store: iStore,
	}
}
