package loaders

import (
	"sync"

	"github.com/Popov-Dmitriy-Ivanovich/genmilk_backend/models"
)

const MAX_CONCURENT_LOADERS = 10

type Loader interface {
	Init(chan any)
	Run()
	Terminate()
	CollectErrors() []error
}

type ModelLoader struct {
	LoaderFunc    func(model any) error
	ModelsChannel chan any
	Errors        []error
	ErrorsMtx     sync.Mutex
}

// Run
// Blocks until all loaders finish their tasks
func (ml *ModelLoader) Run() {
	wg := sync.WaitGroup{}
	for i := 0; i < MAX_CONCURENT_LOADERS; i++ {
		wg.Add(1)
		go func() {
			for model := range ml.ModelsChannel {
				if err := ml.LoaderFunc(model); err != nil {
					ml.ErrorsMtx.Lock()
					ml.Errors = append(ml.Errors, err)
					ml.ErrorsMtx.Unlock()
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func (ml *ModelLoader) Terminate() {
	close(ml.ModelsChannel)
}

func (ml *ModelLoader) CollectErrors() []error {
	return ml.Errors
}

type LoadModelToDbError struct{}

func (lmtde LoadModelToDbError) Error() string {
	return "внутренняя ошибка сервра, LoadModelToDb получил не правильный тип"
}

func LoadModelToDb[modelType any]() func(model any) error {
	return func(model any) error {
		db := models.GetDb()
		typedModel, ok := model.(modelType)
		if !ok {
			return LoadModelToDbError{}
		}
		if err := db.FirstOrCreate(&typedModel).Error; err != nil {
			return err
		}
		if err := db.Save(&typedModel).Error; err != nil {
			return err
		}
		return nil
	}
}
