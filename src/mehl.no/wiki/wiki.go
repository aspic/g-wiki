package main

import (
    "html/template"
    "net/http"
    "log"
    "fmt"
    "flag"
    "io/ioutil"
    "os"
    "path"
    "github.com/russross/blackfriday"
    "os/exec"
    "bytes"
    "bufio"
    "encoding/json"
    "strings"
)

const (
    default_host = "localhost"
    default_port = 8080
    dir = "files"
    log_limit = 10
)

type Node struct {
    Path string
    File string
    Content string
    Template string
    Markdown string
    Log []*Log
    Revision string
    Dirs []string
    Active string
}

type Log struct {
    Hash string
    Message string
    Time string
}
// Add node
func (node *Node) GitAdd() *Node {
    gitCmd(exec.Command("git", "add", node.File))
    return node
}
// Commit node message
func (node *Node) GitCommit(msg string) *Node {
    gitCmd(exec.Command("git", "commit", "-m", msg))
    return node
}
// Fetch node revision
func (node *Node) GitShow() *Node {
    buf := gitCmd(exec.Command("git", "show", node.Revision+":"+node.File))
    node.Markdown = string(blackfriday.MarkdownBasic(buf.Bytes()))
    if node.Markdown == "" {
        node.Markdown = " "
    }
    return node
}
// Fetch node log
func (node *Node) GitLog() *Node {
    buf := gitCmd(exec.Command("git", "log", "--pretty=format:{\"Hash\": \"%h\", \"Message\":\"%s\", \"Time\":\"%ad\"}", "--date=relative", node.File))
    var err error
    b := bufio.NewReader(buf)
    var bytes []byte
    node.Log = make([]*Log, 0)
    for (err == nil) {
        if len(node.Log) >= log_limit {
            break
        }
        bytes, err = b.ReadSlice('\n')
        logLine := &Log{}
        err = json.Unmarshal(bytes, logLine)
        if err != nil {
            break
        }
        node.Log = append(node.Log, logLine)
    }
    return node
}
// Run git command, will currently die on all errors
func gitCmd(cmd *exec.Cmd) (*bytes.Buffer) {
    cmd.Dir = fmt.Sprintf("%s/", dir)
    var out bytes.Buffer
    cmd.Stdout = &out
    runError := cmd.Run()
    if runError != nil {
        log.Fatal(fmt.Sprintf("Command failed with:\n\"%s\n\"", out.String()))
    }
    return &out
}

func (node *Node) DirHir(path string) {
    dirs := strings.Split(path, "/")
    curr := ""
    node.Dirs = make([]string, 0)
    for _,line := range dirs {
        curr += line
        node.Dirs = append(node.Dirs, curr)
    }
}

func wikiHandler(w http.ResponseWriter, r *http.Request) {

    if r.URL.Path == "/favicon.ico" {
        return
    }
    // Params
    content := r.FormValue("content")
    edit := r.FormValue("edit")
    changelog := r.FormValue("msg")
    revision := r.FormValue("show")

    filePath := fmt.Sprintf("%s%s.md", dir, r.URL.Path)
    node := &Node{File: r.URL.Path[1:] + ".md", Path: r.URL.Path}

    entry := r.URL.Path

    node.Active = path.Base(entry)
    if len(path.Dir(entry)) > 1 {
        node.Dirs = strings.Split(path.Dir(entry), "/")
    }

    // Write file
    if content != "" && changelog != "" {
        bytes := []byte(content)
        err := writeFile(bytes, filePath)
        if err != nil {
            log.Printf("Cant write to file %s, error: ", filePath, err)
        } else {
            // Written file, commit
            node.GitAdd().GitCommit(changelog).GitLog()
            node.Markdown = string(blackfriday.MarkdownBasic(bytes))
        }
    } else if(revision != "") {
        node.Revision = revision
        node.GitShow().GitLog()
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
            node.GitLog()
        }
        node.Template = "templates/edit.tpl"
    }
    renderTemplate(w, node)
}

func writeFile(bytes []byte, entry string) error {
    err := os.MkdirAll(path.Dir(entry), 0777)
    if err == nil {
        return ioutil.WriteFile(entry, bytes, 0644)
    }
    return err
}

func renderTemplate(w http.ResponseWriter, node *Node) {

    t := template.New("test")
    var err error

    // Build template
    if node.Markdown != "" {
        tpl := fmt.Sprintf("%s\n%s", "{{ template \"header\" . }}", node.Markdown)
       if node.Revision == "" {
            tpl += "{{ template \"actions\" .}}"
        }
        // Footer
        tpl += "{{ template \"footer\" . }}"
        t.Parse(tpl)
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
