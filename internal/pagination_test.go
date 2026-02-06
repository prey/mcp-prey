package internal

import (
	"net/url"
	"testing"
)

func TestNormalizePage(t *testing.T) {
	if got := NormalizePage(0); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
	if got := NormalizePage(5); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestNormalizePageSize(t *testing.T) {
	if got, err := NormalizePageSize(0); err != nil || got != 20 {
		t.Fatalf("expected default 20, got %d err=%v", got, err)
	}
	if _, err := NormalizePageSize(101); err == nil {
		t.Fatalf("expected error for page_size > 100")
	}
	if got, err := NormalizePageSize(100); err != nil || got != 100 {
		t.Fatalf("expected 100, got %d err=%v", got, err)
	}
}

func TestAddPagination(t *testing.T) {
	q, err := AddPagination(url.Values{}, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.Get("page") != "1" {
		t.Fatalf("expected page=1, got %s", q.Get("page"))
	}
	if q.Get("page_size") != "20" {
		t.Fatalf("expected page_size=20, got %s", q.Get("page_size"))
	}
}

func TestMeta(t *testing.T) {
	meta, err := Meta(2, 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta["page"].(int) != 2 {
		t.Fatalf("expected page=2")
	}
	if meta["page_size"].(int) != 50 {
		t.Fatalf("expected page_size=50")
	}
}
