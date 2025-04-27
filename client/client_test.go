package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkSet(b *testing.B) {
	// Mocking the server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			b.Errorf("expected POST, got %v", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Create a new client using the mocked server URL
	client := NewClient(ts.URL)

	// Run the benchmark for Set method
	b.ResetTimer() // Reset the timer to exclude setup time
	for i := 0; i < b.N; i++ {
		err := client.Set("testKey", "testValue", 60)
		if err != nil {
			b.Fatalf("Error setting value: %v", err)
		}
	}
}

func BenchmarkGet(b *testing.B) {
	// Mocking the server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			b.Errorf("expected GET, got %v", r.Method)
		}
		w.Write([]byte(`"testValue"`))
	}))
	defer ts.Close()

	// Create a new client using the mocked server URL
	client := NewClient(ts.URL)

	// Run the benchmark for Get method
	b.ResetTimer() // Reset the timer to exclude setup time
	for i := 0; i < b.N; i++ {
		_, err := client.Get("testKey")
		if err != nil {
			b.Fatalf("Error getting value: %v", err)
		}
	}
}

func BenchmarkUpdate(b *testing.B) {
	// Mocking the server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			b.Errorf("expected POST, got %v", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Create a new client using the mocked server URL
	client := NewClient(ts.URL)

	// Run the benchmark for Update method
	b.ResetTimer() // Reset the timer to exclude setup time
	for i := 0; i < b.N; i++ {
		err := client.Update("testKey", "newValue")
		if err != nil {
			b.Fatalf("Error updating value: %v", err)
		}
	}
}

func BenchmarkRemove(b *testing.B) {
	// Mocking the server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			b.Errorf("expected DELETE, got %v", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Create a new client using the mocked server URL
	client := NewClient(ts.URL)

	// Run the benchmark for Remove method
	b.ResetTimer() // Reset the timer to exclude setup time
	for i := 0; i < b.N; i++ {
		err := client.Remove("testKey")
		if err != nil {
			b.Fatalf("Error removing value: %v", err)
		}
	}
}

func BenchmarkPush(b *testing.B) {
	// Mocking the server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			b.Errorf("expected POST, got %v", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Create a new client using the mocked server URL
	client := NewClient(ts.URL)

	// Run the benchmark for Push method
	b.ResetTimer() // Reset the timer to exclude setup time
	for i := 0; i < b.N; i++ {
		err := client.Push("testList", "newItem")
		if err != nil {
			b.Fatalf("Error pushing item: %v", err)
		}
	}
}

func BenchmarkPop(b *testing.B) {
	// Mocking the server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			b.Errorf("expected POST, got %v", r.Method)
		}
		w.Write([]byte(`"item"`))
	}))
	defer ts.Close()

	// Create a new client using the mocked server URL
	client := NewClient(ts.URL)

	// Run the benchmark for Pop method
	b.ResetTimer() // Reset the timer to exclude setup time
	for i := 0; i < b.N; i++ {
		_, err := client.Pop("testList")
		if err != nil {
			b.Fatalf("Error popping item: %v", err)
		}
	}
}
