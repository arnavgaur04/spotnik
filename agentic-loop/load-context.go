package agenticloop

import (
	"fmt"
	"spotnik/llm"

	"google.golang.org/genai"
)

func CallGemini(contents []*genai.Content) (string, error) {
	reply, err := llm.CallGemini(contents)
	if err != nil {
		fmt.Println("LLM error:", err)
		return "", err
	}

	return reply, nil
}
