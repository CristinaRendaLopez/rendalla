package repository

type FakeAuthRepository struct {
	Credentials AuthCredentials
}

func (f *FakeAuthRepository) GetAuthCredentials() (*AuthCredentials, error) {
	return &f.Credentials, nil
}
