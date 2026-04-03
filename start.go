package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Пароль
	password := "admin123"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Ошибка генерации хэша:", err)
		return
	}

	fmt.Println("Пароль:", password)
	fmt.Println("Хэш:", string(hash))

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		fmt.Println("Пароль не совпадает!")
		return
	}
	fmt.Println("Пароль проверен успешно!")
}
