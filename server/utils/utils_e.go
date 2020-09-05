package utils

import "errors"

func errorUnexpectedJwtSigningMethod() error {
	return errors.New("unexpected jwt token signing method")
}

func errorInvalidTokenClaims() error {
	return errors.New("invalid token claims")
}
