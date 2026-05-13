package memory

import (
	"spotnik/database"
	"spotnik/memory/methods"
)

func LoadContext(option int, user, prompt string) ([]database.History, error) {
	switch option {
	case 0:
		return methods.FullHistory(user)
	case 1:
		return methods.RAGHistory(user, prompt, 3)
	default:
		return methods.FullHistory(user)
	}
}
