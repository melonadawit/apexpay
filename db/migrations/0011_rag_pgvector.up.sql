-- 0011_rag_pgvector: RAG compliance layer with pgvector 1024 dim multilingual-e5
-- Requires pgvector extension, Python rag-worker ingestion

create extension if not exists vector;

create table rag_documents (
  id            text primary key,
  title         text not null,
  source_type   text not null check (source_type in ('nbe_directive','policy','faq','evidence','merchant_doc')),
  source_url    text,
  content_hash  text not null,
  status        text not null check (status in ('pending','indexed','failed')) default 'pending',
  metadata      jsonb not null default '{}'::jsonb,
  created_at    timestamptz not null default now()
);
create index rag_docs_status_idx on rag_documents (status) where status='pending';
create index rag_docs_source_type_idx on rag_documents (source_type);
create unique index rag_docs_content_hash_uidx on rag_documents (content_hash);

create table rag_chunks (
  id            text primary key,
  document_id   text not null references rag_documents(id) on delete cascade,
  chunk_index   int not null check (chunk_index >=0),
  content       text not null,
  content_en    text,
  content_am    text,
  embedding     vector(1024), -- bge-m3 / multilingual-e5-large 1024 dim optimal
  metadata      jsonb not null default '{}'::jsonb,
  created_at    timestamptz not null default now(),
  unique (document_id, chunk_index)
);
create index rag_chunks_doc_idx on rag_chunks (document_id);
-- ivfflat index for cosine similarity O(log n) search optimal for 100k chunks, lists=100
create index rag_chunks_embedding_idx on rag_chunks using ivfflat (embedding vector_cosine_ops) with (lists=100);

-- Seed sample NBE directive doc for testing RAG ask
insert into rag_documents (id, title, source_type, source_url, content_hash, status, metadata) values
('rdoc_nbe_10_2025', 'NBE ONPS/10/2025 Licensing Amendment - 2FA & Interoperability', 'nbe_directive', 'https://nbe.gov.et/directives/ONPS-10-2025.pdf', 'hash_nbe_10_2025_sample', 'pending', '{"page_count": 12, "year": 2025}'::jsonb),
('rdoc_refund_policy', 'ApexPay Refund Policy', 'policy', 'https://apexpay.et/policies/refund', 'hash_refund_policy_sample', 'pending', '{"version": "1.0"}'::jsonb),
('rdoc_aml_200k', 'FIC AML Reporting ETB 200k Threshold', 'nbe_directive', 'https://fic.gov.et/aml', 'hash_aml_200k', 'pending', '{"threshold": 200000}'::jsonb)
on conflict (id) do nothing;

comment on table rag_chunks is 'RAG compliance: chunk 800 tokens overlap 100, embed multilingual-e5/bge-m3 1024 dim, citations mandatory, no hallucination guard score threshold 0.65, AM/EN support.';
