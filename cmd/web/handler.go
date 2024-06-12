package main

import (
	"fmt"
	"net/http"
	"html/template"
)

var Response string

// Renders homepage
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	render(w, files, &TemplateData{
		
	})
}

// Fetches the Degug response from OpenAI and redirects to home
func (app *application) fetchDebugPage(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/response.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	render(w, files, &TemplateData{
		Response: "",
		PageLayoutData: &PageLayoutData{
			Title: "Debugger",
		},
		PromptMessage: "",
	})
}
func (app *application) fetchDebug(w http.ResponseWriter, r *http.Request) {
	prompt := r.FormValue("Message")

	DebuggerSystem := &ChatSystem{
		SystemMessage: debuggerSystemMessage,
	}

	// Call completion API
	response, err := app.OpenAIClient.GetCompletionResponse(prompt, DebuggerSystem)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	files := []string{
		"./ui/html/response.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	Response:=mdToHTML(response)
	// Response = string(response)
	render(w, files, &TemplateData{
		PageLayoutData: &PageLayoutData{
			Title: "Debugger",
		},
		Response: template.HTML(Response),
		PromptMessage: prompt,
	})

	// fmt.Println("Response:", Response)
}

func (app *application) testerPage(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/response.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	render(w, files, &TemplateData{
		Response: "",
		PageLayoutData: &PageLayoutData{
			Title: "Tester",
		},
		PromptMessage: "",
	})
}
func (app *application) tester(w http.ResponseWriter, r *http.Request) {
	prompt := r.FormValue("Message")

	DebuggerSystem := &ChatSystem{
		SystemMessage: testerSystemMessage,
	}

	// Call completion API
	response, err := app.OpenAIClient.GetCompletionResponse(prompt, DebuggerSystem)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	files := []string{
		"./ui/html/response.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	Response := mdToHTML(response)
	render(w, files, &TemplateData{
		PageLayoutData: &PageLayoutData{
			Title: "Tester",
		},
		Response:      template.HTML(Response),
		PromptMessage: prompt,
	})

	// fmt.Println("Response:", Response)
}

func (app *application) formatterPage(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/response.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	render(w, files, &TemplateData{
		Response:"",
		PageLayoutData: &PageLayoutData{
			Title: "Formatter",
		},
		PromptMessage: "",
	})
}
func (app *application) formatter(w http.ResponseWriter, r *http.Request) {
	prompt := r.FormValue("Message")

	DebuggerSystem := &ChatSystem{
		SystemMessage: formatterSystemMessage,
	}

	// Call completion API
	response, err := app.OpenAIClient.GetCompletionResponse(prompt, DebuggerSystem)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	files := []string{
		"./ui/html/response.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	Response := mdToHTML(response)
	render(w, files, &TemplateData{
		PageLayoutData: &PageLayoutData{
			Title: "Formatter",
		},
		Response:      template.HTML(Response),
		PromptMessage: prompt,
	})
}
// func (app *application) handleQuery(w http.ResponseWriter, r *http.Request) {

// }
