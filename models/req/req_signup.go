package req

type ReqSignUp struct {
	FullName string `validate:"require"` //tags
	Email    string `validate:"require"`
	Password string `validate:"require"`
}
