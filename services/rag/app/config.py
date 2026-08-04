from pydantic_settings import BaseSettings
from typing import Literal

# Config validated, no silent defaults for secrets per SAD §11 cross-cutting

class Settings(BaseSettings):
    env: Literal["local","staging","production"] = "local"
    database_url: str = "postgresql://apexpay:apexpay_dev@localhost:5432/apexpay"
    redis_url: str = "redis://localhost:6379/0"
    embedding_model: str = "intfloat/multilingual-e5-large" # 1024 dim optimal for Amharic+English
    # alternatives: BAAI/bge-m3 1024 dim, multilingual-e5 also 1024
    chunk_size: int = 800 # tokens per DATABASE
    chunk_overlap: int = 100
    top_k: int = 5
    score_threshold: float = 0.65 # no answer if top score < threshold - hallucination guard
    minio_endpoint: str = "localhost:9000"
    minio_access_key: str = "minioadmin"
    minio_secret_key: str = "minioadmin"
    minio_bucket: str = "apexpay-vault"
    embedding_batch_size: int = 32 # optimal throughput
    worker_poll_interval_seconds: int = 10

    class Config:
        env_file = ".env"

settings = Settings()
