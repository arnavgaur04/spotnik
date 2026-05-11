package llm

import (
    "fmt"
    "context"
    "log"
    "google.golang.org/genai"
)


func CallGemini(contents []*genai.Content) (string, error) {
    ctx := context.Background()

    client, err := genai.NewClient(ctx, nil) 
    if err != nil {
        log.Fatal(err)
    }

    result, err := client.Models.GenerateContent(
        ctx,
        "gemini-3-flash-preview",
        contents,
        nil,
    )
    if err != nil {
	log.Fatal(err)
	return "", err
    }

    fmt.Println("Result:")
    fmt.Println(result.Text())

    return result.Text(), nil
}
