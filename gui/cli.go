package gui

import (
	"fmt"

	"github.com/c-bata/go-prompt"
)

// Context main object controlling the cli prompt
type Context struct {
	p *prompt.Prompt
}

// Initialize is the function to initialize the Context for using the prompt
func (c *Context) Initialize() {
	fmt.Println("Please select table.")
	t := prompt.Input("> ", completer)
	fmt.Println("You selected " + t)
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
