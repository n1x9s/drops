package domain

import (
	"regexp"
	"strings"
	"time"
)

var wordRe = regexp.MustCompile(`[\p{L}\p{N}]+`)

func LocalEnrichMemory(text string) Enrichment {
	return Enrichment{
		Summary:  summarize(text),
		Category: detectCategory(text),
		Tags:     detectTags(text),
	}
}

func LocalExtractTask(text string) TaskExtraction {
	title := strings.TrimSpace(text)
	title = strings.TrimPrefix(strings.ToLower(title), "create task")
	if title == "" {
		title = strings.TrimSpace(text)
	}
	extraction := TaskExtraction{
		Title:    strings.TrimSpace(title),
		Priority: PriorityMedium,
		Tags:     detectTags(text),
	}
	lower := strings.ToLower(text)
	if strings.Contains(lower, "urgent") || strings.Contains(lower, "asap") {
		extraction.Priority = PriorityUrgent
	} else if strings.Contains(lower, "important") || strings.Contains(lower, "high") {
		extraction.Priority = PriorityHigh
	}
	if strings.Contains(lower, "tomorrow") {
		due := time.Now().Add(24 * time.Hour).Truncate(time.Minute)
		extraction.DueAt = &due
	}
	return extraction
}

func summarize(text string) string {
	text = strings.TrimSpace(text)
	runes := []rune(text)
	if len(runes) <= 180 {
		return text
	}
	return string(runes[:177]) + "..."
}

func detectCategory(text string) Category {
	lower := strings.ToLower(text)
	switch {
	case containsAny(lower, "meeting", "standup", "sync", "созвон", "встреч"):
		return CategoryMeetings
	case containsAny(lower, "learn", "study", "course", "book", "читать", "изуч"):
		return CategoryLearning
	case containsAny(lower, "project", "roadmap", "release", "ship", "проект"):
		return CategoryProjects
	case containsAny(lower, "idea", "hypothesis", "thought", "идея", "мысл"):
		return CategoryIdeas
	case containsAny(lower, "family", "home", "health", "personal", "личн"):
		return CategoryPersonal
	default:
		return CategoryWork
	}
}

func detectTags(text string) []string {
	words := wordRe.FindAllString(strings.ToLower(text), -1)
	seen := map[string]struct{}{}
	tags := make([]string, 0, 6)
	stop := map[string]struct{}{
		"the": {}, "and": {}, "for": {}, "with": {}, "this": {}, "that": {}, "about": {},
		"это": {}, "как": {}, "что": {}, "для": {}, "надо": {}, "нужно": {}, "про": {},
	}
	for _, w := range words {
		if len([]rune(w)) < 4 {
			continue
		}
		if _, ok := stop[w]; ok {
			continue
		}
		if _, ok := seen[w]; ok {
			continue
		}
		seen[w] = struct{}{}
		tags = append(tags, w)
		if len(tags) == 6 {
			break
		}
	}
	return tags
}

func containsAny(value string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(value, needle) {
			return true
		}
	}
	return false
}
