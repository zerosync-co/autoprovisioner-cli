// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package opencode

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/sst/opencode-sdk-go/internal/apijson"
	"github.com/sst/opencode-sdk-go/internal/param"
	"github.com/sst/opencode-sdk-go/internal/requestconfig"
	"github.com/sst/opencode-sdk-go/option"
	"github.com/sst/opencode-sdk-go/shared"
	"github.com/tidwall/gjson"
)

// SessionService contains methods and other services that help with interacting
// with the opencode API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSessionService] method instead.
type SessionService struct {
	Options []option.RequestOption
}

// NewSessionService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewSessionService(opts ...option.RequestOption) (r *SessionService) {
	r = &SessionService{}
	r.Options = opts
	return
}

// Create a new session
func (r *SessionService) New(ctx context.Context, opts ...option.RequestOption) (res *Session, err error) {
	opts = append(r.Options[:], opts...)
	path := "session"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return
}

// List all sessions
func (r *SessionService) List(ctx context.Context, opts ...option.RequestOption) (res *[]Session, err error) {
	opts = append(r.Options[:], opts...)
	path := "session"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Delete a session and all its data
func (r *SessionService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (res *bool, err error) {
	opts = append(r.Options[:], opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("session/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// Abort a session
func (r *SessionService) Abort(ctx context.Context, id string, opts ...option.RequestOption) (res *bool, err error) {
	opts = append(r.Options[:], opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("session/%s/abort", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return
}

// Create and send a new message to a session
func (r *SessionService) Chat(ctx context.Context, id string, body SessionChatParams, opts ...option.RequestOption) (res *AssistantMessage, err error) {
	opts = append(r.Options[:], opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("session/%s/message", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Analyze the app and create an AGENTS.md file
func (r *SessionService) Init(ctx context.Context, id string, body SessionInitParams, opts ...option.RequestOption) (res *bool, err error) {
	opts = append(r.Options[:], opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("session/%s/init", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// List messages for a session
func (r *SessionService) Messages(ctx context.Context, id string, opts ...option.RequestOption) (res *[]SessionMessagesResponse, err error) {
	opts = append(r.Options[:], opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("session/%s/message", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Share a session
func (r *SessionService) Share(ctx context.Context, id string, opts ...option.RequestOption) (res *Session, err error) {
	opts = append(r.Options[:], opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("session/%s/share", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return
}

// Summarize the session
func (r *SessionService) Summarize(ctx context.Context, id string, body SessionSummarizeParams, opts ...option.RequestOption) (res *bool, err error) {
	opts = append(r.Options[:], opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("session/%s/summarize", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Unshare the session
func (r *SessionService) Unshare(ctx context.Context, id string, opts ...option.RequestOption) (res *Session, err error) {
	opts = append(r.Options[:], opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("session/%s/share", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

type AssistantMessage struct {
	ID         string                 `json:"id,required"`
	Cost       float64                `json:"cost,required"`
	ModelID    string                 `json:"modelID,required"`
	Path       AssistantMessagePath   `json:"path,required"`
	ProviderID string                 `json:"providerID,required"`
	Role       AssistantMessageRole   `json:"role,required"`
	SessionID  string                 `json:"sessionID,required"`
	System     []string               `json:"system,required"`
	Time       AssistantMessageTime   `json:"time,required"`
	Tokens     AssistantMessageTokens `json:"tokens,required"`
	Error      AssistantMessageError  `json:"error"`
	Summary    bool                   `json:"summary"`
	JSON       assistantMessageJSON   `json:"-"`
}

// assistantMessageJSON contains the JSON metadata for the struct
// [AssistantMessage]
type assistantMessageJSON struct {
	ID          apijson.Field
	Cost        apijson.Field
	ModelID     apijson.Field
	Path        apijson.Field
	ProviderID  apijson.Field
	Role        apijson.Field
	SessionID   apijson.Field
	System      apijson.Field
	Time        apijson.Field
	Tokens      apijson.Field
	Error       apijson.Field
	Summary     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantMessage) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantMessageJSON) RawJSON() string {
	return r.raw
}

func (r AssistantMessage) implementsMessage() {}

type AssistantMessagePath struct {
	Cwd  string                   `json:"cwd,required"`
	Root string                   `json:"root,required"`
	JSON assistantMessagePathJSON `json:"-"`
}

// assistantMessagePathJSON contains the JSON metadata for the struct
// [AssistantMessagePath]
type assistantMessagePathJSON struct {
	Cwd         apijson.Field
	Root        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantMessagePath) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantMessagePathJSON) RawJSON() string {
	return r.raw
}

type AssistantMessageRole string

const (
	AssistantMessageRoleAssistant AssistantMessageRole = "assistant"
)

func (r AssistantMessageRole) IsKnown() bool {
	switch r {
	case AssistantMessageRoleAssistant:
		return true
	}
	return false
}

type AssistantMessageTime struct {
	Created   float64                  `json:"created,required"`
	Completed float64                  `json:"completed"`
	JSON      assistantMessageTimeJSON `json:"-"`
}

// assistantMessageTimeJSON contains the JSON metadata for the struct
// [AssistantMessageTime]
type assistantMessageTimeJSON struct {
	Created     apijson.Field
	Completed   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantMessageTime) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantMessageTimeJSON) RawJSON() string {
	return r.raw
}

type AssistantMessageTokens struct {
	Cache     AssistantMessageTokensCache `json:"cache,required"`
	Input     float64                     `json:"input,required"`
	Output    float64                     `json:"output,required"`
	Reasoning float64                     `json:"reasoning,required"`
	JSON      assistantMessageTokensJSON  `json:"-"`
}

// assistantMessageTokensJSON contains the JSON metadata for the struct
// [AssistantMessageTokens]
type assistantMessageTokensJSON struct {
	Cache       apijson.Field
	Input       apijson.Field
	Output      apijson.Field
	Reasoning   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantMessageTokens) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantMessageTokensJSON) RawJSON() string {
	return r.raw
}

type AssistantMessageTokensCache struct {
	Read  float64                         `json:"read,required"`
	Write float64                         `json:"write,required"`
	JSON  assistantMessageTokensCacheJSON `json:"-"`
}

// assistantMessageTokensCacheJSON contains the JSON metadata for the struct
// [AssistantMessageTokensCache]
type assistantMessageTokensCacheJSON struct {
	Read        apijson.Field
	Write       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantMessageTokensCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantMessageTokensCacheJSON) RawJSON() string {
	return r.raw
}

type AssistantMessageError struct {
	// This field can have the runtime type of [shared.ProviderAuthErrorData],
	// [shared.UnknownErrorData], [interface{}].
	Data  interface{}               `json:"data,required"`
	Name  AssistantMessageErrorName `json:"name,required"`
	JSON  assistantMessageErrorJSON `json:"-"`
	union AssistantMessageErrorUnion
}

// assistantMessageErrorJSON contains the JSON metadata for the struct
// [AssistantMessageError]
type assistantMessageErrorJSON struct {
	Data        apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r assistantMessageErrorJSON) RawJSON() string {
	return r.raw
}

func (r *AssistantMessageError) UnmarshalJSON(data []byte) (err error) {
	*r = AssistantMessageError{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [AssistantMessageErrorUnion] interface which you can cast to
// the specific types for more type safety.
//
// Possible runtime types of the union are [shared.ProviderAuthError],
// [shared.UnknownError], [AssistantMessageErrorMessageOutputLengthError],
// [shared.MessageAbortedError].
func (r AssistantMessageError) AsUnion() AssistantMessageErrorUnion {
	return r.union
}

// Union satisfied by [shared.ProviderAuthError], [shared.UnknownError],
// [AssistantMessageErrorMessageOutputLengthError] or [shared.MessageAbortedError].
type AssistantMessageErrorUnion interface {
	ImplementsAssistantMessageError()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*AssistantMessageErrorUnion)(nil)).Elem(),
		"name",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(shared.ProviderAuthError{}),
			DiscriminatorValue: "ProviderAuthError",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(shared.UnknownError{}),
			DiscriminatorValue: "UnknownError",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantMessageErrorMessageOutputLengthError{}),
			DiscriminatorValue: "MessageOutputLengthError",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(shared.MessageAbortedError{}),
			DiscriminatorValue: "MessageAbortedError",
		},
	)
}

type AssistantMessageErrorMessageOutputLengthError struct {
	Data interface{}                                       `json:"data,required"`
	Name AssistantMessageErrorMessageOutputLengthErrorName `json:"name,required"`
	JSON assistantMessageErrorMessageOutputLengthErrorJSON `json:"-"`
}

// assistantMessageErrorMessageOutputLengthErrorJSON contains the JSON metadata for
// the struct [AssistantMessageErrorMessageOutputLengthError]
type assistantMessageErrorMessageOutputLengthErrorJSON struct {
	Data        apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantMessageErrorMessageOutputLengthError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantMessageErrorMessageOutputLengthErrorJSON) RawJSON() string {
	return r.raw
}

func (r AssistantMessageErrorMessageOutputLengthError) ImplementsAssistantMessageError() {}

type AssistantMessageErrorMessageOutputLengthErrorName string

const (
	AssistantMessageErrorMessageOutputLengthErrorNameMessageOutputLengthError AssistantMessageErrorMessageOutputLengthErrorName = "MessageOutputLengthError"
)

func (r AssistantMessageErrorMessageOutputLengthErrorName) IsKnown() bool {
	switch r {
	case AssistantMessageErrorMessageOutputLengthErrorNameMessageOutputLengthError:
		return true
	}
	return false
}

type AssistantMessageErrorName string

const (
	AssistantMessageErrorNameProviderAuthError        AssistantMessageErrorName = "ProviderAuthError"
	AssistantMessageErrorNameUnknownError             AssistantMessageErrorName = "UnknownError"
	AssistantMessageErrorNameMessageOutputLengthError AssistantMessageErrorName = "MessageOutputLengthError"
	AssistantMessageErrorNameMessageAbortedError      AssistantMessageErrorName = "MessageAbortedError"
)

func (r AssistantMessageErrorName) IsKnown() bool {
	switch r {
	case AssistantMessageErrorNameProviderAuthError, AssistantMessageErrorNameUnknownError, AssistantMessageErrorNameMessageOutputLengthError, AssistantMessageErrorNameMessageAbortedError:
		return true
	}
	return false
}

type FilePart struct {
	ID        string       `json:"id,required"`
	MessageID string       `json:"messageID,required"`
	Mime      string       `json:"mime,required"`
	SessionID string       `json:"sessionID,required"`
	Type      FilePartType `json:"type,required"`
	URL       string       `json:"url,required"`
	Filename  string       `json:"filename"`
	JSON      filePartJSON `json:"-"`
}

// filePartJSON contains the JSON metadata for the struct [FilePart]
type filePartJSON struct {
	ID          apijson.Field
	MessageID   apijson.Field
	Mime        apijson.Field
	SessionID   apijson.Field
	Type        apijson.Field
	URL         apijson.Field
	Filename    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FilePart) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r filePartJSON) RawJSON() string {
	return r.raw
}

func (r FilePart) implementsPart() {}

type FilePartType string

const (
	FilePartTypeFile FilePartType = "file"
)

func (r FilePartType) IsKnown() bool {
	switch r {
	case FilePartTypeFile:
		return true
	}
	return false
}

type FilePartParam struct {
	ID        param.Field[string]       `json:"id,required"`
	MessageID param.Field[string]       `json:"messageID,required"`
	Mime      param.Field[string]       `json:"mime,required"`
	SessionID param.Field[string]       `json:"sessionID,required"`
	Type      param.Field[FilePartType] `json:"type,required"`
	URL       param.Field[string]       `json:"url,required"`
	Filename  param.Field[string]       `json:"filename"`
}

func (r FilePartParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r FilePartParam) implementsSessionChatParamsPartUnion() {}

type Message struct {
	ID        string      `json:"id,required"`
	Role      MessageRole `json:"role,required"`
	SessionID string      `json:"sessionID,required"`
	// This field can have the runtime type of [UserMessageTime],
	// [AssistantMessageTime].
	Time interface{} `json:"time,required"`
	Cost float64     `json:"cost"`
	// This field can have the runtime type of [AssistantMessageError].
	Error   interface{} `json:"error"`
	ModelID string      `json:"modelID"`
	// This field can have the runtime type of [AssistantMessagePath].
	Path       interface{} `json:"path"`
	ProviderID string      `json:"providerID"`
	Summary    bool        `json:"summary"`
	// This field can have the runtime type of [[]string].
	System interface{} `json:"system"`
	// This field can have the runtime type of [AssistantMessageTokens].
	Tokens interface{} `json:"tokens"`
	JSON   messageJSON `json:"-"`
	union  MessageUnion
}

// messageJSON contains the JSON metadata for the struct [Message]
type messageJSON struct {
	ID          apijson.Field
	Role        apijson.Field
	SessionID   apijson.Field
	Time        apijson.Field
	Cost        apijson.Field
	Error       apijson.Field
	ModelID     apijson.Field
	Path        apijson.Field
	ProviderID  apijson.Field
	Summary     apijson.Field
	System      apijson.Field
	Tokens      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r messageJSON) RawJSON() string {
	return r.raw
}

func (r *Message) UnmarshalJSON(data []byte) (err error) {
	*r = Message{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [MessageUnion] interface which you can cast to the specific
// types for more type safety.
//
// Possible runtime types of the union are [UserMessage], [AssistantMessage].
func (r Message) AsUnion() MessageUnion {
	return r.union
}

// Union satisfied by [UserMessage] or [AssistantMessage].
type MessageUnion interface {
	implementsMessage()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*MessageUnion)(nil)).Elem(),
		"role",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(UserMessage{}),
			DiscriminatorValue: "user",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantMessage{}),
			DiscriminatorValue: "assistant",
		},
	)
}

type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
)

func (r MessageRole) IsKnown() bool {
	switch r {
	case MessageRoleUser, MessageRoleAssistant:
		return true
	}
	return false
}

type Part struct {
	ID        string   `json:"id,required"`
	MessageID string   `json:"messageID,required"`
	SessionID string   `json:"sessionID,required"`
	Type      PartType `json:"type,required"`
	CallID    string   `json:"callID"`
	Cost      float64  `json:"cost"`
	Filename  string   `json:"filename"`
	Mime      string   `json:"mime"`
	Snapshot  string   `json:"snapshot"`
	// This field can have the runtime type of [ToolPartState].
	State     interface{} `json:"state"`
	Synthetic bool        `json:"synthetic"`
	Text      string      `json:"text"`
	// This field can have the runtime type of [TextPartTime].
	Time interface{} `json:"time"`
	// This field can have the runtime type of [StepFinishPartTokens].
	Tokens interface{} `json:"tokens"`
	Tool   string      `json:"tool"`
	URL    string      `json:"url"`
	JSON   partJSON    `json:"-"`
	union  PartUnion
}

// partJSON contains the JSON metadata for the struct [Part]
type partJSON struct {
	ID          apijson.Field
	MessageID   apijson.Field
	SessionID   apijson.Field
	Type        apijson.Field
	CallID      apijson.Field
	Cost        apijson.Field
	Filename    apijson.Field
	Mime        apijson.Field
	Snapshot    apijson.Field
	State       apijson.Field
	Synthetic   apijson.Field
	Text        apijson.Field
	Time        apijson.Field
	Tokens      apijson.Field
	Tool        apijson.Field
	URL         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r partJSON) RawJSON() string {
	return r.raw
}

func (r *Part) UnmarshalJSON(data []byte) (err error) {
	*r = Part{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [PartUnion] interface which you can cast to the specific types
// for more type safety.
//
// Possible runtime types of the union are [TextPart], [FilePart], [ToolPart],
// [StepStartPart], [StepFinishPart], [SnapshotPart].
func (r Part) AsUnion() PartUnion {
	return r.union
}

// Union satisfied by [TextPart], [FilePart], [ToolPart], [StepStartPart],
// [StepFinishPart] or [SnapshotPart].
type PartUnion interface {
	implementsPart()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*PartUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(TextPart{}),
			DiscriminatorValue: "text",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(FilePart{}),
			DiscriminatorValue: "file",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ToolPart{}),
			DiscriminatorValue: "tool",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(StepStartPart{}),
			DiscriminatorValue: "step-start",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(StepFinishPart{}),
			DiscriminatorValue: "step-finish",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SnapshotPart{}),
			DiscriminatorValue: "snapshot",
		},
	)
}

type PartType string

const (
	PartTypeText       PartType = "text"
	PartTypeFile       PartType = "file"
	PartTypeTool       PartType = "tool"
	PartTypeStepStart  PartType = "step-start"
	PartTypeStepFinish PartType = "step-finish"
	PartTypeSnapshot   PartType = "snapshot"
)

func (r PartType) IsKnown() bool {
	switch r {
	case PartTypeText, PartTypeFile, PartTypeTool, PartTypeStepStart, PartTypeStepFinish, PartTypeSnapshot:
		return true
	}
	return false
}

type Session struct {
	ID       string        `json:"id,required"`
	Time     SessionTime   `json:"time,required"`
	Title    string        `json:"title,required"`
	Version  string        `json:"version,required"`
	ParentID string        `json:"parentID"`
	Revert   SessionRevert `json:"revert"`
	Share    SessionShare  `json:"share"`
	JSON     sessionJSON   `json:"-"`
}

// sessionJSON contains the JSON metadata for the struct [Session]
type sessionJSON struct {
	ID          apijson.Field
	Time        apijson.Field
	Title       apijson.Field
	Version     apijson.Field
	ParentID    apijson.Field
	Revert      apijson.Field
	Share       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *Session) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionJSON) RawJSON() string {
	return r.raw
}

type SessionTime struct {
	Created float64         `json:"created,required"`
	Updated float64         `json:"updated,required"`
	JSON    sessionTimeJSON `json:"-"`
}

// sessionTimeJSON contains the JSON metadata for the struct [SessionTime]
type sessionTimeJSON struct {
	Created     apijson.Field
	Updated     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SessionTime) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionTimeJSON) RawJSON() string {
	return r.raw
}

type SessionRevert struct {
	MessageID string            `json:"messageID,required"`
	Part      float64           `json:"part,required"`
	Snapshot  string            `json:"snapshot"`
	JSON      sessionRevertJSON `json:"-"`
}

// sessionRevertJSON contains the JSON metadata for the struct [SessionRevert]
type sessionRevertJSON struct {
	MessageID   apijson.Field
	Part        apijson.Field
	Snapshot    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SessionRevert) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionRevertJSON) RawJSON() string {
	return r.raw
}

type SessionShare struct {
	URL  string           `json:"url,required"`
	JSON sessionShareJSON `json:"-"`
}

// sessionShareJSON contains the JSON metadata for the struct [SessionShare]
type sessionShareJSON struct {
	URL         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SessionShare) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionShareJSON) RawJSON() string {
	return r.raw
}

type SnapshotPart struct {
	ID        string           `json:"id,required"`
	MessageID string           `json:"messageID,required"`
	SessionID string           `json:"sessionID,required"`
	Snapshot  string           `json:"snapshot,required"`
	Type      SnapshotPartType `json:"type,required"`
	JSON      snapshotPartJSON `json:"-"`
}

// snapshotPartJSON contains the JSON metadata for the struct [SnapshotPart]
type snapshotPartJSON struct {
	ID          apijson.Field
	MessageID   apijson.Field
	SessionID   apijson.Field
	Snapshot    apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SnapshotPart) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r snapshotPartJSON) RawJSON() string {
	return r.raw
}

func (r SnapshotPart) implementsPart() {}

type SnapshotPartType string

const (
	SnapshotPartTypeSnapshot SnapshotPartType = "snapshot"
)

func (r SnapshotPartType) IsKnown() bool {
	switch r {
	case SnapshotPartTypeSnapshot:
		return true
	}
	return false
}

type StepFinishPart struct {
	ID        string               `json:"id,required"`
	Cost      float64              `json:"cost,required"`
	MessageID string               `json:"messageID,required"`
	SessionID string               `json:"sessionID,required"`
	Tokens    StepFinishPartTokens `json:"tokens,required"`
	Type      StepFinishPartType   `json:"type,required"`
	JSON      stepFinishPartJSON   `json:"-"`
}

// stepFinishPartJSON contains the JSON metadata for the struct [StepFinishPart]
type stepFinishPartJSON struct {
	ID          apijson.Field
	Cost        apijson.Field
	MessageID   apijson.Field
	SessionID   apijson.Field
	Tokens      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *StepFinishPart) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r stepFinishPartJSON) RawJSON() string {
	return r.raw
}

func (r StepFinishPart) implementsPart() {}

type StepFinishPartTokens struct {
	Cache     StepFinishPartTokensCache `json:"cache,required"`
	Input     float64                   `json:"input,required"`
	Output    float64                   `json:"output,required"`
	Reasoning float64                   `json:"reasoning,required"`
	JSON      stepFinishPartTokensJSON  `json:"-"`
}

// stepFinishPartTokensJSON contains the JSON metadata for the struct
// [StepFinishPartTokens]
type stepFinishPartTokensJSON struct {
	Cache       apijson.Field
	Input       apijson.Field
	Output      apijson.Field
	Reasoning   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *StepFinishPartTokens) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r stepFinishPartTokensJSON) RawJSON() string {
	return r.raw
}

type StepFinishPartTokensCache struct {
	Read  float64                       `json:"read,required"`
	Write float64                       `json:"write,required"`
	JSON  stepFinishPartTokensCacheJSON `json:"-"`
}

// stepFinishPartTokensCacheJSON contains the JSON metadata for the struct
// [StepFinishPartTokensCache]
type stepFinishPartTokensCacheJSON struct {
	Read        apijson.Field
	Write       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *StepFinishPartTokensCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r stepFinishPartTokensCacheJSON) RawJSON() string {
	return r.raw
}

type StepFinishPartType string

const (
	StepFinishPartTypeStepFinish StepFinishPartType = "step-finish"
)

func (r StepFinishPartType) IsKnown() bool {
	switch r {
	case StepFinishPartTypeStepFinish:
		return true
	}
	return false
}

type StepStartPart struct {
	ID        string            `json:"id,required"`
	MessageID string            `json:"messageID,required"`
	SessionID string            `json:"sessionID,required"`
	Type      StepStartPartType `json:"type,required"`
	JSON      stepStartPartJSON `json:"-"`
}

// stepStartPartJSON contains the JSON metadata for the struct [StepStartPart]
type stepStartPartJSON struct {
	ID          apijson.Field
	MessageID   apijson.Field
	SessionID   apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *StepStartPart) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r stepStartPartJSON) RawJSON() string {
	return r.raw
}

func (r StepStartPart) implementsPart() {}

type StepStartPartType string

const (
	StepStartPartTypeStepStart StepStartPartType = "step-start"
)

func (r StepStartPartType) IsKnown() bool {
	switch r {
	case StepStartPartTypeStepStart:
		return true
	}
	return false
}

type TextPart struct {
	ID        string       `json:"id,required"`
	MessageID string       `json:"messageID,required"`
	SessionID string       `json:"sessionID,required"`
	Text      string       `json:"text,required"`
	Type      TextPartType `json:"type,required"`
	Synthetic bool         `json:"synthetic"`
	Time      TextPartTime `json:"time"`
	JSON      textPartJSON `json:"-"`
}

// textPartJSON contains the JSON metadata for the struct [TextPart]
type textPartJSON struct {
	ID          apijson.Field
	MessageID   apijson.Field
	SessionID   apijson.Field
	Text        apijson.Field
	Type        apijson.Field
	Synthetic   apijson.Field
	Time        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *TextPart) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r textPartJSON) RawJSON() string {
	return r.raw
}

func (r TextPart) implementsPart() {}

type TextPartType string

const (
	TextPartTypeText TextPartType = "text"
)

func (r TextPartType) IsKnown() bool {
	switch r {
	case TextPartTypeText:
		return true
	}
	return false
}

type TextPartTime struct {
	Start float64          `json:"start,required"`
	End   float64          `json:"end"`
	JSON  textPartTimeJSON `json:"-"`
}

// textPartTimeJSON contains the JSON metadata for the struct [TextPartTime]
type textPartTimeJSON struct {
	Start       apijson.Field
	End         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *TextPartTime) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r textPartTimeJSON) RawJSON() string {
	return r.raw
}

