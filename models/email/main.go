package model

type WelcomeAccountMailContent struct {
	Code int
	Name string
}

type ResetPasswordMailContent struct {
	Code int
}

type DeleteAccountMailContent struct {
	Code    int
	Name    string
	Message string
}
