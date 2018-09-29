package exchi

import (
	"net/http"
	"strconv"

	"local/apikit"

	"github.com/go-chi/chi"
)

type cparams struct {
	r *http.Request
}

func URLParams(r *http.Request) apikit.Params {
	p := new(cparams)
	p.r = r
	return p
}

func (p *cparams) Get(key string) string {
	return chi.URLParam(p.r, key)
}

func (p *cparams) AsInt(key string) (int, error) {
	ret, err := strconv.ParseInt(p.Get(key), 0, 0)
	return int(ret), err
}

func (p *cparams) AsInt64(key string) (int64, error) {
	return strconv.ParseInt(p.Get(key), 0, 64)
}
