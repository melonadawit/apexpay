"""
RAG Worker - polls pending docs every 10s, ingestion + embedding
Optimal: single process with backoff, no duplicate processing via FOR UPDATE SKIP LOCKED
"""

import time
import psycopg
from psycopg.rows import dict_row
from .config import settings
from .ingestion import ingest_document
from pgvector.psycopg import register_vector
import traceback

def fetch_and_lock_pending(conn, limit=3):
    """Fetch pending docs with row-level lock SKIP LOCKED for concurrency safe O(1) per worker"""
    with conn.cursor() as cur:
        cur.execute("""
        SELECT id, title, metadata FROM rag_documents
        WHERE status='pending'
        ORDER BY created_at ASC
        FOR UPDATE SKIP LOCKED
        LIMIT %s
        """, (limit,))
        rows = cur.fetchall()
        # Mark as processing to avoid re-pick? Use indexed as final, pending as queue, we update to indexed after
        return rows

def worker_loop():
    print("[RAG Worker] Starting poll interval", settings.worker_poll_interval_seconds)
    embedder_type = "mock" if settings.env=="local" else "e5-large"
    print(f"[RAG Worker] Embedder: {embedder_type}, threshold {settings.score_threshold}")

    while True:
        try:
            with psycopg.connect(settings.database_url, row_factory=dict_row) as conn:
                register_vector(conn)
                docs = fetch_and_lock_pending(conn)
                if not docs:
                    time.sleep(settings.worker_poll_interval_seconds)
                    continue

                for doc in docs:
                    try:
                        print(f"[RAG Worker] Ingesting {doc['id']} {doc['title']}")
                        # Fake text for skeleton - real would download from MinIO/S3
                        fake_text = f"{doc['title']} content. " + "NBE ONPS/10/2025 2FA >5000 ETB PIN OTP biometric. " * 50
                        result = ingest_document(doc['id'], doc['title'], fake_text)
                        print(f"[RAG Worker] Indexed {result}")
                    except Exception as e:
                        print(f"[RAG Worker] Failed {doc['id']}: {e}")
                        traceback.print_exc()
                        with conn.cursor() as cur:
                            cur.execute("UPDATE rag_documents SET status='failed', metadata=jsonb_set(metadata, '{error}', %s) WHERE id=%s", (f'"{str(e)}"', doc['id']))
                        conn.commit()
        except Exception as outer:
            print(f"[RAG Worker] Outer error: {outer}")
            traceback.print_exc()
            time.sleep(5)

if __name__ == "__main__":
    worker_loop()
