package auth

type TokenVersionValidator interface {
	ValidateTokenVersion(int, int) (bool, error)
}
