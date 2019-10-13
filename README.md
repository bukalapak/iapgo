# iapgo

[![Build Status](https://travis-ci.com/saifulwebid/iapgo.svg?branch=master)](https://travis-ci.com/saifulwebid/iapgo)
[![codecov](https://codecov.io/gh/saifulwebid/iapgo/branch/master/graph/badge.svg)](https://codecov.io/gh/saifulwebid/iapgo)

iapgo is a Go library to help authenticating access to endpoints behind Google [Cloud Identity-Aware Proxy](https://cloud.google.com/iap/).

This library is heavily using [`golang.org/x/oauth2/google`](https://godoc.org/golang.org/x/oauth2/google) to handle credentials parsing and [authentication](https://cloud.google.com/iap/docs/authentication-howto).
