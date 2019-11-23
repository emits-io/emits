package data

import "strings"

// Node structure used to support the Emits structure.
type Node struct {
	Appending     bool     `json:"-"`
	Collapsing    bool     `json:"-"`
	Newline       bool     `json:"-"`
	Configuration bool     `json:"-"`
	ParentNode    *Node    `json:"-"`
	Comment       Comment  `json:"-"`
	Parent        int      `json:"parent,omitempty"`
	Line          int      `json:"line,omitempty"`
	Index         int      `json:"index,omitempty"`
	Keyword       string   `json:"keyword,omitempty"`
	Value         string   `json:"value,omitempty"`
	Children      []Node   `json:"data,omitempty"`
	Separator     bool     `json:"separator,omitempty"`
	Flags         []string `json:"flags,omitempty"`
}

// Comment structure
type Comment struct {
	BlockOpen  bool `json:"blockOpen,omitempty"`
	BlockClose bool `json:"blockClose,omitempty"`
	BlockLine  bool `json:"blockLine,omitempty"`
	Inline     bool `json:"inline,omitempty"`
}

// AppendChild helper function appends a node to the children field.
func (n *Node) AppendChild(node Node) {
	n.Children = append(n.Children, node)
}

// HasChildren return true if the children length is greater than zero; must not be a configuration node.
func (n Node) HasChildren() bool {
	return len(n.Children) > 0 && !n.IsConfiguration()
}

// HasKeyword returns true if the keyword length is greater than zero; must not be a configuration node.
func (n Node) HasKeyword() bool {
	return len(n.Keyword) > 0 && !n.IsConfiguration()
}

// HasValue returns true if the value length is greater than zero; must not be a configuration node.
func (n Node) HasValue() bool {
	return len(n.Value) > 0 && !n.IsConfiguration()
}

// IsAppending returns the node appending value.
func (n *Node) IsAppending() bool {
	return n.Appending
}

// IsComment returns true if any of the comment types are true.
func (n *Node) IsComment() bool {
	return n.Comment.BlockOpen || n.Comment.BlockClose || n.Comment.BlockLine || n.Comment.Inline
}

// IsConfiguration returns the configuration (bool) value.
func (n Node) IsConfiguration() bool {
	return n.Configuration
}

// IsEmpty returns true if the keyword and value lengths are greater than zero; must not be a separator node.
func (n Node) IsEmpty() bool {
	return n.HasKeyword() == false && n.HasValue() == false && n.Separator == false || n.HasKeyword() == false && n.IsAppending() == false
}

// HasData returns true if the node has relevant data to output.
func (n Node) HasData() bool {
	data := false
	if n.HasKeyword() || n.HasValue() {
		data = true
	} else if n.Separator {
		data = true
	}
	return data
}

// LastNode (recursive) returns the last node in the tree structure.
func (n *Node) LastNode() (node *Node) {
	if n.HasChildren() {
		return n.Children[len(n.Children)-1].LastNode()
	}
	return n
}

// LastAppendingNodeOrRoot (recursive) returns the first node in the tree structure that is expecting appending values; or the tree root node.
func (n *Node) LastAppendingNodeOrRoot() (node *Node) {
	if n.ParentNode != nil && !n.IsAppending() {
		return n.ParentNode.LastAppendingNodeOrRoot()
	}
	return n
}

// PreviousIndexNode (recursive) returns the first node in the tree structure that is less than the current index; an exact match is not required in order to support stranded indexes.
func (n *Node) PreviousIndexNode(index int) (node *Node) {
	if n.ParentNode != nil && n.Index > index {
		return n.ParentNode.PreviousIndexNode(index)
	}
	return n
}

// HasInstanceOfKeyword (recursive) returns true for the first instance of
func (n *Node) HasInstanceOfKeyword(keyword []string) bool {
	for _, k := range keyword {
		if n.Keyword == k && len(k) > 0 {
			return true
		}
		for _, c := range n.Children {
			if c.HasInstanceOfKeyword(keyword) {
				return true
			}
		}
	}
	return false
}

// CollapseAppending func
func (n *Node) CollapseAppending() {
	for i, c := range n.Children {
		if c.Collapsing {
			value := c.collapseValues(c.Newline)
			if c.Newline {
				value = strings.TrimPrefix(value, "\n")
				offset := 0
				offset, value = cleanSpace(value)
				if offset > 0 {
					value = strings.ReplaceAll(value, "\n"+strings.Repeat(" ", offset), "\n")
				}
			}
			n.Children[i].Value = value
			n.Children[i].Children = nil
		} else {
			c.CollapseAppending()
		}
	}
}

func (n *Node) collapseValues(newline bool) (value string) {
	for _, c := range n.Children {
		linebreak := ""
		if newline {
			linebreak = "\n"
		}
		value += linebreak + strings.Repeat(" ", c.Index) + c.Value
		if c.HasChildren() {
			value += c.collapseValues(newline)
		}
	}
	return value
}
