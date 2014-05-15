package main

import (
    "html/template"
    "net/http"
    "log"
    "fmt"
    "flag"
    "io/ioutil"
    "os"
    "regexp"
    "path"
    "github.com/russross/blackfriday"
    "os/exec"
    "bytes"
    "bufio"
    "strings"
    "strconv"
)

const (
    DIRECTORY = "files"
    LOG_LIMIT = "5"
)

type Node struct {
    Path string
    File string
    Content string
    Template string
    Revision string
    Bytes []byte
    Dirs []*Directory
    Log []*Log
    Markdown template.HTML

    Revisions bool // Show revisions
}

type Directory struct {
    Path string
    Name string
    Active bool
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
func (node *Node) GitCommit(msg string, author string) *Node {
    if author != "" {
        gitCmd(exec.Command("git", "commit", "-m", msg, fmt.Sprintf("--author='%s <system@g-wiki>'", author)))
    } else {
        gitCmd(exec.Command("git", "commit", "-m", msg))
    }
    return node
}
// Fetch node revision
func (node *Node) GitShow() *Node {
    buf := gitCmd(exec.Command("git", "show", node.Revision+":"+node.File))
    node.Bytes = buf.Bytes()
    return node
}
// Fetch node log
func (node *Node) GitLog() *Node {
    buf := gitCmd(exec.Command("git", "log", "--pretty=format:%h %ad %s", "--date=relative", "-n", LOG_LIMIT, node.File))
    var err error
    b := bufio.NewReader(buf)
    var bytes []byte
    node.Log = make([]*Log, 0)
    for (err == nil) {
        bytes, err = b.ReadSlice('\n')
        logLine := parseLog(bytes)
        if logLine == nil {
            continue
        } else if logLine.Hash != node.Revision {
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

func parseLog(bytes []byte) *Log {
    line := string(bytes)
    re := regexp.MustCompile(`(.{0,7}) (\d+ \w+ ago) (.*)`)
    matches := re.FindStringSubmatch(line)
    if len(matches) == 4 {
        return &Log{Hash: matches[1], Time: matches[2], Message: matches[3]}
    }
    return nil
}

func listDirectories(path string) []*Directory {
    s := make([]*Directory, 0)
    dirPath := ""
    for i,dir := range (strings.Split(path, "/")) {
        if i == 0 {
            dirPath += dir
        } else {
            dirPath += "/" + dir
        }
        s = append(s, &Directory{Path: dirPath, Name: dir})
    }
    if len(s) > 0 {
        s[len(s)-1].Active = true
    }
    return s
}

// Soft reset to specific revision
func (node *Node) GitRevert() *Node {
    log.Printf("Reverts %s to revision %s", node, node.Revision)
    gitCmd(exec.Command("git", "checkout", node.Revision, "--", node.File))
    return node
}
// Run git command, will currently die on all errors
func gitCmd(cmd *exec.Cmd) (*bytes.Buffer) {
    cmd.Dir = fmt.Sprintf("%s/", DIRECTORY)
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
    node.Markdown = template.HTML(string(blackfriday.MarkdownCommon(node.Bytes)))
}

func ParseBool(value string) bool {
    boolValue, err := strconv.ParseBool(value)
    if err != nil {
        return false
    }
    return boolValue
}

func wikiHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/favicon.ico" {
        return
    }
    // Params
    content := r.FormValue("content")
    edit := r.FormValue("edit")
    changelog := r.FormValue("msg")
    author := r.FormValue("author")
    reset := r.FormValue("revert")
    revision := r.FormValue("revision")

    filePath := fmt.Sprintf("%s%s.md", DIRECTORY, r.URL.Path)
    node := &Node{File: r.URL.Path[1:] + ".md", Path: r.URL.Path}
    node.Revisions = ParseBool(r.FormValue("revisions"))

    node.Dirs = listDirectories(r.URL.Path)

    // We have content, update
    if content != "" && changelog != "" {
        bytes := []byte(content)
        err := writeFile(bytes, filePath)
        if err != nil {
            log.Printf("Cant write to file %s, error: ", filePath, err)
        } else {
            // Wrote file, commit
            node.Bytes = bytes
            node.GitAdd().GitCommit(changelog, author).GitLog()
            node.ToMarkdown()
        }
    } else if reset != "" {
        // Reset to revision
        node.Revision = reset
        node.GitRevert().GitCommit("Reverted to: " + node.Revision, author)
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

    t := template.New("wiki")
    var err error

    // Build template
    if node.Markdown != "" {
        tpl := "{{ template \"header\" . }}"
        if node.isHead() {
            tpl += "{{ template \"actions\" .}}"
        } else if node.Revision != "" {
            tpl += "{{ template \"revision\" . }}"
        }
        // Add node
        tpl += "{{ template \"node\" . }}"
        // Show revisions
        if node.Revisions {
            tpl += "{{ template \"revisions\" . }}"
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
    t.ParseFiles("templates/header.tpl", "templates/footer.tpl",
    "templates/actions.tpl", "templates/revision.tpl",
    "templates/revisions.tpl", "templates/node.tpl")
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

    var local = flag.String("local", "", "serve as webserver, example: 0.0.0.0:8000")

    flag.Parse()
    var err error

    if *local != "" {
        err = http.ListenAndServe(*local, nil)
    }
    if err != nil {
        panic("ListenAndServe: " + err.Error())
    }
}
