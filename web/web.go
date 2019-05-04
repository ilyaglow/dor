package dorweb

import (
	"encoding/json"
	"net/http"

	"github.com/ilyaglow/dor"
	"github.com/julienschmidt/httprouter"
)

// Serve represents a web interaction with the DomainRank
func Serve(address string, d *dor.App) error {
	router := httprouter.New()
	router.GET("/rank/:domain", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		result, err := d.Find(ps.ByName("domain"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		enc := json.NewEncoder(w)
		if err := enc.Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	return http.ListenAndServe(address, router)
}
