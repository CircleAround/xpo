package web

import (
	"local/gaekit"
	"local/xpo/app"
	"log"
	"net/http"
)

func NewAdminHandler() *gaekit.AdminHandler {
	ah := gaekit.NewAdminHandler()
	ah.AddMenu("identity_name", func(w http.ResponseWriter, r *http.Request) {
		err := app.NewXUserService().MigrateUniqueIndex(Context(r))
		if err != nil {
			log.Fatalf("マイグレーション中のエラー: %v", err)
		}
	}, "IdentityNameへ移行")
	return ah
}
