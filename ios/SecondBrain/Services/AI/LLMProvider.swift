import Foundation

struct MemoryEnrichment: Codable, Sendable {
    let summary: String
    let category: MemoryCategory
    let tags: [String]
}

protocol LLMProvider: Sendable {
    func enrichMemory(_ text: String) async throws -> MemoryEnrichment
}

struct LocalMemoryEnricher {
    func enrich(_ text: String) -> MemoryEnrichment {
        let lower = text.lowercased()
        let category: MemoryCategory
        if lower.contains("meeting") || lower.contains("созвон") {
            category = .meetings
        } else if lower.contains("study") || lower.contains("learn") || lower.contains("изуч") {
            category = .learning
        } else if lower.contains("idea") || lower.contains("идея") {
            category = .ideas
        } else {
            category = .work
        }
        let words = lower
            .split { !$0.isLetter && !$0.isNumber }
            .map(String.init)
            .filter { $0.count > 3 }
        return MemoryEnrichment(summary: String(text.prefix(180)), category: category, tags: Array(Set(words)).prefix(6).map(\.self))
    }
}
