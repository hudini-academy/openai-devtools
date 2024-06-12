package main
 
import (
    "fmt"
    "html/template"
    "strings"
    "regexp"
 
    // "strings"
 
    //"log"
    "net/http"
)
 
var Response string
 
// isValidEmail checks if the given email address is valid
func isValidEmail(email string) bool {
    // Regular expression for validating email addresses
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    regex := regexp.MustCompile(pattern)
    return regex.MatchString(email)
}
 
// Renders homepage
func (app *application) Base(w http.ResponseWriter, r *http.Request) {
    files := []string{
        "./ui/html/base.page.tmpl",
        "./ui/html/base.layout.tmpl",
 
    }
 
    render(w, files,nil)
}
 
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
    username, ok := app.session.Get(r, "username").(string)
    if !ok {
        username = ""
    }
    files := []string{
        "./ui/html/home.page.tmpl",
        "./ui/html/base.layout.tmpl",
 
    }
	fmt.Println(username)
 
    render(w, files, &TemplateData{
        Response: Response,
        Username: username,
    })
}
 
// Fetches the Degug response from OpenAI and redirects to home
func (app *application) fetchDebugPage(w http.ResponseWriter, r *http.Request) {
    files := []string{
        "./ui/html/response.page.tmpl",
        "./ui/html/base.layout.tmpl",
    }
 
    render(w, files, &TemplateData{
        Response: Response,
        PageLayoutData: &PageLayoutData{
            Title: "debugger",
        },
    })
}
func (app *application) fetchDebug(w http.ResponseWriter, r *http.Request) {
    prompt := r.FormValue("Message")
 
    DebuggerSystem := &ChatSystem{
        SystemMessage: debuggerSystemMessage,
    }
 
    // Call completion API
    response, err := client.CompleteText(prompt, DebuggerSystem)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    files := []string{
        "./ui/html/response.page.tmpl",
        "./ui/html/base.layout.tmpl",
    }
    Response = string(response)
    render(w, files, &TemplateData{
        Response: Response,
    })
 
    // fmt.Println("Response:", Response)
}
 
func (app *application) handleQuery(w http.ResponseWriter, r *http.Request) {
 
}
 
func (app *application) signUpForm(w http.ResponseWriter, r *http.Request) {
    files := []string{
        "./ui/html/signup.page.tmpl",
        "./ui/html/base.layout.tmpl",
    }
    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.errorLog.Println(err.Error())
        http.Error(w, "internal server error", 500)
        return
    }
    ts.Execute(w, app.session.PopString(r, "flash"))
}
 
func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
    userName := r.FormValue("name")
    userEmail := r.FormValue("email")
    userPass := r.FormValue("password")
 
 
    if !isValidEmail(userEmail) {
        app.session.Put(r, "flash", "Invalid email address")
        http.Redirect(w, r, "/user/signup", http.StatusSeeOther)
        return
    }
 
    if len(strings.TrimSpace(userName)) != 0 && len(strings.TrimSpace(userEmail)) !=0 && len(strings.TrimSpace(userPass)) !=0{
        err := app.users.Insert(userName, userEmail, userPass)
        if err != nil {
            app.errorLog.Println(err.Error())
            http.Error(w, "Internal Server Error", 500)
        }
        app.session.Put(r, "flash", "SignUp Successful")
        http.Redirect(w, r, "/user/login", http.StatusSeeOther)
    } else {
        app.session.Put(r, "flash", " item field cant be empty!")
        http.Redirect(w, r, "/user/signup", http.StatusSeeOther)
    }
   
   
}
 
func (app *application) loginForm(w http.ResponseWriter, r *http.Request) {
    files := []string{
        "./ui/html/login.page.tmpl",
        "./ui/html/base.layout.tmpl",
    }
    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.errorLog.Println(err.Error())
        http.Error(w, "internal server error", 500)
        return
    }
    ts.Execute(w, app.session.PopString(r, "flash"))
}
 
func (app *application) login(w http.ResponseWriter, r *http.Request) {
 
    userEmail := r.FormValue("email")
    userPass := r.FormValue("password")
 
    if !isValidEmail(userEmail) {
        app.session.Put(r, "flash", "Invalid email address")
        http.Redirect(w, r, "/user/login", http.StatusSeeOther)
        return
    }
    isUser, err := app.users.Authenticate(userEmail, userPass)
    if err != nil {
        app.errorLog.Println(err.Error())
        http.Error(w, "internal server error", 500)
        return
    }
    if isUser {
        app.session.Put(r, "Authentication", true)
        app.session.Put(r, "flash", "Login Successful")
        app.session.Put(r, "username", isUser)
        http.Redirect(w, r, "/home", http.StatusSeeOther)
    } else {
        app.session.Put(r, "flash", "Login failed")
        http.Redirect(w, r, "/user/login", http.StatusSeeOther)
        app.session.Put(r, "Authentiaction", false)
 
    }
}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
    app.session.Put(r, "Authenticated", false)
    http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}