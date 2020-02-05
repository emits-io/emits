package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	name = "emits.json"
)

// File struct
type File struct {
	Groups []Group `json:"group,omitempty"`
	Tasks  []Task  `json:"task,omitempty"`
}

// Index struct
type Index struct {
	Files []string `json:"file"`
}

// Task struct
type Task struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Comment       Comment  `json:"comment"`
	Source        bool     `json:"source"`
	File          Pattern  `json:"file"`
	Keyword       Pattern  `json:"keyword"`
	Configuration Pattern  `json:"configuration"`
	Grammar       []string `json:"grammar"`
}

// Cache struct
type Cache struct {
	Grammar []interface{} `json:"-"`
}

// Group struct
type Group struct {
	Name  string   `json:"name"`
	Tasks []string `json:"tasks"`
}

// Sanitize func
func (t *Task) Sanitize() Task {
	t.File = t.File.santize()
	t.Keyword = t.Keyword.santize()
	t.Configuration = t.Configuration.santize()
	return *t
}

// Comment struct
type Comment struct {
	Block  Block  `json:"block,omitempty"`
	Inline string `json:"inline,omitempty"`
}

// Block struct
type Block struct {
	Open  string `json:"open,omitempty"`
	Line  string `json:"line,omitempty"`
	Close string `json:"close,omitempty"`
}

// Pattern struct
type Pattern struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

func (p *Pattern) santize() Pattern {
	p.Include = deduplicate(p.Include)
	p.Exclude = deduplicate(p.Exclude)
	return *p
}

func deduplicate(values []string) (santized []string) {
	for _, v := range values {
		value := strings.TrimSpace(v)
		if len(value) > 0 {
			duplicate := false
			for _, vd := range santized {
				if vd == value {
					duplicate = true
					break
				}
			}
			if !duplicate {
				santized = append(santized, value)
			}
		}
	}
	return santized
}

// Open func
func Open() (file File, err error) {
	err = file.unmarshal()

	// Santize
	for i, t := range file.Tasks {
		file.Tasks[i] = t.Sanitize()
	}

	err = file.Write()

	return file, err
}

// Unmarshal func
func (f *File) unmarshal() (err error) {

	if !f.exists() {
		err := f.initialize()
		if err != nil {
			return err
		}
	}

	open, err := os.Open(name)
	if err != nil {
		return err
	}

	read, err := ioutil.ReadAll(open)
	if err != nil {
		return err
	}

	err = json.Unmarshal(read, &f)
	if err != nil {
		return err
	}

	return nil
}

func (f *File) initialize() (err error) {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return f.Write()
	}
	return nil
}

// Exists returns a bool based on emits.json existing or not.
func (f *File) exists() bool {
	info, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// HasTask returns a boolean if the task exists or not.
func (f *File) HasTask(task Task) bool {
	for _, t := range f.Tasks {
		if strings.ToLower(strings.TrimSpace(t.Name)) == strings.ToLower(strings.TrimSpace(task.Name)) {
			return true
		}
	}
	return false
}

// HasGroup returns a boolean if the group exists or not.
func (f *File) HasGroup(group Group) bool {
	for _, g := range f.Groups {
		if strings.ToLower(strings.TrimSpace(g.Name)) == strings.ToLower(strings.TrimSpace(group.Name)) {
			return true
		}
	}
	return false
}

// GetTask returns a Task by name
func (f *File) GetTask(task Task) Task {
	for _, t := range f.Tasks {
		if strings.ToLower(strings.TrimSpace(t.Name)) == strings.ToLower(strings.TrimSpace(task.Name)) {
			return t
		}
	}
	return Task{}
}

// GetGroup returns a Group by name
func (f *File) GetGroup(group Group) Group {
	for _, g := range f.Groups {
		if strings.ToLower(strings.TrimSpace(g.Name)) == strings.ToLower(strings.TrimSpace(group.Name)) {
			return g
		}
	}
	return Group{}
}

// CreateTask creates a task
func (f *File) CreateTask(task Task) (err error) {
	if !f.HasTask(task) {
		f.Tasks = append(f.Tasks, task)
		return nil
	}
	return fmt.Errorf("%s task already exists; cannot create task", task.Name)
}

// DeleteTask deletes a task by name
func (f *File) DeleteTask(task Task) (success bool) {
	for i, t := range f.Tasks {
		if t.Name == task.Name {
			f.Tasks = append(f.Tasks[:i], f.Tasks[i+1:]...)
			return true
		}
	}
	return false
}

// UpdateTask updates a task by name
func (f *File) UpdateTask(task Task) (err error) {
	for i, t := range f.Tasks {
		if t.Name == task.Name {
			f.Tasks[i] = task
			break
		}
	}
	return f.Write()
}

// ListTasks lists tasks and desriptions
func (f *File) ListTasks() {

}

// Write file
func (f *File) Write() (err error) {
	file, err := json.MarshalIndent(f, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(name, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Files file
func (t *Task) Files() (matches []string, err error) {
	// Includes
	var includePattern []string
	var includes []string
	for _, pattern := range t.File.Include {
		includePatternChecked := false
		for i := range includePattern {
			if includePattern[i] == pattern {
				includePatternChecked = true
				break
			}
		}
		if includePatternChecked == false {
			includePattern = append(includePattern, pattern)
			include, err := filepath.Glob(pattern)
			if err != nil {
				return nil, err
			}
			for _, name := range include {
				includes = append(includes, name)
			}
		}
	}
	includes = uniqueFile(includes)
	// Excludes
	var excludePattern []string
	var excludes []string
	for _, pattern := range t.File.Exclude {
		excludePatternChecked := false
		for i := range excludePattern {
			if excludePattern[i] == pattern {
				excludePatternChecked = true
				break
			}
		}
		if excludePatternChecked == false {
			excludePattern = append(excludePattern, pattern)
			exclude, err := filepath.Glob(pattern)
			if err != nil {
				return nil, err
			}
			for _, name := range exclude {
				excludes = append(excludes, name)
			}
		}
	}
	excludes = uniqueFile(excludes)
	// Includes - Excludes
	for _, exclude := range excludes {
		for i, include := range includes {
			if exclude == include {
				includes = append(includes[:i], includes[i+1:]...)
			}
		}
	}
	return includes, nil
}

func uniqueFile(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
