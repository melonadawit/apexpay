"""
Eval harness for RAG - citation precision, no hallucination guard, Amharic/English
Optimal data structure: list of test cases O(n), metrics map
"""

from typing import List, Dict
from .config import settings

test_cases = [
    {
        "query": "When is 2FA required?",
        "expected_contains": "5000",
        "expected_citation": "ONPS/10/2025",
        "lang": "en",
        "should_have_citation": True,
    },
    {
        "query": "What is the refund policy requirement?",
        "expected_contains": "refund",
        "expected_citation": "PayAtlas",
        "lang": "en",
        "should_have_citation": True,
    },
    {
        "query": "What is ETB 200k reporting?",
        "expected_contains": "200,000",
        "expected_citation": "FIC",
        "lang": "en",
        "should_have_citation": True,
    },
    {
        "query": "What is the weather today?",
        "expected_contains": "Not in compliance corpus",
        "expected_citation": "",
        "lang": "en",
        "should_have_citation": False, # should have no citation, no answer
    },
    {
        "query": "2FA መቼ ያስፈልጋል?", # Amharic: When is 2FA required?
        "expected_contains": "5000",
        "expected_citation": "ONPS",
        "lang": "am",
        "should_have_citation": True,
    },
]

def evaluate(ask_fn):
    """
    ask_fn: callable that takes query -> AskResponse
    Returns metrics dict
    """
    total = len(test_cases)
    passed = 0
    no_hallucination_pass = 0
    citation_precision = 0

    for tc in test_cases:
        resp = ask_fn(tc["query"], tc["lang"])
        answer = resp["answer"]
        citations = resp["citations"]
        no_answer = resp["no_answer"]

        # Check expected contains
        contains_ok = tc["expected_contains"].lower() in answer.lower()

        # Citation check
        has_citation = len(citations) > 0
        citation_ok = has_citation == tc["should_have_citation"]

        # No hallucination: if no_answer true, citations must be empty
        no_halluc_ok = True
        if no_answer and has_citation:
            no_halluc_ok = False

        if contains_ok and citation_ok and no_halluc_ok:
            passed += 1

        if no_halluc_ok:
            no_hallucination_pass += 1

        if has_citation and tc["expected_citation"].lower() in " ".join([c["document_id"]+c["content"] for c in citations]).lower():
            citation_precision += 1

    return {
        "total": total,
        "passed": passed,
        "accuracy": passed/total,
        "no_hallucination_rate": no_hallucination_pass/total,
        "citation_precision": citation_precision/total if total>0 else 0,
        "threshold": 0.8,
        "passed_threshold": (passed/total) >= 0.8
    }

if __name__ == "__main__":
    # Mock ask_fn for skeleton
    class MockResp:
        def __init__(self, answer, citations, no_answer):
            self.answer = answer
            self.citations = citations
            self.no_answer = no_answer

    def mock_ask(q, lang):
        if "2FA" in q or "5000" in q or "2FA መቼ" in q:
            return MockResp("Transactions above 5000 ETB require 2FA per ONPS/10/2025 [1]", [{"document_id": "rdoc_nbe_10_2025", "content": "2FA >5000"}], False)
        if "refund" in q.lower():
            return MockResp("Refund policy required per PayAtlas [2]", [{"document_id": "rdoc_refund_policy", "content": "refund"}], False)
        if "200k" in q or "200,000" in q:
            return MockResp("ETB 200,000 reporting to FIC [3]", [{"document_id": "rdoc_aml_200k", "content": "200k"}], False)
        return MockResp("Not in compliance corpus", [], True)

    metrics = evaluate(mock_ask)
    print(metrics)
    assert metrics["passed_threshold"], f"Eval failed: {metrics}"
    print("Eval passed gold")
