import Foundation

enum MemoryCategory: String, Codable, CaseIterable, Identifiable {
    case work = "Work"
    case learning = "Learning"
    case personal = "Personal"
    case projects = "Projects"
    case meetings = "Meetings"
    case ideas = "Ideas"

    var id: String { rawValue }
}

enum TaskPriority: String, Codable, CaseIterable, Identifiable {
    case low
    case medium
    case high
    case urgent

    var id: String { rawValue }
}

enum TaskStatus: String, Codable, CaseIterable, Identifiable {
    case inbox
    case today
    case upcoming
    case overdue
    case completed

    var id: String { rawValue }
}

struct TagDTO: Codable, Hashable, Identifiable {
    let id: UUID
    let name: String
}

struct MemoryDTO: Codable, Identifiable, Hashable {
    let id: UUID
    let content: String
    let summary: String
    let category: MemoryCategory
    let source: String
    let tags: [TagDTO]
    let createdAt: Date

    enum CodingKeys: String, CodingKey {
        case id, content, summary, category, source, tags
        case createdAt = "created_at"
    }
}

struct TaskDTO: Codable, Identifiable, Hashable {
    let id: UUID
    let title: String
    let notes: String
    let priority: TaskPriority
    let status: TaskStatus
    let dueAt: Date?
    let tags: [TagDTO]

    enum CodingKeys: String, CodingKey {
        case id, title, notes, priority, status, tags
        case dueAt = "due_at"
    }
}

struct SearchResultDTO: Codable, Identifiable, Hashable {
    let id: UUID
    let type: String
    let title: String
    let snippet: String
    let score: Double
    let category: String
}
