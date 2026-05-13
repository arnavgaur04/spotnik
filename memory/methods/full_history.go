package methods

import (
	"fmt"
	"spotnik/database"
)

func FullHistory(user string) ([]database.History, error) {
	rows, err := database.DB.Query(
		"SELECT content, output FROM spotnik.messages WHERE \"user\" = $1 ORDER BY created_at ASC",
		user,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var his []database.History
	for rows.Next() {
		var h database.History
		rows.Scan(&h.Content, &h.Output)
		fmt.Printf("  -> Content (User): %s\n", h.Content)
		fmt.Printf("  -> Output  (AI):   %s\n", h.Output)
		fmt.Printf("---------------------------------------\n")
		his = append(his, h)
	}
	return his, nil
}