type TextPartParam struct {
	ID        param.Field[string]            `json:"id,required"`
	MessageID param.Field[string]            `json:"messageID,required"`
	SessionID param.Field[string]            `json:"sessionID,required"`
	Text      param.Field[string]            `json:"text,required"`
	Type      param.Field[TextPartType]      `json:"type,required"`
	Synthetic param.Field[bool]              `json:"synthetic"`
	Time      param.Field[TextPartTimeParam] `json:"time"`
}

func (r TextPartParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r TextPartParam) implementsSessionChatParamsPartUnion() {}

type TextPartTimeParam struct {
	Start param.Field[float64] `json:"start,required"`
	End   param.Field[float64] `json:"end"`
}

func (r TextPartTimeParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ToolPart struct {
	ID        string        `json:"id,required"`
	CallID    string        `json:"callID,required"`
	MessageID string        `json:"messageID,required"`
	SessionID string        `json:"sessionID,required"`
	State     ToolPartState `json:"state,required"`
	Tool      string        `json:"tool,required"`
	Type      ToolPartType  `json:"type,required"`
	JSON      toolPartJSON  `json:"-"`
}

// toolPartJSON contains the JSON metadata for the struct [ToolPart]
type toolPartJSON struct {
	ID          apijson.Field
	CallID      apijson.Field
	MessageID   apijson.Field
	SessionID   apijson.Field
	State       apijson.Field
	Tool        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ToolPart) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r toolPartJSON) RawJSON() string {
	return r.raw
}

