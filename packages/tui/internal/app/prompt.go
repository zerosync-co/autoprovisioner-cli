package app

import (
	"time"

	"github.com/sst/opencode-sdk-go"
	"github.com/sst/opencode/internal/attachment"
	"github.com/sst/opencode/internal/id"
)

type Prompt struct {
	Text        string                   `toml:"text"`
	Attachments []*attachment.Attachment `toml:"attachments"`
}

func (p Prompt) ToMessage(
	messageID string,
	sessionID string,
) Message {
	message := opencode.UserMessage{
		ID:        messageID,
		SessionID: sessionID,
		Role:      opencode.UserMessageRoleUser,
		Time: opencode.UserMessageTime{
			Created: float64(time.Now().UnixMilli()),
		},
	}

	text := p.Text
	textAttachments := []*attachment.Attachment{}
	for _, attachment := range p.Attachments {
		if attachment.Type == "text" {
			textAttachments = append(textAttachments, attachment)
		}
	}
	for i := 0; i < len(textAttachments)-1; i++ {
		for j := i + 1; j < len(textAttachments); j++ {
			if textAttachments[i].StartIndex < textAttachments[j].StartIndex {
				textAttachments[i], textAttachments[j] = textAttachments[j], textAttachments[i]
			}
		}
	}
	for _, att := range textAttachments {
		source, _ := att.GetTextSource()
		text = text[:att.StartIndex] + source.Value + text[att.EndIndex:]
	}

	parts := []opencode.PartUnion{opencode.TextPart{
		ID:        id.Ascending(id.Part),
		MessageID: messageID,
		SessionID: sessionID,
		Type:      opencode.TextPartTypeText,
		Text:      text,
	}}
	for _, attachment := range p.Attachments {
		text := opencode.FilePartSourceText{
			Start: int64(attachment.StartIndex),
			End:   int64(attachment.EndIndex),
			Value: attachment.Display,
		}
		var source *opencode.FilePartSource
		switch attachment.Type {
		case "text":
			continue
		case "file":
			fileSource, _ := attachment.GetFileSource()
			source = &opencode.FilePartSource{
				Text: text,
				Path: fileSource.Path,
				Type: opencode.FilePartSourceTypeFile,
			}
		case "symbol":
			symbolSource, _ := attachment.GetSymbolSource()
			source = &opencode.FilePartSource{
				Text: text,
				Path: symbolSource.Path,
				Type: opencode.FilePartSourceTypeSymbol,
				Kind: int64(symbolSource.Kind),
				Name: symbolSource.Name,
				Range: opencode.SymbolSourceRange{
					Start: opencode.SymbolSourceRangeStart{
						Line:      float64(symbolSource.Range.Start.Line),
						Character: float64(symbolSource.Range.Start.Char),
					},
					End: opencode.SymbolSourceRangeEnd{
						Line:      float64(symbolSource.Range.End.Line),
						Character: float64(symbolSource.Range.End.Char),
					},
				},
			}
		}
		parts = append(parts, opencode.FilePart{
			ID:        id.Ascending(id.Part),
			MessageID: messageID,
			SessionID: sessionID,
			Type:      opencode.FilePartTypeFile,
			Filename:  attachment.Filename,
			Mime:      attachment.MediaType,
			URL:       attachment.URL,
			Source:    *source,
		})
	}
	return Message{
		Info:  message,
		Parts: parts,
	}
}

