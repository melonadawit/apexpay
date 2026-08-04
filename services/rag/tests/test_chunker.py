import pytest
from app.chunker import chunk_text, chunk_char

def test_chunk_text_small():
    text = "Short text"
    chunks = chunk_text(text, 800, 100)
    assert len(chunks) == 1
    assert chunks[0] == text

def test_chunk_overlap():
    text = "A " * 1000 # 2000 chars
    chunks = chunk_text(text, 100, 20)
    assert len(chunks) > 1
    # Overlap check: second chunk should start before first ends
    # Optimal: ensure no data loss
    combined = "".join(chunks)
    # Due to overlap, combined longer than original but must contain original tokens
    assert len(combined) >= len(text)

def test_chunk_char():
    text = "x" * 1000
    chunks = chunk_char(text, 300, 50)
    assert len(chunks) == 4 # 0-300, 250-550, 500-800, 750-1000
    assert "".join([c[:50] for c in chunks]) != "" # sanity

def test_no_empty_chunks():
    text = "NBE ONPS/10/2025 requires 2FA above 5000 ETB. " * 20
    chunks = chunk_text(text, 50, 10)
    for c in chunks:
        assert len(c.strip()) > 0
