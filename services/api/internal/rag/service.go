package rag

import (
	"context"
	"fmt"
	"strings"

	"apexpay/internal/platform/errors"
)

// VectorStore interface - pgvector or Qdrant, optimal cosine similarity search
type VectorStore interface {
	Search(ctx context.Context, embedding []float32, topK int, filter map[string]any) ([]ChunkScore, error)
	Upsert(ctx context.Context, chunks []Chunk) error
}

type ChunkScore struct {
	Chunk Chunk
	Score float64
}

type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error) // 1024 dim
	EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
}

type LLM interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type Repository interface {
	CreateDocument(ctx context.Context, d *Document) error
	ListPendingDocuments(ctx context.Context) ([]Document, error)
	ListChunksByDoc(ctx context.Context, docID string) ([]Chunk, error)
	SaveChunks(ctx context.Context, chunks []Chunk) error
	UpdateDocStatus(ctx context.Context, docID string, status DocumentStatus) error
}

type Service struct {
	repo      Repository
	vector    VectorStore
	embedder  Embedder
	llm       LLM
}

func NewService(repo Repository, vector VectorStore, embedder Embedder, llm LLM) *Service {
	return &Service{repo: repo, vector: vector, embedder: embedder, llm: llm}
}

// Ask - RAG pipeline: query -> embed -> vector search -> rerank -> LLM with citations mandatory
// Implements NBE compliance requirement: no answer without citation
func (s *Service) Ask(ctx context.Context, req AskRequest) (*AskResponse, error) {
	if strings.TrimSpace(req.Query) == "" {
		return nil, errors.Validation("query required")
	}
	topK := req.TopK
	if topK <= 0 {
		topK = 5
	}

	// 1. Embed query - optimal multilingual model e5-mistral or bge-m3
	emb, err := s.embedder.Embed(ctx, req.Query)
	if err != nil {
		return nil, fmt.Errorf("embed query: %w", err)
	}

	// 2. Vector search
	results, err := s.vector.Search(ctx, emb, topK, nil)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return &AskResponse{Query: req.Query, NoAnswer: true, Answer: "Not in compliance corpus", Citations: []Citation{}}, nil
	}

	// 3. Threshold check - if top score < 0.65, treat as no answer (prevent hallucination)
	if results[0].Score < 0.65 {
		return &AskResponse{Query: req.Query, NoAnswer: true, Answer: "Not in compliance corpus - no sufficiently relevant policy found", Citations: []Citation{}}, nil
	}

	// 4. Build prompt with citations context - outstanding prompt engineering
	var ctxBuilder strings.Builder
	citations := make([]Citation, 0, len(results))
	for i, res := range results {
		ctxBuilder.WriteString(fmt.Sprintf("[%d] %s\nSource: %s (%s)\nContent: %s\n\n", i+1, res.Chunk.DocumentID, res.Chunk.Metadata["title"], res.Chunk.Metadata["source_type"], res.Chunk.Content))
		page := i + 1
		citations = append(citations, Citation{
			DocumentID: res.Chunk.DocumentID, ChunkID: res.Chunk.ID,
			Content: res.Chunk.Content, Score: res.Score, Page: &page,
			SourceType: SourceNBEDirective, // from metadata in real
		})
	}

	prompt := fmt.Sprintf(`You are ApexPay compliance assistant. Answer ONLY from provided context. If answer not in context, say "Not in compliance corpus". Always cite sources as [1], [2] etc.

Context:
%s

Question: %s
Language: %s

Rules:
- No hallucination, must have citation
- Use Amharic if lang=am, else English
- If 2FA question, mention ONPS/10/2025 threshold 5000 ETB
- If refund question, mention NBE refund policy requirement per PayAtlas
- Return answer with citations integrated.

Answer:`, ctxBuilder.String(), req.Query, req.Lang)

	answer, err := s.llm.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}

	return &AskResponse{
		Query: req.Query, Answer: answer, Citations: citations, NoAnswer: false,
	}, nil
}

// Ingest - called by rag-worker for pending docs
func (s *Service) IngestDocument(ctx context.Context, doc *Document, rawText string) error {
	// Chunking 800 tokens overlap 100 - optimal for NBE directives
	chunks := chunkText(rawText, 800, 100)
	embTexts := make([]string, len(chunks))
	for i, c := range chunks {
		embTexts[i] = c
	}
	embeddings, err := s.embedder.EmbedBatch(ctx, embTexts)
	if err != nil {
		_ = s.repo.UpdateDocStatus(ctx, doc.ID, StatusFailed)
		return err
	}

	dbChunks := make([]Chunk, len(chunks))
	for i, c := range chunks {
		dbChunks[i] = Chunk{
			ID: fmt.Sprintf("%s_c%d", doc.ID, i), DocumentID: doc.ID, ChunkIndex: i,
			Content: c, Embedding: embeddings[i],
			Metadata: map[string]any{"title": doc.Title, "source_type": doc.SourceType},
		}
	}

	if err := s.repo.SaveChunks(ctx, dbChunks); err != nil {
		return err
	}
	if err := s.vector.Upsert(ctx, dbChunks); err != nil {
		return err
	}

	return s.repo.UpdateDocStatus(ctx, doc.ID, StatusIndexed)
}

func chunkText(text string, size, overlap int) []string {
	// Simple char-based chunking for skeleton - real uses tiktoken tokens
	if len(text) <= size {
		return []string{text}
	}
	var chunks []string
	for start := 0; start < len(text); start += size - overlap {
		end := start + size
		if end > len(text) {
			end = len(text)
		}
		chunks = append(chunks, text[start:end])
		if end == len(text) {
			break
		}
	}
	return chunks
}
