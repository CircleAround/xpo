package xpo

import "local/xpo/web"

//go:generate go-assets-builder --output=assets/reserved_username_list.go -p=assets ../../../assets/reserved_username_list

func init() {
	web.Routes()
}
