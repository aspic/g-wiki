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
    log_limit = "5"
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
    Bytes []byte
}

type Log struct {
    Hash string
    Message string
    Time string
    Link bool
}
func (node *Node) isHead() bool {
    return len(node.Log) > 0 && node.Revision == node.Log[0].Hash
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
    node.Bytes = buf.Bytes()
    log.Print(len(node.Bytes))
    log.Print(node.Revision)
    return node
}
// Fetch node log
func (node *Node) GitLog() *Node {
    buf := gitCmd(exec.Command("git", "log", "--pretty=format:{\"Hash\": \"%h\", \"Message\":\"%s\", \"Time\":\"%ad\"}", "--date=relative", "-n", log_limit, node.File))
    var err error
    b := bufio.NewReader(buf)
    var bytes []byte
    node.Log = make([]*Log, 0)
    for (err == nil) {
        bytes, err = b.ReadSlice('\n')
        logLine := &Log{}
        err = json.Unmarshal(bytes, logLine)
        if err != nil {
            log.Print(err)
            break
        }
        if logLine.Hash != node.Revision {
            logLine.Link = true
        }
        node.Log = append(node.Log, logLine)
    }
    if node.Revision == "" && len(node.Log) > 0 {
        node.Revision = node.Log[0].Hash
        node.Log[0].Link = false
    }
    return node
}
// Soft reset to specific revision
func (node *Node) GitRevert() *Node {
    log.Printf("Reverts %s to revision %s", node, node.Revision)
    gitCmd(exec.Command("git", "checkout", node.Revision, "--", node.File))
    return node
}
// Run git command, will currently die on all errors
func gitCmd(cmd *exec.Cmd) (*bytes.Buffer) {
    cmd.Dir = fmt.Sprintf("%s/", dir)
    var out bytes.Buffer
    cmd.Stdout = &out
    runError := cmd.Run()
    if runError != nil {
        log.Print(fmt.Sprintf("Error: command failed with:\n\"%s\n\"", out.String()))
        return bytes.NewBuffer([]byte{})
    }
    return &out
}
// Process node contents
func (node *Node) ToMarkdown() {
    node.Markdown = string(blackfriday.MarkdownCommon(node.Bytes))
}

func wikiHandler(w http.ResponseWriter, r *http.Request) {

    if r.URL.Path == "/favicon.ico" {
        return
    }
    log.Printf("ip: %s, route: %s", r.RemoteAddr, r.URL.Path)

    // Params
    content := r.FormValue("content")
    edit := r.FormValue("edit")
    changelog := r.FormValue("msg")
    reset := r.FormValue("revert")
    revision := r.FormValue("revision")

    filePath := fmt.Sprintf("%s%s.md", dir, r.URL.Path)
    node := &Node{File: r.URL.Path[1:] + ".md", Path: r.URL.Path}

    entry := r.URL.Path

    node.Active = path.Base(entry)
    if len(path.Dir(entry)) > 1 {
        node.Dirs = strings.Split(path.Dir(entry), "/")
    }

    // We have content, update
    if content != "" && changelog != "" {
        bytes := []byte(content)
        err := writeFile(bytes, filePath)
        if err != nil {
            log.Printf("Cant write to file %s, error: ", filePath, err)
        } else {
            // Wrote file, commit
            node.Bytes = bytes
            node.GitAdd().GitCommit(changelog).GitLog()
            node.ToMarkdown()
        }
    } else if(reset != "") {
        // Reset to revision
        node.Revision = reset
        node.GitRevert().GitCommit("Reverted to: " + node.Revision)
        node.Revision = ""
        node.GitShow().GitLog()
        node.ToMarkdown()
    } else {
        // Show specific revision
        node.Revision = revision
        node.GitShow().GitLog()
        if edit == "true" || len(node.Bytes) == 0 {
            node.Content = string(node.Bytes)
            node.Template = "templates/edit.tpl"
        } else {
            node.ToMarkdown()
        }
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
        tpl := "{{ template \"header\" . }}"
        if node.isHead() {
            tpl += "{{ template \"actions\" .}}"
        } else if node.Revision != "" {
            tpl += "{{ template \"revision\" . }}"
        }
        // Add content
        tpl += node.Markdown
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
    t.ParseFiles("templates/header.tpl", "templates/footer.tpl",
    "templates/actions.tpl", "templates/revision.tpl")
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