func (r ToolPart) implementsPart() {}

type ToolPartState struct {
	Status ToolPartStateStatus `json:"status,required"`
	Error  string              `json:"error"`
	// This field can have the runtime type of [interface{}], [map[string]interface{}].
	Input interface{} `json:"input"`
	// This field can have the runtime type of [map[string]interface{}].
	Metadata interface{} `json:"metadata"`
	Output   string      `json:"output"`
	// This field can have the runtime type of [ToolStateRunningTime],
	// [ToolStateCompletedTime], [ToolStateErrorTime].
	Time  interface{}       `json:"time"`
	Title string            `json:"title"`
	JSON  toolPartStateJSON `json:"-"`
	union ToolPartStateUnion
}

// toolPartStateJSON contains the JSON metadata for the struct [ToolPartState]
type toolPartStateJSON struct {
	Status      apijson.Field
	Error       apijson.Field
	Input       apijson.Field
	Metadata    apijson.Field
	Output      apijson.Field
	Time        apijson.Field
	Title       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r toolPartStateJSON) RawJSON() string {
	return r.raw
}

func (r *ToolPartState) UnmarshalJSON(data []byte) (err error) {
	*r = ToolPartState{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [ToolPartStateUnion] interface which you can cast to the
// specific types for more type safety.
//
// Possible runtime types of the union are [ToolStatePending], [ToolStateRunning],
// [ToolStateCompleted], [ToolStateError].
func (r ToolPartState) AsUnion() ToolPartStateUnion {
	return r.union
}

// Union satisfied by [ToolStatePending], [ToolStateRunning], [ToolStateCompleted]
// or [ToolStateError].
type ToolPartStateUnion interface {
	implementsToolPartState()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ToolPartStateUnion)(nil)).Elem(),
		"status",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ToolStatePending{}),
			DiscriminatorValue: "pending",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ToolStateRunning{}),
			DiscriminatorValue: "running",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ToolStateCompleted{}),
			DiscriminatorValue: "completed",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ToolStateError{}),
			DiscriminatorValue: "error",
		},
	)
}

