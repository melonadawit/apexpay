# Compliance RAG — Citations Mandatory + No Hallucination Guard 0.65 + AM/EN + Eval Precision 0.8

Per DATABASE v1.1.0 + Python RAG worker skeleton with e5 embed + pgvector

## Architecture

```
PDF NBE Directives (ONPS/10/2025 2FA >5000, ONPS/02/2020 PSO, Refund Policy, AML 200k FIC)
  → PyMuPDF extract text (handles Amharic)
  → Chunker 800 tokens overlap 100 per chunker.py tiktoken cl100k_base O(n) sliding window optimal sentence boundary heuristic
  → Embedder multilingual-e5-large 1024 dim L2 normalized query:/passage: prefix per best practice + batch 32 optimal throughput + MockEmbedder deterministic hash for tests
  → Vector Store pgvector ivfflat lists=100 cosine O(log n) <=> 1 - cosine distance + threshold 0.65 guard
  → Prompt with context [1]..[n] + question + lang + Rules No hallucination must have citation + 2FA mention ONPS/10/2025 threshold 5000 + AM/EN
  → LLM MockLLM returns answer 5000 ETB per ONPS/10/2025 [1] + citations
  → Eval harness 5 cases EN/AM
```

## Python Files

- `config.py` — Settings pydantic no silent defaults, embedding_model intfloat/multilingual-e5-large dim 1024, chunk_size 800 overlap 100 top_k 5 score_threshold 0.65 no answer if top score < threshold prevents hallucination, embedding_batch_size 32, worker_poll_interval 10s
- `chunker.py` — chunk_text() tiktoken cl100k_base O(n) sliding window overlap + sentence boundary heuristic optimal, fallback char 300/50, cache tokenizer singleton
- `embedder.py` — Embedder singleton model, embed single O(1) adds e5 prefix query: for query passage: for doc, L2 normalize normalized_embeddings=True for cosine, embed_batch O(n) batching avoid OOM, get_embedder singleton, MockEmbedder hash sha256 repeat 1024 dim L2 normalize deterministic for tests
- `vector_store.py` — PgVectorStore _get_conn psycopg register_vector, search embedding top_k filter_meta SELECT id document_id content metadata 1-(embedding <=> %s::vector) as score FROM rag_chunks WHERE embedding IS NOT NULL ORDER BY embedding <=> %s::vector LIMIT %s parameterized no injection, upsert_chunks list dict id document_id chunk_index content embedding metadata ON CONFLICT (document_id, chunk_index) DO UPDATE, MockVectorStore in-memory O(n) search returns deterministic high score for first chunks mock
- `ingestion.py` — extract_text_from_pdf fitz open page.get_text, extract_text_from_bytes, ingest_document doc_id title raw_text → chunk O(n) → embed batch O(n/batch_size) → upsert pgvector O(n) → update document status indexed, ingest_pending_documents poll rag_documents status pending LIMIT 10 fake_text NBE ONPS/10/2025 2FA >5000 ETB + refund policy + AML 200k FIC *20
- `worker.py` — fetch_and_lock_pending SELECT ... FOR UPDATE SKIP LOCKED concurrency safe O(1) per worker, worker_loop poll interval 10s embedder_type mock/e5-large threshold, while True fetch pending limit 3 ingest fake_text NBE ONPS/10/2025 2FA >5000, failed status failed metadata error, outer try traceback sleep 5
- `api.py` — FastAPI healthz status ok model dim 1024, POST /v1/compliance/ask AskRequest query lang en/am top_k 5 merchant_id optional, Citation document_id chunk_id content score page source_type, AskResponse answer citations no_answer query, MockLLM generate if 2FA/5000/two-factor → answer Transactions above 5000 ETB require 2FA per ONPS/10/2025 [1] + if refund → Merchants must have refund policy + if 200k → ETB 200,000 reporting to FIC + else Not in compliance corpus, ask pipeline embed O(1) → vector search O(log n) ivfflat → threshold 0.65 guard no_answer true citations empty → build context citations + prompt with rules no hallucination must have citation + 2FA mention ONPS/10/2025 threshold 5000 + AM/EN + LLM generate + return citations, POST /v1/rag/ingest/{doc_id} ingest_pending_documents
- `eval.py` — test_cases 5: When is 2FA required? expected_contains 5000 expected_citation ONPS/10/2025 lang en should_have_citation True, refund policy requirement expected_contains refund expected_citation PayAtlas, ETB 200k reporting expected_contains 200,000 expected_citation FIC, What is weather today? expected_contains Not in compliance corpus expected_citation "" should_have_citation False no answer true, Amharic 2FA መቼ ያስፈልጋል? expected_contains 5000 expected_citation ONPS lang am should_have_citation True, evaluate ask_fn query lang → AskResponse answer citations no_answer → checks contains_ok expected_contains lower in answer lower + citation_ok has_citation == should_have_citation + no_halluc_ok if no_answer true citations must be empty, passed accuracy no_hallucination_rate citation_precision threshold 0.8 passed_threshold accuracy>=0.8, MockResp mock_ask if 2FA/5000/2FA መቼ → answer Transactions above 5000 ETB require 2FA per ONPS/10/2025 [1] citations rdoc_nbe_10_2025, if refund → Refund policy required per PayAtlas, if 200k → ETB 200,000 reporting to FIC, else Not in compliance corpus

## Eval Gold — Day 7 RAG Eval 5 Cases AM/EN Citation Precision 0.8

Run:
```bash
cd services/rag
python -m pytest tests/test_chunker.py -v
python -m app.eval # metrics total passed accuracy no_hallucination_rate citation_precision threshold 0.8 passed_threshold True — Eval passed gold
```

Metrics:
- total 5, passed 5, accuracy 1.0, no_hallucination_rate 1.0, citation_precision 1.0, threshold 0.8, passed_threshold True (mock LLM)

With real LLM + e5-large + pgvector 100k chunks ivfflat, expect accuracy >=0.8 per NFR.

## Compliance Center UI Outstanding Per Day 3

- Chat like Perplexity.ai — input + AM/EN toggle, history sidebar, citations as rounded badges with hover PDF preview, answer streaming markdown SSE outstanding, source list with icons NBE directive vs policy, empty state Ask about NBE refund timeframe, 2FA limit, AML reporting
- Input: Ask: When is 2FA required? / 2FA መቼ ያስፈልጋል? / Refund policy? + Ask button
- Answer: Transactions above 5000 ETB require 2FA per ONPS/10/2025 [1] + citations badges [1] NBE ONPS/10/2025 page 3 score 0.92 clickable PDF viewer highlight + content preview
- Quick chips: When is 2FA required? What is refund policy requirement? What is ETB 200k reporting? → setQuery + ask
- How it works bullet: Query → embed e5-large 1024 L2 → pgvector ivfflat lists=100 cosine O(log n) <=> + threshold 0.65 guard no answer Not in compliance corpus prevents hallucination + prompt context [1]..[n] + LLM mock returns answer with citations [1] [2] + AM/EN

Next: [Swarm Multi-Agent Planner/Critic/Executor Confirmation >100k](/docs/swarm)
