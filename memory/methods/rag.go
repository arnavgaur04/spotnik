package methods

import (
    "context"
    "fmt"

    "github.com/pgvector/pgvector-go"
    "google.golang.org/genai"
    "spotnik/database"
)

func GetEmbedding(text string) ([]float32, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("embedding client error: %w", err)
	}

	dim := int32(768)

	result, err := client.Models.EmbedContent(ctx,
		"gemini-embedding-001",
		[]*genai.Content{
    			{
			    Parts: []*genai.Part{{Text: text}},
			},
		},
		&genai.EmbedContentConfig{
        	    OutputDimensionality: &dim,
    		},
	)
	if err != nil {
		return nil, fmt.Errorf("embed content error: %w", err)
	}

	return result.Embeddings[0].Values, nil
}

func RAGHistory(user, currentMsg string, limit int) ([]database.History, error) {
	// Get embedding for current message
	embedding, err := GetEmbedding(currentMsg)
	if err != nil {
		return nil, fmt.Errorf("embedding error: %w", err)
	}

	rows, err := database.DB.Query(
		`SELECT content, output
		 FROM spotnik.messages
		 WHERE "user" = $1
		 AND output IS NOT NULL
		 AND output != ''
		 AND embedding IS NOT NULL
		 ORDER BY embedding <=> $2::vector
		 LIMIT $3`,
		user, pgvector.NewVector(embedding), limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []database.History
	for rows.Next() {
		var his database.History
		if err := rows.Scan(&his.Content, &his.Output); err != nil {
			return nil, err
		}
		history = append(history, his)
	}

	return history, rows.Err()
}

func SaveEmbedding(messageID string, embedding []float32) error {
	_, err := database.DB.Exec(
		`UPDATE spotnik.messages SET embedding = $1 WHERE id = $2`,
		pgvector.NewVector(embedding), messageID,
	)
	return err
}
