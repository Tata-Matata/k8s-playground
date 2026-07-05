package main

import (
	"net/http/httptest"
	"os"
	"testing"
)

func TestResolveNamespaceFromQuery(t *testing.T) {
	t.Setenv("POD_NAMESPACE", "env-namespace")
	req := httptest.NewRequest("GET", "/pods?namespace=query-namespace", nil)

	if got := resolveNamespace(req); got != "query-namespace" {
		t.Fatalf("resolveNamespace() = %q, want %q", got, "query-namespace")
	}
}

func TestResolveNamespaceFromEnv(t *testing.T) {
	t.Setenv("POD_NAMESPACE", "env-namespace")
	req := httptest.NewRequest("GET", "/pods", nil)

	if got := resolveNamespace(req); got != "env-namespace" {
		t.Fatalf("resolveNamespace() = %q, want %q", got, "env-namespace")
	}
}

func TestResolveNamespaceDefault(t *testing.T) {
	_ = os.Unsetenv("POD_NAMESPACE")
	req := httptest.NewRequest("GET", "/pods", nil)

	if got := resolveNamespace(req); got != "default" {
		t.Fatalf("resolveNamespace() = %q, want %q", got, "default")
	}
}
