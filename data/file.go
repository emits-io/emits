package data

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/emits-io/emits/configuration"
)

const (
	// appending constant referenced by process function
	appending = separator + separator + separator
	// collapsing constant referenced by process function
	collapsing = ":"
	// config constant referenced by process function
	config = emits + separator
	// fileExtension constant referenced by the Write function
	fileExtension = ".json"
	// emits constant referenced by the configuration constant
	emits = "emits"
	// grammarFilePrefix constant referenced by the process function
	grammarFilePrefix = emits + separator + "grammar" + separator
	// separator constant referenced by process function and configuration constant
	separator = "."
	// escape character referenced by the process function
	escape = "\\"
	// flag
	flag = "`"
	// flagSeparator
	flagSeparator = ","
	// indent character referenced by the process function
	indent = ">"
	// outdent character referenced by the process function
	outdent = "<"
)

// emit structure is used to write the json file format.
type emit struct {
	File          file   `json:"file,omitempty"`
	Configuration []Node `json:"configuration,omitempty"`
	Data          []Node `json:"data,omitempty"`
}

// file structure components available within the emits structure.
type file struct {
	Path      string `json:"path,omitempty"`
	Name      string `json:"name,omitempty"`
	Extension string `json:"extension,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// FileComment structure
type FileComment struct {
	BlockOpen  string
	BlockLine  string
	BlockClose string
	Inline     string
}

// write the emit file format to persistant storage.
func (e emit) write(file string) (err error) {
	data, err := json.MarshalIndent(e, "", "\t")
	if err == nil {
		err = os.MkdirAll(filepath.Dir(file), os.ModePerm)
		if err == nil {
			err = ioutil.WriteFile(file, data, 0644)
		}
	}
	return err
}

func cleanSpace(line string) (index int, clean string) {
	for _, c := range line {
		if unicode.IsSpace(c) {
			index++
		} else {
			break
		}
	}
	return index, strings.TrimSpace(line)
}

func cleanPrefix(line string, prefix string, optional bool) (hasPrefix bool, clean string) {
	if strings.HasPrefix(line, prefix) {
		line = line[len(prefix):]
		index := 0
		for _, c := range line {
			if unicode.IsSpace(c) {
				index++
			} else {
				break
			}
		}
		return true, line[index:]
	}
	return optional, line
}
func cleanSuffix(line string, suffix string, optional bool) (hasPrefix bool, clean string) {
	if strings.HasSuffix(line, suffix) {
		line = line[:len(line)-len(suffix)]
		index := 0
		for _, c := range line {
			if unicode.IsSpace(c) {
				index++
			} else {
				break
			}
		}
		return true, line[index:]
	}
	return optional, line
}

func keywordValueFlagIndex(line string, index int) (keyword string, value string, flags []string, indexDelta int) {
	indexDelta = 0
	split := strings.SplitN(line, separator, 2)
	if len(split) == 2 {
		split = strings.SplitN(split[1], " ", 2)
		if len(split) == 2 {
			keyword = strings.TrimSpace(split[0])
			value = strings.TrimSpace(split[1])
		} else {
			keyword = strings.TrimSpace(split[0])
		}
		//
		keywordMeta := ""
		//
		keywordOverride := ""
		for i, c := range keyword {
			valid := unicode.IsLetter(c) || unicode.IsDigit(c) || string(c) == separator
			if index == 0 && strings.HasPrefix(keyword, emits) && i > 0 && !valid {
				valid = string(c) == separator
			}
			if valid {
				keywordOverride += string(c)
			} else {
				break
			}
		}
		keywordMeta = keyword[len(keywordOverride):]
		keyword = keywordOverride
		//
		if strings.HasPrefix(keywordMeta, outdent) {
			for i, c := range keywordMeta {
				if string(c) == outdent {
					index = index - 1
				} else {
					keywordMeta = keywordMeta[i:]
					break
				}
			}
		}
		if strings.HasPrefix(keywordMeta, indent) {
			for i, c := range keywordMeta {
				if string(c) == indent {
					index = index + 1
				} else {
					keywordMeta = keywordMeta[i:]
					break
				}
			}
		}
		//
		flagsOverride := ""
		if strings.HasPrefix(keywordMeta, flag) && strings.HasSuffix(keywordMeta, flag) {
			for _, c := range keywordMeta[1 : len(keywordMeta)-1] {
				if unicode.IsLetter(c) || unicode.IsDigit(c) || string(c) == flagSeparator {
					flagsOverride += string(c)
				}
			}
		}
		flagsArray := strings.Split(flagsOverride, flagSeparator)
		for _, f := range flagsArray {
			if len(f) > 0 {
				flags = append(flags, f)
			}
		}
	}
	return keyword, value, flags, index
}

// process returns a node structure based on simple string conditions.
func process(name string, line string, lineNumber int, task configuration.Task, cache Cache, previous Node) Node {
	// Options
	isAppending, isCollapsing, isNewline, isConfiguration, isSeparator, isCommentInline, isCommentBlockOpen, isCommentBlockClose, isCommentBlockLine := false, false, false, false, false, false, false, false, false
	keyword, value := "", ""
	var flags []string
	index := 0
	// Clean Up
	index, line = cleanSpace(line)

	// Comments
	if len(strings.TrimSpace(task.Comment.Inline)) > 0 {
		isCommentInline, line = cleanPrefix(line, task.Comment.Inline, false)
	}
	if len(strings.TrimSpace(task.Comment.Block.Open)) > 0 {
		isCommentBlockOpen, line = cleanPrefix(line, task.Comment.Block.Open, false)
	}
	if len(strings.TrimSpace(task.Comment.Block.Close)) > 0 {
		isCommentBlockClose, line = cleanSuffix(line, task.Comment.Block.Close, false)
	}
	if isCommentBlockOpen && isCommentBlockClose {
		isCommentBlockOpen = false
		isCommentBlockLine = false
		isCommentBlockClose = false
		isCommentInline = true
	}
	if previous.Comment.BlockOpen || previous.Comment.BlockLine && !previous.Comment.BlockClose {
		if len(strings.TrimSpace(task.Comment.Block.Line)) > 0 {
			isCommentBlockLine, line = cleanPrefix(line, task.Comment.Block.Line, true)
		}
	}
	if isCommentBlockOpen || isCommentBlockLine || isCommentBlockClose || isCommentInline {
		keyword, value, flags, index = keywordValueFlagIndex(line, index)
		if index == 0 && strings.HasPrefix(keyword, config) {
			// Configuration
			keyword = keyword[len(config):]
			isConfiguration = true
		} else if strings.HasPrefix(keyword, separator) {
			// Separator (Syntax)
			keyword = keyword[1:] // remove the separator character
			isSeparator = true
		} else if strings.HasPrefix(line, escape) {
			// Escape
			index++      // index must be greater to create a child node
			keyword = "" // a keyword is not indended; clear it.
			value = line[len(escape):]
		} else if strings.HasPrefix(value, appending) {
			if strings.HasSuffix(value, collapsing) {
				isCollapsing = true
				if strings.HasSuffix(value, collapsing+collapsing) {
					isNewline = true
				}
			}
			// Appending Initialized
			value = "" // a value is not indended; clear it.
			isAppending = true
		} else if len(keyword) == 0 && len(value) == 0 && isCommentBlockLine {
			// Appending Inline
			index++ // index must be greater to create a child node
			value = line
			isCommentInline = true
		}
	}
	return processGrammar(Node{
		Line:          lineNumber,
		Index:         index,
		Keyword:       keyword,
		Value:         value,
		Flags:         flags,
		Separator:     isSeparator,
		Configuration: isConfiguration,
		Appending:     isAppending,
		Collapsing:    isCollapsing,
		Newline:       isNewline,
		Comment: Comment{
			BlockOpen:  isCommentBlockOpen,
			BlockLine:  isCommentBlockLine,
			BlockClose: isCommentBlockClose,
			Inline:     isCommentInline,
		},
	}, task, cache, filepath.Ext(filepath.Base(name)))
}

// Parse returns a node tree and configuration node array.
func Parse(name string, task configuration.Task, cache Cache) (tree Node, config []Node, err error) {
	file, err := os.Open(name)
	if err != nil {
		return tree, config, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	line := 0
	previousNode := Node{}
	for scanner.Scan() {
		line++
		text := scanner.Text()
		node := process(name, text, line, task, cache, previousNode)
		previousNode = node
		if node.IsComment() {
			// Data
			if node.HasData() && !node.IsConfiguration() {
				if !tree.HasChildren() {
					node.ParentNode = &tree
					tree.AppendChild(node)
				} else {
					lastNode := tree.LastNode()
					if node.Index == lastNode.Index {
						node.ParentNode = lastNode.ParentNode
						node.Parent = lastNode.ParentNode.Line
						lastNode.ParentNode.AppendChild(node)
					} else if node.Index > lastNode.Index {
						node.ParentNode = lastNode
						node.Parent = lastNode.Line
						lastNode.AppendChild(node)
					} else if node.Index < lastNode.Index {
						previousNode := lastNode.PreviousIndexNode(node.Index)
						if previousNode.ParentNode != nil {
							node.ParentNode = previousNode.ParentNode
							node.Parent = previousNode.ParentNode.Line
							previousNode.ParentNode.AppendChild(node)
						} else {
							node.ParentNode = previousNode
							node.Parent = previousNode.Line
							previousNode.AppendChild(node)
						}
					}
				}
			}
			// Config
			if node.IsConfiguration() {
				config = append(config, node)
			}
		} else {
			// Explicit flag required to expose source code; default's to false.
			value := ""
			if task.Source {
				value = strings.TrimSpace(text)
			}
			appendNode := tree.LastNode().LastAppendingNodeOrRoot()
			if appendNode.IsAppending() {
				appendNode.AppendChild(Node{
					Index:      appendNode.Index + 1,
					ParentNode: appendNode,
					Parent:     appendNode.Line,
					Line:       line,
					Value:      value,
					Comment:    Comment{Inline: true},
				})
			}
		}
	}
	return tree, config, scanner.Err()
}

// Write the emits json file to an optional prefix directory.
func Write(name string, task configuration.Task, cache Cache, prefixDirectory ...string) (err error) {

	nodes, configurations, err := Parse(name, task, cache)

	if len(task.Keyword.Include) > 0 && !nodes.HasInstanceOfKeyword(task.Keyword.Include) {
		err = fmt.Errorf("keyword include not found")
	}

	if len(task.Keyword.Exclude) > 0 && nodes.HasInstanceOfKeyword(task.Keyword.Exclude) {
		err = fmt.Errorf("keyword exclude found")
	}

	for _, n := range configurations {
		if len(task.Configuration.Include) > 0 {
			found := false
			for _, k := range task.Configuration.Include {
				if n.Keyword == k {
					found = true
					break
				}
			}
			if !found {
				err = fmt.Errorf("configuration include not found")
			}
		}
		if len(task.Configuration.Exclude) > 0 {
			found := false
			for _, k := range task.Configuration.Exclude {
				if n.Keyword == k {
					found = true
					break
				}
			}
			if found {
				err = fmt.Errorf("configuration exclude found")
			}
		}
	}

	if err == nil {

		nodes.CollapseAppending()

		file := emit{
			File: file{
				Path:      filepath.Dir(name),
				Name:      strings.TrimSuffix(filepath.Base(name), filepath.Ext(filepath.Base(name))),
				Extension: strings.TrimPrefix(filepath.Ext(name), "."),
				Timestamp: time.Now().UTC().String(),
			},
			Configuration: configurations,
			Data:          nodes.Children,
		}
		return file.write(filepath.Join(filepath.Join(prefixDirectory...), name+fileExtension))
	}
	return err
}

func processGrammar(n Node, task configuration.Task, cache Cache, extension string) (node Node) {
	if !n.HasKeyword() && n.HasValue() && cache.HasGrammar() {
		return cache.ProcessGrammar(n, extension)
	}
	return n
}