type ToolPartStateStatus string

const (
	ToolPartStateStatusPending   ToolPartStateStatus = "pending"
	ToolPartStateStatusRunning   ToolPartStateStatus = "running"
	ToolPartStateStatusCompleted ToolPartStateStatus = "completed"
	ToolPartStateStatusError     ToolPartStateStatus = "error"
)

func (r ToolPartStateStatus) IsKnown() bool {
	switch r {
	case ToolPartStateStatusPending, ToolPartStateStatusRunning, ToolPartStateStatusCompleted, ToolPartStateStatusError:
		return true
	}
	return false
}

type ToolPartType string

const (
	ToolPartTypeTool ToolPartType = "tool"
)

func (r ToolPartType) IsKnown() bool {
	switch r {
	case ToolPartTypeTool:
		return true
	}
	return false
}

type ToolStateCompleted struct {
	Input    map[string]interface{}   `json:"input,required"`
	Metadata map[string]interface{}   `json:"metadata,required"`
	Output   string                   `json:"output,required"`
	Status   ToolStateCompletedStatus `json:"status,required"`
	Time     ToolStateCompletedTime   `json:"time,required"`
	Title    string                   `json:"title,required"`
	JSON     toolStateCompletedJSON   `json:"-"`
}

// toolStateCompletedJSON contains the JSON metadata for the struct
// [ToolStateCompleted]
type toolStateCompletedJSON struct {
	Input       apijson.Field
	Metadata    apijson.Field
	Output      apijson.Field
	Status      apijson.Field
	Time        apijson.Field
	Title       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ToolStateCompleted) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r toolStateCompletedJSON) RawJSON() string {
	return r.raw
}

