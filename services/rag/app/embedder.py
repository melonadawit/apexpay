"""
Embedder: multilingual-e5-large 1024 dim optimal for Amharic + English per DATABASE v1.1.0
Best practice: singleton model, batch embedding O(batch) throughput, L2 normalize for cosine similarity pgvector.
Implements fallback to bge-m3 if e5 not available.
"""

from typing import List
import numpy as np
from functools import lru_cache

# Lazy imports for startup speed

class Embedder:
    def __init__(self, model_name: str = "intfloat/multilingual-e5-large"):
        self.model_name = model_name
        self._model = None
        self.dim = 1024

    def _load(self):
        if self._model is None:
            try:
                from sentence_transformers import SentenceTransformer
                # trust_remote_code=False for security
                self._model = SentenceTransformer(self.model_name, trust_remote_code=False)
                self.dim = self._model.get_sentence_embedding_dimension()
                print(f"[Embedder] Loaded {self.model_name} dim={self.dim}")
            except Exception as e:
                print(f"[Embedder] Failed to load {self.model_name}: {e}, trying BAAI/bge-m3")
                from sentence_transformers import SentenceTransformer
                self._model = SentenceTransformer("BAAI/bge-m3")
                self.dim = 1024

    def embed(self, text: str) -> List[float]:
        """Embed single text O(1) - adds e5 prefix query/passage per model spec"""
        self._load()
        # e5 requires prefix "query: " for query, "passage: " for doc
        prefixed = f"query: {text}" if len(text) < 500 else f"passage: {text}"
        emb = self._model.encode(prefixed, normalize_embeddings=True) # L2 norm for cosine
        return emb.tolist()

    def embed_batch(self, texts: List[str], batch_size: int = 32) -> List[List[float]]:
        """
        Batch embed optimal O(n) with batching to avoid OOM.
        Uses normalized embeddings for cosine similarity pgvector.
        """
        self._load()
        # Add e5 prefixes per best practice
        prefixed = [f"passage: {t}" for t in texts]
        embeddings = self._model.encode(
            prefixed,
            batch_size=batch_size,
            show_progress_bar=False,
            normalize_embeddings=True,
            convert_to_numpy=True,
        )
        return embeddings.tolist()

# Singleton optimal pattern
_embedder_instance = None

def get_embedder() -> Embedder:
    global _embedder_instance
    if _embedder_instance is None:
        from .config import settings
        _embedder_instance = Embedder(settings.embedding_model)
    return _embedder_instance

# Mock embedder for tests - deterministic hash based embedding to avoid model download
class MockEmbedder(Embedder):
    def __init__(self):
        super().__init__("mock")
        self.dim = 1024

    def embed(self, text: str) -> List[float]:
        # deterministic pseudo embedding from hash - optimal for tests
        import hashlib
        h = hashlib.sha256(text.encode()).digest()
        # repeat hash to fill 1024 dim float32 normalized
        vals = []
        for i in range(1024):
            vals.append(((h[i % len(h)] / 255.0) * 2 - 1))  # -1..1
        # L2 normalize
        arr = np.array(vals, dtype=np.float32)
        norm = np.linalg.norm(arr)
        if norm > 0:
            arr = arr / norm
        return arr.tolist()

    def embed_batch(self, texts: List[str], batch_size: int = 32) -> List[List[float]]:
        return [self.embed(t) for t in texts]
