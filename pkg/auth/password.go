package auth

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmptyPassword    = errors.New("password cannot be empty")
	ErrPasswordTooShort = errors.New("password is too short")
)

// HashPassword хеширует пароль используя bcrypt
func HashPassword(password string) (string, error) {
	// 1. Проверка того, что пароль не пустой
	if password == "" {
		return "", ErrEmptyPassword
	}

	// 2. Использование bcrypt - для хеширования
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// 3. Вернуть хешированный пароль как строку
	return string(hashedPassword), nil
}

// CheckPassword проверяет соответствие пароля и его хеша
func CheckPassword(password, hash string) bool {
	// 1. Сравнение пароля с хешем, используя bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePasswordStrength проверяет надежность пароля
func ValidatePasswordStrength(password string) error {
	// Проверка минимальной длины
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Проверка наличия различных типов символов
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// Проверка наличия хотя бы 3-х из 4-х типов символов
	count := 0
	if hasUpper {
		count++
	}
	if hasLower {
		count++
	}
	if hasNumber {
		count++
	}
	if hasSpecial {
		count++
	}

	if count < 3 {
		return errors.New("password must contain at least 3 different character types (uppercase, lowercase, numbers, special characters)")
	}

	return nil
}
