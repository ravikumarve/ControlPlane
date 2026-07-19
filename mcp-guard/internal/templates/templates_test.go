package templates

import "testing"

func TestGet(t *testing.T) {
	content, ok := Get("github")
	if !ok {
		t.Fatal("Get(github) returned false")
	}
	if len(content) == 0 {
		t.Fatal("Get(github) returned empty content")
	}
}

func TestGetCaseInsensitive(t *testing.T) {
	_, ok := Get("GitHub")
	if !ok {
		t.Fatal("Get(GitHub) should be case insensitive")
	}
}

func TestGetUnknown(t *testing.T) {
	_, ok := Get("nonexistent")
	if ok {
		t.Fatal("Get(nonexistent) should return false")
	}
}

func TestList(t *testing.T) {
	list := List()
	if len(list) != 5 {
		t.Fatalf("expected 5 templates, got %d: %v", len(list), list)
	}
}

func TestAllTemplatesNonEmpty(t *testing.T) {
	for _, name := range List() {
		content, ok := Get(name)
		if !ok {
			t.Fatalf("Get(%q) returned false", name)
		}
		if len(content) == 0 {
			t.Fatalf("Get(%q) returned empty content", name)
		}
	}
}
