package server

import "errors"

func errorConfigIsNil() error {
	return errors.New("error in server: config is nil")
}

func errorPortIsEmpty() error {
	return errors.New("error in server: tcp port is empty")
}
