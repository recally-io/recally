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

const chatMessageWithRagPrompt = `
I have searched our knowledge base and chat history for relevant information based on the user's query. Please use this information to construct a well-informed response.

Knowledge from our database:
{{.KnowledgeBaseInfo}}

Relevant chat history:
{{.ChatHistory}}

User's query:
{{.UserQuery}}

Please provide a comprehensive and accurate response to the user's query, incorporating the provided knowledge and relevant chat history where appropriate. If the information provided is not sufficient to answer the query fully, please state so and offer to assist with finding more information.
`

var chatMessageWithRagPromptTmpl = template.Must(template.New("promptChatMessageWithRag").Parse(chatMessageWithRagPrompt))

func getChatMessageWithRagPrompt(knowledgeBaseInfo, chatHistory, userQuery string) (string, error) {
	var tpl bytes.Buffer
	if err := chatMessageWithRagPromptTmpl.Execute(&tpl, map[string]interface{}{
		"KnowledgeBaseInfo": knowledgeBaseInfo,
		"ChatHistory":       chatHistory,
		"UserQuery":         userQuery,
	}); err != nil {
		return "", err
	}
	return tpl.String(), nil
}
