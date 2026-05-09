package main

import (
    "context"
    "bufio"
    "os"
    "fmt"
    "log"
    "google.golang.org/genai"
)

func main() {
    ctx := context.Background()
    // The client gets the API key from the environment variable `GEMINI_API_KEY`.
    client, err := genai.NewClient(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Enter prompt below:")
    reader := bufio.NewReader(os.Stdin)
    prompt, _ := reader.ReadString('\n')

    fmt.Printf("Calling Gemini with prompt %s\n", prompt)

    result, err := client.Models.GenerateContent(
        ctx,
        "gemini-3-flash-preview",
        genai.Text(prompt),
        nil,
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Result:")
    fmt.Println(result.Text())
}
