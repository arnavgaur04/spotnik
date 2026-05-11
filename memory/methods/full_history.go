package methods

import (
    "spotnik/database"
)

func FullHistory(user string) ([] database.History, error) {
    rows, err := database.DB.Query(
        "SELECT content, output FROM spotnik.messages WHERE \"user\" = $1 ORDER BY created_at ASC",
        user,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var his [] database.History
    for rows.Next() {
        var h database.History
        rows.Scan(&h.Content, &h.Output)
        his = append(his, h)
    }
    return his, nil
}
