package attachment

import (
	"github.com/google/uuid"
)

type TextSource struct {
	Value string `toml:"value"`
}

type FileSource struct {
	Path string `toml:"path"`
	Mime string `toml:"mime"`
	Data []byte `toml:"data,omitempty"` // Optional for image data
}

type SymbolSource struct {
	Path  string      `toml:"path"`
	Name  string      `toml:"name"`
	Kind  int         `toml:"kind"`
	Range SymbolRange `toml:"range"`
}

type SymbolRange struct {
	Start Position `toml:"start"`
	End   Position `toml:"end"`
}

type Position struct {
	Line int `toml:"line"`
	Char int `toml:"char"`
}

type Attachment struct {
	ID         string `toml:"id"`
	Type       string `toml:"type"`
	Display    string `toml:"display"`
	URL        string `toml:"url"`
	Filename   string `toml:"filename"`
	MediaType  string `toml:"media_type"`
	StartIndex int    `toml:"start_index"`
	EndIndex   int    `toml:"end_index"`
	Source     any    `toml:"source,omitempty"`
}

// NewAttachment creates a new attachment with a unique ID
func NewAttachment() *Attachment {
	return &Attachment{
		ID: uuid.NewString(),
	}
}

func (a *Attachment) GetTextSource() (*TextSource, bool) {
	if a.Type != "text" {
		return nil, false
	}
	ts, ok := a.Source.(*TextSource)
	return ts, ok
}

// GetFileSource returns the source as FileSource if the attachment is a file type
func (a *Attachment) GetFileSource() (*FileSource, bool) {
	if a.Type != "file" {
		return nil, false
	}
	fs, ok := a.Source.(*FileSource)
	return fs, ok
}

// GetSymbolSource returns the source as SymbolSource if the attachment is a symbol type
func (a *Attachment) GetSymbolSource() (*SymbolSource, bool) {
	if a.Type != "symbol" {
		return nil, false
	}
	ss, ok := a.Source.(*SymbolSource)
	return ss, ok
}
