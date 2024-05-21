package response

const (
	Ok                    = 0o01
	CantDecodeRequestBody = 0o02
	InvalidPayload        = 0o03
	EntityAlreadyExists   = 0o04
	InternalError         = 0o05
	NotFound              = 0o06
)
