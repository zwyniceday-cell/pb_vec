package handlers

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

type TestVecHandler struct{}

func (h *TestVecHandler) Register(se *core.ServeEvent) error {
	RegisterTestVecHandler(se.App.(*pocketbase.PocketBase), se)
	return nil
}

func RegisterTestVecHandler(app *pocketbase.PocketBase, se *core.ServeEvent) {
	se.Router.GET("/api/testvec", func(re *core.RequestEvent) error {
		var result struct {
			Version  string  `db:"version"`
			Distance float64 `db:"distance"`
		}

		// Query both the version and a sample distance calculation
		err := app.DB().NewQuery("SELECT vec_version() as version, vec_distance_cosine('[1, 2, 3]', '[4, 5, 6]') as distance").
			One(&result)

		if err != nil {
			return re.JSON(500, map[string]any{
				"error":   "Failed to execute vector query",
				"details": err.Error(),
			})
		}

		return re.JSON(200, map[string]any{
			"status":   "success",
			"version":  result.Version,
			"distance": result.Distance,
			"message":  "sqlite-vec extension is active and working correctly.",
		})
	})
}
