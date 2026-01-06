package domain

type RequestPayloadType string

const (
	SimpleUtterance RequestPayloadType = "SimpleUtterance"
	ButtonPressed   RequestPayloadType = "ButtonPressed"
)

type RequestSession struct {
	New       bool   `json:"new"`
	MessageID int    `json:"message_id"`
	SessionID string `json:"session_id"`
	SkillID   string `json:"skill_id"`
	UserID    string `json:"user_id"`
}

type RequestMeta struct {
	Locale     string                 `json:"locale"`
	Timezone   string                 `json:"timezone"`
	ClientID   string                 `json:"client_id"`
	Interfaces map[string]interface{} `json:"interfaces"`
}

type RequestMarkup struct {
	DangerousContext bool `json:"dangerous_context"`
}

type Entity struct {
	Tokens EntityTokens `json:"tokens"`
	Type   EntityType   `json:"type"`
	Value  interface{}  `json:"value"`
}

type EntityTokens struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type EntityType string

const (
	EntityYandexFio      EntityType = "YANDEX.FIO"
	EntityYandexGeo      EntityType = "YANDEX.GEO"
	EntityYandexDateTime EntityType = "YANDEX.DATETIME"
	EntityYandexNumber   EntityType = "YANDEX.NUMBER"
	EntityYandexString   EntityType = "YANDEX.STRING"
)

type Slot struct {
	Type   EntityType   `json:"type"`
	Tokens EntityTokens `json:"tokens"`
	Value  interface{}  `json:"value"`
}

type SlotsClassAndDate struct {
	ClassNumber    Slot `json:"class_number,omitempty"`
	ClassCharacter Slot `json:"class_charachter,omitempty"`
	DayOfWeek      Slot `json:"day_of_week,omitempty"`
	Date           Slot `json:"date,omitempty"`
}

type IntentClassAndDate struct {
	Slots SlotsClassAndDate `json:"slots"`
}

type Intents struct {
	ClassAndDate IntentClassAndDate `json:"class_and_date,omitempty"`
}

type RequestNLU struct {
	Tokens   []string `json:"tokens"`
	Entities []Entity `json:"entities"`
	Intents  Intents  `json:"intents"`
}

type RequestPayload struct {
	Command           string             `json:"command"`
	OriginalUtterance string             `json:"original_utterance"`
	Type              RequestPayloadType `json:"type"`
	Markup            RequestMarkup      `json:"markup"`
	Payload           interface{}        `json:"payload,omitempty"`
	NLU               RequestNLU         `json:"nlu"`
}

type Event struct {
	Version string         `json:"version"`
	Meta    RequestMeta    `json:"meta"`
	Request RequestPayload `json:"request"`
	Session RequestSession `json:"session"`
}

type ResponseSession struct {
	SessionID string `json:"session_id"`
	MessageID int    `json:"message_id"`
	UserID    string `json:"user_id"`
}

type ResponseCardType string

const (
	BigImage  ResponseCardType = "BigImage"
	ItemsList ResponseCardType = "ItemsList"
)

type ResponseCard struct {
	Type ResponseCardType `json:"type"`

	// For BigImage type.
	*ResponseCardItem `json:",omitempty"`

	// For ItemsList type.
	*ResponseCardItemsList `json:",omitempty"`
}

type ResponseCardItemsList struct {
	Header *ResponseCardHeader `json:"header,omitempty"`
	Items  []ResponseCardItem  `json:"items,omitempty"`
	Footer *ResponseCardFooter `json:"footer,omitempty"`
}

type ResponseCardFooter struct {
	Text   string              `json:"text"`
	Button *ResponseCardButton `json:"button,omitempty"`
}

type ResponseCardButton struct {
	Text    string      `json:"text"`
	URL     string      `json:"url,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

type ResponseCardHeader struct {
	Text string `json:"text"`
}

type ResponseCardItem struct {
	ImageID     string              `json:"image_id,omitempty"`
	Title       string              `json:"title,omitempty"`
	Description string              `json:"description,omitempty"`
	Button      *ResponseCardButton `json:"button,omitempty"`
}

type ResponseButton struct {
	Title   string      `json:"title"`
	Payload interface{} `json:"payload,omitempty"`
	URL     string      `json:"url,omitempty"`
	Hide    bool        `json:"hide,omitempty"`
}

type ResponsePayload struct {
	Text       string           `json:"text"`
	Tts        string           `json:"tts,omitempty"`
	Card       *ResponseCard    `json:"card,omitempty"`
	Buttons    []ResponseButton `json:"buttons,omitempty"`
	EndSession bool             `json:"end_session"`
}

type Response struct {
	Response ResponsePayload `json:"response"`
	Session  ResponseSession `json:"session"`
	Version  string          `json:"version"`
}