func (r ToolStateCompleted) implementsToolPartState() {}

type ToolStateCompletedStatus string

const (
	ToolStateCompletedStatusCompleted ToolStateCompletedStatus = "completed"
)

func (r ToolStateCompletedStatus) IsKnown() bool {
	switch r {
	case ToolStateCompletedStatusCompleted:
		return true
	}
	return false
}

type ToolStateCompletedTime struct {
	End   float64                    `json:"end,required"`
	Start float64                    `json:"start,required"`
	JSON  toolStateCompletedTimeJSON `json:"-"`
}

// toolStateCompletedTimeJSON contains the JSON metadata for the struct
// [ToolStateCompletedTime]
type toolStateCompletedTimeJSON struct {
	End         apijson.Field
	Start       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ToolStateCompletedTime) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r toolStateCompletedTimeJSON) RawJSON() string {
	return r.raw
}

type ToolStateError struct {
	Error  string                 `json:"error,required"`
	Input  map[string]interface{} `json:"input,required"`
	Status ToolStateErrorStatus   `json:"status,required"`
	Time   ToolStateErrorTime     `json:"time,required"`
	JSON   toolStateErrorJSON     `json:"-"`
}

// toolStateErrorJSON contains the JSON metadata for the struct [ToolStateError]
type toolStateErrorJSON struct {
	Error       apijson.Field
	Input       apijson.Field
	Status      apijson.Field
	Time        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ToolStateError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r toolStateErrorJSON) RawJSON() string {
	return r.raw
}

