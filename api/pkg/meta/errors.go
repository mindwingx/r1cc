package meta

import (
	"microservice/pkg/meta/status"
)

var (
	Validate      = ServiceErr(status.Validate)
	Failed        = ServiceErr(status.Failed)
	NotFound      = ServiceErr(status.NotFound)
	Conflict      = ServiceErr(status.Conflict)
	ItemExist     = ServiceErr(status.ItemExist)
	Unauthorized  = ServiceErr(status.Unauthorized)
	Forbidden     = ServiceErr(status.Forbidden)
	InvalidClient = ServiceErr(status.InvalidClient)
	TokenExpired  = ServiceErr(status.TokenExpired)
	DtoBindErr    = ServiceErr(status.DtoBindErr)
)
