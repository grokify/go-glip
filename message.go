package glip

type GlipWebhookMessage struct {
	Icon           string       `json:"icon,omitempty"`
	Activity       string       `json:"activity,omitempty"`
	Title          string       `json:"title,omitempty"`
	Body           string       `json:"body,omitempty"`
	AttachmentType string       `json:"attachment_type,omitempty"`
	Attachments    []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Type         string  `json:"card,omitempty"`
	Color        string  `json:"color,omitempty"`
	Pretext      string  `json:"pretext,omitempty"`
	AuthorName   string  `json:"author_name,omitempty"`
	AuthorLink   string  `json:"author_link,omitempty"`
	AuthorIcon   string  `json:"author_icon,omitempty"`
	Title        string  `json:"title,omitempty"`
	TitleLink    string  `json:"title_link,omitempty"`
	Fallback     string  `json:"fallback,omitempty"`
	Fields       []Field `json:"fields,omitempty"`
	Text         string  `json:"text,omitempty"`
	ImageURL     string  `json:"image_url,omitempty"`
	ThumbnailURL string  `json:"thumbnail_url,omitempty"`
	Footer       string  `json:"footer,omitempty"`
	FooterIcon   string  `json:"footer_icon,omitempty"`
	TS           int64   `json:"ts,omitempty"`
}

type Author struct {
	Name    string `json:"name,omitempty"`
	URI     string `json:"uri,omitempty"`
	IconURI string `json:"iconUri,omitempty"`
}

type Footnote struct {
	Text    string `json:"text,omitempty"`
	IconURI string `json:"iconUri,omitempty"`
}

type Field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
	Style string `json:"style,omitempty"`
}
