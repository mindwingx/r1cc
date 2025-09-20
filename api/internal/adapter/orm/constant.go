package orm

const (
	// DbLockUpdate : no DB read and write are allowed by others while updating
	DbLockUpdate string = "UPDATE"
	// DbLockShare : DB read is allowed, but no write is allowed by others while updating
	DbLockShare string = "SHARE"
)
