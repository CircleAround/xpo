package apikit

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type gparams struct {
	vars map[string]string
}

func URLParams(r *http.Request) Params {
	p := new(gparams)
	p.vars = mux.Vars(r)
	return p
}

func (p *gparams) Get(key string) string {
	return p.vars[key]
}

func (p *gparams) AsInt(key string) (int, error) {
	ret, err := strconv.ParseInt(p.vars[key], 0, 0)
	return int(ret), err
}

func (p *gparams) AsInt64(key string) (int64, error) {
	return strconv.ParseInt(p.vars[key], 0, 64)
}
