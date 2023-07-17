package helpers

import "golang.org/x/crypto/bcrypt"

func ValidatePassword(hashedPassword string, plainPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)); err != nil {
		return err
	}

	return nil
}
