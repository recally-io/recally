package assistants

import (
	"bytes"
	"text/template"
)

const titleGenerationPrompt = `Create a concise, 3-5 word title with an emoji as a title for the prompt in the given language. \
Suitable Emojis for the summary can be used to enhance understanding but avoid quotation marks or special formatting. RESPOND ONLY WITH THE TITLE TEXT.

Examples of titles:
📉 Stock Market Trends
🍪 Perfect Chocolate Chip Recipe
Evolution of Music Streaming
Remote Work Productivity Tips
Artificial Intelligence in Healthcare
🎮 Video Game Development Insights


<prompt>
{{ .Prompt }}
</prompt>

Title: 
`

var titleGenerationPromptTmpl = template.Must(template.New("promptTitleGeneration").Parse(titleGenerationPrompt))

func getTitleGenerationPrompt(conversation string) (string, error) {
	var tpl bytes.Buffer
	if err := titleGenerationPromptTmpl.Execute(&tpl, map[string]interface{}{
		"Prompt": conversation,
	}); err != nil {
		return "", err
	}
	return tpl.String(), nil
}
