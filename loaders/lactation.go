package loaders

import "github.com/Popov-Dmitriy-Ivanovich/genmilk_backend/models"

type LactationLoader struct {
	*ModelLoader
}

func (ll *LactationLoader) Init(mc chan any) {
	ll.LoaderFunc = LoadModelToDb[models.Lactation]()
	ll.ModelsChannel = mc
}
