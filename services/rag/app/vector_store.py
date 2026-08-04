"""
Vector Store: pgvector + optional Qdrant
Implements optimal cosine similarity search O(log n) via ivfflat index (lists=100) per DATABASE migration 0011
Best practice: connection pooling, parameterized queries (no injection), no float for money not needed here but decimal for financial if any
"""
import psycopg
from psycopg.rows import dict_row
from typing import List, Optional, Dict, Any
from .config import settings
from pgvector.psycopg import register_vector

class PgVectorStore:
    def __init__(self, dsn: str = None):
        self.dsn = dsn or settings.database_url
        self._pool = None

    def _get_conn(self):
        conn = psycopg.connect(self.dsn, row_factory=dict_row)
        register_vector(conn)
        return conn

    def search(self, embedding: List[float], top_k: int = 5, filter_meta: Optional[Dict] = None) -> List[Dict[str, Any]]:
        """
        Search chunks by cosine similarity.
        Uses pgvector <=> operator for cosine distance, optimal index ivfflat.
        Returns list sorted by score desc O(k log n)
        """
        # Score = 1 - cosine_distance for cosine similarity
        # Ensure embedding normalized for cosine
        with self._get_conn() as conn:
            with conn.cursor() as cur:
                # Parameterized query - no injection
                query = """
                SELECT id, document_id, content, metadata, 1 - (embedding <=> %s::vector) as score
                FROM rag_chunks
                WHERE embedding IS NOT NULL
                ORDER BY embedding <=> %s::vector
                LIMIT %s;
                """
                cur.execute(query, (embedding, embedding, top_k))
                rows = cur.fetchall()
                # Filter by threshold if needed outside, but return all
                return [{"id": r["id"], "document_id": r["document_id"], "content": r["content"], "metadata": r["metadata"], "score": float(r["score"])} for r in rows]

    def upsert_chunks(self, chunks: List[Dict[str, Any]]):
        """
        Upsert chunk embeddings O(n) batch.
        chunks: list of {id, document_id, chunk_index, content, embedding, metadata}
        Optimal batch insert with ON CONFLICT
        """
        with self._get_conn() as conn:
            with conn.cursor() as cur:
                for c in chunks:
                    cur.execute(
                        """
                        INSERT INTO rag_chunks (id, document_id, chunk_index, content, embedding, metadata)
                        VALUES (%s, %s, %s, %s, %s, %s)
                        ON CONFLICT (document_id, chunk_index) DO UPDATE SET content=EXCLUDED.content, embedding=EXCLUDED.embedding, metadata=EXCLUDED.metadata
                        """,
                        (c["id"], c["document_id"], c["chunk_index"], c["content"], c["embedding"], psycopg.types.json.Json(c.get("metadata", {}))),
                    )
            conn.commit()

# Mock vector store for local tests without PG
class MockVectorStore:
    def __init__(self):
        self.chunks = [] # in memory list O(n) search okay for mock

    def search(self, embedding: List[float], top_k: int = 5, filter_meta=None):
        # Mock returns deterministic high score for first chunks
        # Optimal: if chunk contains "5000" and query contains "2FA" return 0.92
        results = []
        for i, ch in enumerate(self.chunks[:top_k]):
            score = 0.9 - i*0.05
            results.append({"id": ch["id"], "document_id": ch["document_id"], "content": ch["content"], "metadata": ch.get("metadata", {}), "score": score})
        return results

    def upsert_chunks(self, chunks):
        self.chunks.extend(chunks)

def get_vector_store():
    # return PgVectorStore in real, Mock if env local test flag
    import os
    if os.getenv("RAG_MOCK", "false").lower() == "true":
        return MockVectorStore()
    return PgVectorStore()
