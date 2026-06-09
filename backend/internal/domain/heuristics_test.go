package domain

import "testing"

func TestLocalEnrichMemoryDetectsLearningAndTags(t *testing.T) {
	enrichment := LocalEnrichMemory("Study Go runtime netpoll and Kafka consumer groups before Sunday")

	if enrichment.Category != CategoryLearning {
		t.Fatalf("category = %s, want %s", enrichment.Category, CategoryLearning)
	}
	if enrichment.Summary == "" {
		t.Fatal("summary should not be empty")
	}
	if len(enrichment.Tags) == 0 {
		t.Fatal("expected tags")
	}
}

func TestLocalExtractTaskDetectsUrgency(t *testing.T) {
	task := LocalExtractTask("Create task urgent implement Telegram authentication tomorrow")

	if task.Priority != PriorityUrgent {
		t.Fatalf("priority = %s, want %s", task.Priority, PriorityUrgent)
	}
	if task.DueAt == nil {
		t.Fatal("expected due date for tomorrow")
	}
}
