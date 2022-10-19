// Copyright (C) 2017 Michał Matczuk
// Use of this source code is governed by an AGPL-style
// license that can be found in the LICENSE file.

package tunnel

import (
	"errors"
	"fmt"
	"io/ioutil"
)

var (
	errClientNotSubscribed    = errors.New("client not subscribed")
	errClientNotConnected     = errors.New("client not connected")
	errClientAlreadyConnected = errors.New("client already connected")

	errUnauthorised = errors.New("unauthorised")
)

func errClientNotSubscribedHtml(p string) error {
	content, err := ioutil.ReadFile("../templates/client_not_subscribed.html")

	if err != nil {
		return errors.New("Client Not subscribed")
	}
	fmt.Println(string(content) + " " + p)

	return errors.New(string(content) + " " + p)
}
