package v1

import "fmt"

var ErrInvalidRecordID = fmt.Errorf("invalid organisation id")
var ErrRecordNotFound = fmt.Errorf("organisation not found")
var ErrInvalidRequestOptions = fmt.Errorf("invalid request options")
var ErrInvalidUserID = fmt.Errorf("invalid user id")
var ErrInvalidJWTClaims = fmt.Errorf("invalid jwt claims")
