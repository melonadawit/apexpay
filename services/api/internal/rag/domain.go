package rag

import "time"

type DocumentStatus string
const (
	StatusPending DocumentStatus = "pending"
	StatusIndexed DocumentStatus = "indexed"
	StatusFailed  DocumentStatus = "failed"
)

type SourceType string
const (
	SourceNBEDirective SourceType = "nbe_directive"
	SourcePolicy       SourceType = "policy"
	SourceFAQ          SourceType = "faq"
	SourceEvidence     SourceType = "evidence"
)

type Document struct {
	ID          string
	Title       string
	SourceType  SourceType
	SourceURL   string
	ContentHash string
	Status      DocumentStatus
	Metadata    map[string]any
	CreatedAt   time.Time
}

type Chunk struct {
	ID         string
	DocumentID string
	ChunkIndex int
	Content    string
	ContentEN  string
	ContentAM  string
	Embedding  []float32 // 1024 dim for pgvector
	Metadata   map[string]any
	CreatedAt  time.Time
}

type AskRequest struct {
	Query     string `validate:"required,min=3"`
	Lang      string `validate:"omitempty,oneof=en am"` // default en
	MerchantID string
	TopK      int // default 5
}

type Citation struct {
	DocumentID   string
	DocumentTitle string
	ChunkID      string
	Content      string
	Score        float64
	Page         *int
	SourceURL    string
	SourceType   SourceType
}

type AskResponse struct {
	Answer    string
	Citations []Citation
	NoAnswer  bool // if not in corpus
	Query     string
}
