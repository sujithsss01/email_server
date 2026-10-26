package emailmodule

type emailHandler struct {
	service *emailService
}

func newEmailHandler(service *emailService) *emailHandler {
	return &emailHandler{
		service: service,
	}
}
