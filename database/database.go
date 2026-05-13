package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type History struct {
	Output  string // Refers to the output of a task
	Content string // Refers to the input of a task
}

var DB *sql.DB

func Connect(connStr string) error {
	var err error
	DB, err = sql.Open("postgres", connStr)
	return err
}

func SaveMessage(messageID, taskID, content, user string) error {
	_, err := DB.Exec(
		"INSERT INTO spotnik.messages (id, task_id, content, \"user\") VALUES ($1, $2, $3, $4)",
		messageID, taskID, content, user,
	)
	return err
}

func UpdateMessage(msgID, output string) error {
	_, err := DB.Exec(
		"UPDATE spotnik.messages SET output = $1 WHERE id = $2",
		output, msgID,
	)
	return err
}

func LoadHistory(user string) ([]History, error) {
	rows, err := DB.Query(
		"SELECT content, output FROM spotnik.messages WHERE \"user\" = $1 ORDER BY created_at ASC",
		user,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var his []History
	for rows.Next() {
		var h History
		rows.Scan(&h.Content, &h.Output)
		his = append(his, h)
	}
	return his, nil
}

func CreateTask(taskID string) error {
	_, err := DB.Exec(
		"INSERT INTO spotnik.tasks (id, status) VALUES ($1, $2)",
		taskID, "created",
	)
	return err
}

func UpdateTaskStatus(taskID, status string) error {
	_, err := DB.Exec(
		"UPDATE spotnik.tasks SET status = $1 WHERE id = $2",
		status, taskID,
	)
	return err
}
