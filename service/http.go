//go:generate $GOPATH/bin/mockgen -destination=../mock/mock_http_service.go -package=mock go-app/service HTTP
package service

import "net/http"

type HTTP interface {
	Get(url string) (resp *http.Response, err error)
}

func (h *HTTPImpl) Get(url string) (resp *http.Response, err error) {
	return http.Get(url)
}
