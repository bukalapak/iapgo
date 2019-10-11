package iapgo_test

import (
	"net/http"

	"github.com/saifulwebid/iapgo"
)

var _ http.RoundTripper = new(iapgo.Transport)
