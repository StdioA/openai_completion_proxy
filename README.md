# completion-proxy
A proxy to translate [OpenAI completion API](https://platform.openai.com/docs/api-reference/completions) (`/v1/completions`) to [OpenAI chat completion API](https://platform.openai.com/docs/api-reference/chat) (`/v1/chat/completions`).

## Usage
Installation: `go install github.com/stdioa/openai_completion_proxy`.

Usage:
```
Usage of ./openai_completion_proxy:
  -endpoint string
        API Endpoint Base (default "https://api.openai.com/v1")
  -listen string
        Listen address (default ":8080")
```
