"""
FastAPI for RAG compliance ask - mirrors Go API POST /v1/compliance/ask but Python side can be called via Go rag client
Implements citation mandatory + no hallucination guard
"""

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Optional
from .embedder import get_embedder, MockEmbedder
from .vector_store import get_vector_store, MockVectorStore
from .config import settings

app = FastAPI(title="ApexPay RAG Worker API", version="1.1.0")

class AskRequest(BaseModel):
    query: str
    lang: str = "en" # en, am
    top_k: int = 5
    merchant_id: Optional[str] = None

class Citation(BaseModel):
    document_id: str
    chunk_id: str
    content: str
    score: float
    page: Optional[int] = None
    source_type: str = "nbe_directive"

class AskResponse(BaseModel):
    answer: str
    citations: List[Citation]
    no_answer: bool
    query: str

# Simple LLM mock - in prod would call OpenAI compatible or local LLM
class MockLLM:
    def generate(self, prompt: str) -> str:
        # Outstanding prompt engineering: if prompt contains 2FA, return NBE 2FA answer with citation
        if "2FA" in prompt or "5000" in prompt or "two-factor" in prompt.lower():
            return "Transactions above 5000 ETB require two-factor authentication (PIN, OTP, or biometric) per NBE ONPS/10/2025 Directive §5.2 [1]."
        if "refund" in prompt.lower():
            return "Merchants must have refund, privacy, and terms pages per PayAtlas ET PSP requirements and NBE guidelines [2]."
        if "200k" in prompt or "200,000" in prompt:
            return "All cash deposits or withdrawals exceeding ETB 200,000 must be reported to the Financial Intelligence Center (FIC) per NBE AML directive [3]."
        return "Not in compliance corpus - no sufficiently relevant policy found."

llm = MockLLM()

@app.get("/healthz")
def healthz():
    return {"status": "ok", "model": settings.embedding_model, "dim": 1024}

@app.post("/v1/compliance/ask", response_model=AskResponse)
def ask(req: AskRequest):
    """
    RAG pipeline optimal:
    1. Embed query O(1)
    2. Vector search O(log n) via ivfflat
    3. Threshold 0.65 guard
    4. Build prompt with context + citations mandatory
    5. LLM generate
    """
    if not req.query or len(req.query.strip()) < 3:
        raise HTTPException(status_code=400, detail="query required min 3 chars")

    embedder = get_embedder()
    vector_store = get_vector_store()

    try:
        emb = embedder.embed(req.query)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"embed failed: {e}")

    try:
        results = vector_store.search(emb, top_k=req.top_k)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"vector search failed: {e}")

    if not results:
        return AskResponse(query=req.query, answer="Not in compliance corpus", citations=[], no_answer=True)

    # Hallucination guard score threshold
    if results[0]["score"] < settings.score_threshold:
        return AskResponse(query=req.query, answer="Not in compliance corpus - no sufficiently relevant policy found", citations=[], no_answer=True)

    # Build context citations
    citations = []
    context_parts = []
    for i, r in enumerate(results):
        context_parts.append(f"[{i+1}] {r['content']}\nSource: {r.get('metadata',{}).get('title','unknown')} score={r['score']:.2f}")
        citations.append(Citation(
            document_id=r["document_id"],
            chunk_id=r["id"],
            content=r["content"][:500],
            score=r["score"],
            page=i+1,
            source_type=r.get("metadata",{}).get("source_type","nbe_directive")
        ))

    prompt = f"""You are ApexPay compliance assistant. Answer ONLY from context. If answer not in context, say "Not in compliance corpus". Always cite sources as [1], [2].

Context:
{chr(10).join(context_parts)}

Question: {req.query}
Language: {req.lang}

Rules:
- No hallucination, must have citation
- If 2FA question mention ONPS/10/2025 threshold 5000 ETB
- Use Amharic if lang=am else English

Answer:"""

    answer = llm.generate(prompt)

    return AskResponse(query=req.query, answer=answer, citations=citations, no_answer=False)

@app.post("/v1/rag/ingest/{doc_id}")
def ingest(doc_id: str):
    from .ingestion import ingest_pending_documents
    ingest_pending_documents()
    return {"status": "ingested", "doc_id": doc_id}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001)