func (r ToolStateError) implementsToolPartState() {}

type ToolStateErrorStatus string

const (
	ToolStateErrorStatusError ToolStateErrorStatus = "error"
)

func (r ToolStateErrorStatus) IsKnown() bool {
	switch r {
	case ToolStateErrorStatusError:
		return true
	}
	return false
}

type ToolStateErrorTime struct {
	End   float64                `json:"end,required"`
	Start float64                `json:"start,required"`
	JSON  toolStateErrorTimeJSON `json:"-"`
}

// toolStateErrorTimeJSON contains the JSON metadata for the struct
// [ToolStateErrorTime]
type toolStateErrorTimeJSON struct {
	End         apijson.Field
	Start       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ToolStateErrorTime) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r toolStateErrorTimeJSON) RawJSON() string {
	return r.raw
}

type ToolStatePending struct {
	Status ToolStatePendingStatus `json:"status,required"`
	JSON   toolStatePendingJSON   `json:"-"`
}

// toolStatePendingJSON contains the JSON metadata for the struct
// [ToolStatePending]
type toolStatePendingJSON struct {
	Status      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ToolStatePending) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r toolStatePendingJSON) RawJSON() string {
	return r.raw
}

func (r ToolStatePending) implementsToolPartState() {}

type ToolStatePendingStatus string

