import Foundation
import SwiftData

@Model
final class CachedMemory {
    @Attribute(.unique) var id: UUID
    var content: String
    var summary: String
    var categoryRaw: String
    var source: String
    var tags: [String]
    var createdAt: Date
    var updatedAt: Date
    var needsSync: Bool

    init(id: UUID = UUID(), content: String, summary: String = "", category: MemoryCategory = .ideas, source: String = "manual", tags: [String] = [], createdAt: Date = .now, updatedAt: Date = .now, needsSync: Bool = false) {
        self.id = id
        self.content = content
        self.summary = summary
        self.categoryRaw = category.rawValue
        self.source = source
        self.tags = tags
        self.createdAt = createdAt
        self.updatedAt = updatedAt
        self.needsSync = needsSync
    }

    var category: MemoryCategory {
        MemoryCategory(rawValue: categoryRaw) ?? .ideas
    }
}

@Model
final class CachedTask {
    @Attribute(.unique) var id: UUID
    var title: String
    var notes: String
    var priorityRaw: String
    var statusRaw: String
    var dueAt: Date?
    var tags: [String]
    var createdAt: Date
    var needsSync: Bool

    init(id: UUID = UUID(), title: String, notes: String = "", priority: TaskPriority = .medium, status: TaskStatus = .inbox, dueAt: Date? = nil, tags: [String] = [], createdAt: Date = .now, needsSync: Bool = false) {
        self.id = id
        self.title = title
        self.notes = notes
        self.priorityRaw = priority.rawValue
        self.statusRaw = status.rawValue
        self.dueAt = dueAt
        self.tags = tags
        self.createdAt = createdAt
        self.needsSync = needsSync
    }

    var priority: TaskPriority { TaskPriority(rawValue: priorityRaw) ?? .medium }
    var status: TaskStatus { TaskStatus(rawValue: statusRaw) ?? .inbox }
}

@Model
final class PendingCapture {
    @Attribute(.unique) var id: UUID
    var text: String
    var source: String
    var createdAt: Date
    var retryCount: Int

    init(id: UUID = UUID(), text: String, source: String = "siri", createdAt: Date = .now, retryCount: Int = 0) {
        self.id = id
        self.text = text
        self.source = source
        self.createdAt = createdAt
        self.retryCount = retryCount
    }
}
