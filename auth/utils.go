package auth

import "golang.org/x/crypto/bcrypt"

type Password struct {
	Password string
}

func (p *Password) HashPassword() (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (p *Password) ComparePasswordHash(passwordHash string)(bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(p.Password))

	if err != nil {
		return false, err
	}
	return true, nil

}

