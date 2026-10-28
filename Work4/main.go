package main

import "fmt"

type Profile struct {
	Name       string
	Age        int
	Email      string
	Newsletter bool
}

type ProfilePatch struct {
	Name       *string
	Age        *int
	Email      *string
	Newsletter *bool
}

// ApplyPatch применяет частичное обновление к профилю.
func ApplyPatch(p *Profile, patch ProfilePatch) {
	// Защита от паники, если передан nil вместо профиля
	if p == nil {
		return
	}

	// Проверяем каждое поле: если указатель не nil, значит значение нужно обновить
	if patch.Name != nil {
		p.Name = *patch.Name
	}
	if patch.Age != nil {
		p.Age = *patch.Age
	}
	if patch.Email != nil {
		p.Email = *patch.Email
	}
	if patch.Newsletter != nil {
		p.Newsletter = *patch.Newsletter
	}
}

func main() {
	// Исходный профиль
	profile := Profile{
		Name:       "Анна",
		Age:        30,
		Email:      "anna@example.com",
		Newsletter: true,
	}

	// Создаем значения для патча
	newAge := 31
	newNewsletter := false

	// Формируем патч (Name и Email не трогаем, оставляем nil)
	patch := ProfilePatch{
		Age:        &newAge,
		Newsletter: &newNewsletter,
	}

	// Применяем патч
	ApplyPatch(&profile, patch)
	fmt.Printf("После патча: %+v\n", profile)
	// Вывод: После патча: {Name:Анна Age:31 Email:anna@example.com Newsletter:false}

	// Проверка пустого патча
	ApplyPatch(&profile, ProfilePatch{})
	fmt.Printf("После пустого патча: %+v\n", profile)
	// Вывод не изменится

	// Проверка на nil профиль (не должно быть паники)
	ApplyPatch(nil, patch)
	fmt.Println("Вызов ApplyPatch(nil, ...) не вызвал панику")
}