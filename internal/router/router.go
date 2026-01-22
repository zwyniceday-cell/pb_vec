package router

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

type Registrar interface {
	Register(se *core.ServeEvent) error
}

func RegisterRoutes(app *pocketbase.PocketBase, registrars ...Registrar) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		for _, r := range registrars {
			if err := r.Register(se); err != nil {
				return err
			}
		}
		return se.Next()
	})
}