const (
	ToolStatePendingStatusPending ToolStatePendingStatus = "pending"
)

func (r ToolStatePendingStatus) IsKnown() bool {
	switch r {
	case ToolStatePendingStatusPending:
		return true
	}
	return false
}

type ToolStateRunning struct {
	Status   ToolStateRunningStatus `json:"status,required"`
	Time     ToolStateRunningTime   `json:"time,required"`
	Input    interface{}            `json:"input"`
	Metadata map[string]interface{} `json:"metadata"`
	Title    string                 `json:"title"`
	JSON     toolStateRunningJSON   `json:"-"`
}

// toolStateRunningJSON contains the JSON metadata for the struct
// [ToolStateRunning]
type toolStateRunningJSON struct {
	Status      apijson.Field
	Time        apijson.Field
	Input       apijson.Field
	Metadata    apijson.Field
	Title       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ToolStateRunning) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r toolStateRunningJSON) RawJSON() string {
	return r.raw
}

func (r ToolStateRunning) implementsToolPartState() {}

type ToolStateRunningStatus string

const (
	ToolStateRunningStatusRunning ToolStateRunningStatus = "running"
)

func (r ToolStateRunningStatus) IsKnown() bool {
	switch r {
	case ToolStateRunningStatusRunning:
		return true
	}
	return false
}

