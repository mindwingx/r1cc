package status

import "net/http"

var MappedStatuses = map[HttpMappedStatus]int{
	// general
	Success:      http.StatusOK,
	Created:      http.StatusCreated,
	Updated:      http.StatusNoContent,
	Deleted:      http.StatusNoContent,
	Validate:     http.StatusUnprocessableEntity,
	Failed:       http.StatusBadRequest,
	NotFound:     http.StatusNotFound,
	Conflict:     http.StatusConflict,
	Assigned:     http.StatusNoContent,
	Revoked:      http.StatusNoContent,
	ItemExist:    http.StatusConflict,
	Unauthorized: http.StatusUnauthorized,
	Forbidden:    http.StatusForbidden,
	// custom
	InvalidClient: http.StatusUnauthorized,
	TokenExpired:  http.StatusUnauthorized,
	DtoBindErr:    http.StatusBadRequest,
}
