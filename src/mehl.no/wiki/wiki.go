package main

import (
    "html/template"
    "net/http"
    "log"
    "fmt"
    "flag"
    "io/ioutil"
    "github.com/russross/blackfriday"
)

const (
    default_host = "localhost"
    default_port = 8080
    dir = "files"
)

type Node struct {
    Path string
    Content string
    Template string
    Markdown string
}

func wikiHandler(w http.ResponseWriter, r *http.Request) {

    if r.URL.Path == "/favicon.ico" {
        return
    }
    content := r.FormValue("content")
    edit := r.FormValue("edit")

    filePath := fmt.Sprintf("%s%s.md", dir, r.URL.Path)

    node := &Node{Path: r.URL.Path}

    // Write file
    if content != "" {
        bytes := []byte(content)
        err := ioutil.WriteFile(filePath, bytes, 0644)
        if err != nil {
            log.Print("Cant write to file", filePath)
        } else {
            node.Markdown = string(blackfriday.MarkdownBasic(bytes))
        }
    } else {
        bytes, err := ioutil.ReadFile(filePath)
        if err != nil {
            log.Printf("No file with path: %s", filePath)
        } else {
            if edit == "true" {
                node.Content = string(bytes)
            } else {
                node.Markdown = string(blackfriday.MarkdownBasic(bytes))
            }
        }
        node.Template = "templates/edit.tpl"
    }
    renderTemplate(w, node)
}

func renderTemplate(w http.ResponseWriter, node *Node) {

    t := template.New("test")
    var err error

    if node.Markdown != "" {
        t.Parse(fmt.Sprintf("%s\n%s\n%s\n%s", "{{ template \"header\" .}}", node.Markdown , "{{ template \"actions\"}}", "{{ template \"footer\" .}}" ))
    } else if node.Template != "" {
        t, err = template.ParseFiles(node.Template)
        if err != nil {
            log.Print("Could not parse template", err)
        }
    }

    // Include the rest
    t.ParseFiles("templates/header.tpl", "templates/footer.tpl", "templates/actions.tpl")
	err = t.Execute(w, node)
    if err != nil {
        log.Print("Could not execute template: ", err)
    }
}

func main() {
    // Handlers
    http.HandleFunc("/", wikiHandler)

    // Static resources
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

    var host = flag.String("h", default_host, "Hostname")
    var port = flag.Int("p", default_port, "Port")
    flag.Parse()

    err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)
    if err != nil {
        panic("ListenAndServe: " + err.Error())
    }
}
