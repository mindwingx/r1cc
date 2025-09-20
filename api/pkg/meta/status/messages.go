package status

type HttpMappedStatus string

// general

const (
	Success      HttpMappedStatus = "resp_done"
	Created      HttpMappedStatus = "create_done"
	Updated      HttpMappedStatus = "update_done"
	Deleted      HttpMappedStatus = "delete_done"
	Validate     HttpMappedStatus = "validation_err"
	Failed       HttpMappedStatus = "resp_fail"
	NotFound     HttpMappedStatus = "not_found"
	Conflict     HttpMappedStatus = "conflict_detected"
	Assigned     HttpMappedStatus = "items_assigned"
	Revoked      HttpMappedStatus = "tokens_revoked"
	ItemExist    HttpMappedStatus = "item_exist"
	Unauthorized HttpMappedStatus = "unauthorized"
	Forbidden    HttpMappedStatus = "forbidden"
)

// custom

const (
	InvalidClient HttpMappedStatus = "invalid_client"
	TokenExpired  HttpMappedStatus = "access_token_exp"
	DtoBindErr    HttpMappedStatus = "dto_bind_err"
)