type ToolStateRunningTime struct {
	Start float64                  `json:"start,required"`
	JSON  toolStateRunningTimeJSON `json:"-"`
}

// toolStateRunningTimeJSON contains the JSON metadata for the struct
// [ToolStateRunningTime]
type toolStateRunningTimeJSON struct {
	Start       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ToolStateRunningTime) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r toolStateRunningTimeJSON) RawJSON() string {
	return r.raw
}

type UserMessage struct {
	ID        string          `json:"id,required"`
	Role      UserMessageRole `json:"role,required"`
	SessionID string          `json:"sessionID,required"`
	Time      UserMessageTime `json:"time,required"`
	JSON      userMessageJSON `json:"-"`
}

// userMessageJSON contains the JSON metadata for the struct [UserMessage]
type userMessageJSON struct {
	ID          apijson.Field
	Role        apijson.Field
	SessionID   apijson.Field
	Time        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *UserMessage) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r userMessageJSON) RawJSON() string {
	return r.raw
}

func (r UserMessage) implementsMessage() {}

type UserMessageRole string

const (
	UserMessageRoleUser UserMessageRole = "user"
)

func (r UserMessageRole) IsKnown() bool {
	switch r {
	case UserMessageRoleUser:
		return true
	}
	return false
}

type UserMessageTime struct {
	Created float64             `json:"created,required"`
	JSON    userMessageTimeJSON `json:"-"`
}

// userMessageTimeJSON contains the JSON metadata for the struct [UserMessageTime]
type userMessageTimeJSON struct {
	Created     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *UserMessageTime) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r userMessageTimeJSON) RawJSON() string {
	return r.raw
}

type SessionMessagesResponse struct {
	Info  Message                     `json:"info,required"`
	Parts []Part                      `json:"parts,required"`
	JSON  sessionMessagesResponseJSON `json:"-"`
}

// sessionMessagesResponseJSON contains the JSON metadata for the struct
// [SessionMessagesResponse]
type sessionMessagesResponseJSON struct {
	Info        apijson.Field
	Parts       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SessionMessagesResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionMessagesResponseJSON) RawJSON() string {
	return r.raw
}

type SessionChatParams struct {
	MessageID  param.Field[string]                       `json:"messageID,required"`
	Mode       param.Field[string]                       `json:"mode,required"`
	ModelID    param.Field[string]                       `json:"modelID,required"`
	Parts      param.Field[[]SessionChatParamsPartUnion] `json:"parts,required"`
	ProviderID param.Field[string]                       `json:"providerID,required"`
}

func (r SessionChatParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SessionChatParamsPart struct {
	ID        param.Field[string]                     `json:"id,required"`
	MessageID param.Field[string]                     `json:"messageID,required"`
	SessionID param.Field[string]                     `json:"sessionID,required"`
	Type      param.Field[SessionChatParamsPartsType] `json:"type,required"`
	Filename  param.Field[string]                     `json:"filename"`
	Mime      param.Field[string]                     `json:"mime"`
	Synthetic param.Field[bool]                       `json:"synthetic"`
	Text      param.Field[string]                     `json:"text"`
	Time      param.Field[interface{}]                `json:"time"`
	URL       param.Field[string]                     `json:"url"`
}

func (r SessionChatParamsPart) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SessionChatParamsPart) implementsSessionChatParamsPartUnion() {}

// Satisfied by [FilePartParam], [TextPartParam], [SessionChatParamsPart].
type SessionChatParamsPartUnion interface {
	implementsSessionChatParamsPartUnion()
}

type SessionChatParamsPartsType string

const (
	SessionChatParamsPartsTypeFile SessionChatParamsPartsType = "file"
	SessionChatParamsPartsTypeText SessionChatParamsPartsType = "text"
)

func (r SessionChatParamsPartsType) IsKnown() bool {
	switch r {
	case SessionChatParamsPartsTypeFile, SessionChatParamsPartsTypeText:
		return true
	}
	return false
}

type SessionInitParams struct {
	MessageID  param.Field[string] `json:"messageID,required"`
	ModelID    param.Field[string] `json:"modelID,required"`
	ProviderID param.Field[string] `json:"providerID,required"`
}

func (r SessionInitParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SessionSummarizeParams struct {
	ModelID    param.Field[string] `json:"modelID,required"`
	ProviderID param.Field[string] `json:"providerID,required"`
}

func (r SessionSummarizeParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
