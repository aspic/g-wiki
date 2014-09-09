/*
GNU GPLv3 - see LICENSE
*/

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

// GitAdd node
func (node *Node) GitAdd() *Node {
	gitCmd(exec.Command("git", "add", node.File))
	return node
}

// GitCommit node message.
func (node *Node) GitCommit(msg string, author string) *Node {
	if author != "" {
		gitCmd(exec.Command("git", "commit", "-m", msg,
			fmt.Sprintf("--author='%s <system@g-wiki>'", author)))
	} else {
		gitCmd(exec.Command("git", "commit", "-m", msg))
	}
	return node
}

// GitShow fetches the node revision.
func (node *Node) GitShow() *Node {
	buf := gitCmd(exec.Command("git", "show", node.Revision+":"+node.File))
	node.Bytes = buf.Bytes()
	return node
}

// GitLog fetches the node log.
func (node *Node) GitLog() *Node {
	buf := gitCmd(exec.Command(
		"git", "log", "--pretty=format:%h %ad %s", "--date=relative",
		"-n", logLimitS, "--", node.File))
	var err error
	b := bufio.NewReader(buf)
	var bytes []byte
	node.Log = make([]*Log, 0)
	for err == nil {
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
	for i, dir := range strings.Split(path, "/") {
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

// GitRevert soft resets to the node's specific revision.
func (node *Node) GitRevert() *Node {
	log.Printf("Reverts %v to revision %s", node, node.Revision)
	gitCmd(exec.Command("git", "checkout", node.Revision, "--", node.File))
	return node
}

// Run git command, will currently die on all errors
func gitCmd(cmd *exec.Cmd) *bytes.Buffer {
	cmd.Dir = fmt.Sprintf("%s/", directory)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		log.Printf("Error: command %q failed (%v) with: %s",
			strings.Join(cmd.Args, " "), err, errBuf.String())
		return &bytes.Buffer{}
	}
	return &outBuf
}


