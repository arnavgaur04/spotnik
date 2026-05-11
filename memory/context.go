package memory

import (
    "fmt"
    "spotnik/database"
    "spotnik/memory/methods"
)

func LoadContext(option int, user, prompt string) ([] database.History, error) {
    fmt.Println("Option Selected Below:")
    fmt.Println(option)

    switch option {
        case 0:
	    return methods.FullHistory(user)
	case 1:
	    return methods.RAGHistory(user, prompt, 3)
        default:
	    return methods.FullHistory(user)
    }
}