func (m Message) ToSessionChatParams() []opencode.SessionChatParamsPartUnion {
	parts := []opencode.SessionChatParamsPartUnion{}
	for _, part := range m.Parts {
		switch p := part.(type) {
		case opencode.TextPart:
			parts = append(parts, opencode.TextPartInputParam{
				ID:        opencode.F(p.ID),
				Type:      opencode.F(opencode.TextPartInputTypeText),
				Text:      opencode.F(p.Text),
				Synthetic: opencode.F(p.Synthetic),
				Time: opencode.F(opencode.TextPartInputTimeParam{
					Start: opencode.F(p.Time.Start),
					End:   opencode.F(p.Time.End),
				}),
			})
		case opencode.FilePart:
			var source opencode.FilePartSourceUnionParam
			switch p.Source.Type {
			case "file":
				source = opencode.FileSourceParam{
					Type: opencode.F(opencode.FileSourceTypeFile),
					Path: opencode.F(p.Source.Path),
					Text: opencode.F(opencode.FilePartSourceTextParam{
						Start: opencode.F(int64(p.Source.Text.Start)),
						End:   opencode.F(int64(p.Source.Text.End)),
						Value: opencode.F(p.Source.Text.Value),
					}),
				}
			case "symbol":
				source = opencode.SymbolSourceParam{
					Type: opencode.F(opencode.SymbolSourceTypeSymbol),
					Path: opencode.F(p.Source.Path),
					Name: opencode.F(p.Source.Name),
					Kind: opencode.F(p.Source.Kind),
					Range: opencode.F(opencode.SymbolSourceRangeParam{
						Start: opencode.F(opencode.SymbolSourceRangeStartParam{
							Line:      opencode.F(float64(p.Source.Range.(opencode.SymbolSourceRange).Start.Line)),
							Character: opencode.F(float64(p.Source.Range.(opencode.SymbolSourceRange).Start.Character)),
						}),
						End: opencode.F(opencode.SymbolSourceRangeEndParam{
							Line:      opencode.F(float64(p.Source.Range.(opencode.SymbolSourceRange).End.Line)),
							Character: opencode.F(float64(p.Source.Range.(opencode.SymbolSourceRange).End.Character)),
						}),
					}),
					Text: opencode.F(opencode.FilePartSourceTextParam{
						Value: opencode.F(p.Source.Text.Value),
						Start: opencode.F(p.Source.Text.Start),
						End:   opencode.F(p.Source.Text.End),
					}),
				}
			}
			parts = append(parts, opencode.FilePartInputParam{
				ID:       opencode.F(p.ID),
				Type:     opencode.F(opencode.FilePartInputTypeFile),
				Mime:     opencode.F(p.Mime),
				URL:      opencode.F(p.URL),
				Filename: opencode.F(p.Filename),
				Source:   opencode.F(source),
			})
		}
	}
	return parts
}

func (p Prompt) ToSessionChatParams() []opencode.SessionChatParamsPartUnion {
	parts := []opencode.SessionChatParamsPartUnion{
		opencode.TextPartInputParam{
			Type: opencode.F(opencode.TextPartInputTypeText),
			Text: opencode.F(p.Text),
		},
	}
	for _, att := range p.Attachments {
		filePart := opencode.FilePartInputParam{
			Type:     opencode.F(opencode.FilePartInputTypeFile),
			Mime:     opencode.F(att.MediaType),
			URL:      opencode.F(att.URL),
			Filename: opencode.F(att.Filename),
		}
		switch att.Type {
		case "file":
			if fs, ok := att.GetFileSource(); ok {
				filePart.Source = opencode.F(
					opencode.FilePartSourceUnionParam(opencode.FileSourceParam{
						Type: opencode.F(opencode.FileSourceTypeFile),
						Path: opencode.F(fs.Path),
						Text: opencode.F(opencode.FilePartSourceTextParam{
							Start: opencode.F(int64(att.StartIndex)),
							End:   opencode.F(int64(att.EndIndex)),
							Value: opencode.F(att.Display),
						}),
					}),
				)
			}
		case "symbol":
			if ss, ok := att.GetSymbolSource(); ok {
				filePart.Source = opencode.F(
					opencode.FilePartSourceUnionParam(opencode.SymbolSourceParam{
						Type: opencode.F(opencode.SymbolSourceTypeSymbol),
						Path: opencode.F(ss.Path),
						Name: opencode.F(ss.Name),
						Kind: opencode.F(int64(ss.Kind)),
						Range: opencode.F(opencode.SymbolSourceRangeParam{
							Start: opencode.F(opencode.SymbolSourceRangeStartParam{
								Line:      opencode.F(float64(ss.Range.Start.Line)),
								Character: opencode.F(float64(ss.Range.Start.Char)),
							}),
							End: opencode.F(opencode.SymbolSourceRangeEndParam{
								Line:      opencode.F(float64(ss.Range.End.Line)),
								Character: opencode.F(float64(ss.Range.End.Char)),
							}),
						}),
						Text: opencode.F(opencode.FilePartSourceTextParam{
							Start: opencode.F(int64(att.StartIndex)),
							End:   opencode.F(int64(att.EndIndex)),
							Value: opencode.F(att.Display),
						}),
					}),
				)
			}
		}
		parts = append(parts, filePart)
	}
	return parts
}
