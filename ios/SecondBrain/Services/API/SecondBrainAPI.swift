import Foundation

struct CreateMemoryBody: Encodable {
    let content: String
    let source: String
}

struct CreateTaskBody: Encodable {
    let text: String?
    let title: String?
    let notes: String?
    let dueAt: Date?
    let priority: String?
    let tags: [String]?

    enum CodingKeys: String, CodingKey {
        case text, title, notes, priority, tags
        case dueAt = "due_at"
    }
}

struct SearchQuery {
    let text: String
    let limit: Int
}

actor SecondBrainAPI {
    private let client: APIClient

    init(client: APIClient) {
        self.client = client
    }

    func createMemory(content: String, source: String) async throws -> MemoryDTO {
        try await client.post("/memories", body: CreateMemoryBody(content: content, source: source))
    }

    func memories(query: String? = nil, category: MemoryCategory? = nil, tag: String? = nil) async throws -> [MemoryDTO] {
        var items: [URLQueryItem] = []
        if let query, !query.isEmpty { items.append(.init(name: "q", value: query)) }
        if let category { items.append(.init(name: "category", value: category.rawValue)) }
        if let tag, !tag.isEmpty { items.append(.init(name: "tag", value: tag)) }
        return try await client.get("/memories", query: items)
    }

    func createTask(text: String) async throws -> TaskDTO {
        try await client.post("/tasks", body: CreateTaskBody(text: text, title: nil, notes: nil, dueAt: nil, priority: nil, tags: nil))
    }

    func tasks(status: TaskStatus? = nil) async throws -> [TaskDTO] {
        let query = status.map { [URLQueryItem(name: "status", value: $0.rawValue)] } ?? []
        return try await client.get("/tasks", query: query)
    }

    func completeTask(id: UUID) async throws -> TaskDTO {
        struct Empty: Encodable {}
        return try await client.post("/tasks/\(id.uuidString)/complete", body: Empty())
    }

    func search(_ query: SearchQuery) async throws -> [SearchResultDTO] {
        try await client.get("/search", query: [
            .init(name: "q", value: query.text),
            .init(name: "limit", value: "\(query.limit)")
        ])
    }
}
