package context

type CtxKey string

const UserIDKey = CtxKey("user_id")
const DeviceIDKey = CtxKey("device_id")

type UserContext struct {
	ID       int
	DeviceID int
}

func NewUserCtx(id, device int) *UserContext {
	return &UserContext{id, device}
}
