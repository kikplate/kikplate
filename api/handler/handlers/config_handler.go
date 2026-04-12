package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kickplate/api/lib"
)

type ConfigHandler struct {
	env lib.Env
}

func NewConfigHandler(env lib.Env) ConfigHandler {
	return ConfigHandler{env: env}
}

type appConfigResponse struct {
	lib.Customization
	PlateCategories []lib.PlateCategoryConfig `json:"plate_categories"`
}

func (h ConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appConfigResponse{
		Customization:   h.env.Customization,
		PlateCategories: lib.EffectivePlateCategories(h.env),
	})
}
