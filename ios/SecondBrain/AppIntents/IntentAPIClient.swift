import AppIntents
import Foundation

actor IntentAPIClient {
    static let shared = IntentAPIClient()

    private let defaults = UserDefaults(suiteName: "group.secondbrain.app") ?? .standard
    private let encoder = JSONEncoder()
    private let decoder = JSONDecoder()

    init() {
        encoder.dateEncodingStrategy = .iso8601
        decoder.dateDecodingStrategy = .iso8601
    }

    func saveMemory(_ text: String) async throws -> MemoryDTO {
        struct Body: Encodable {
            let content: String
            let source: String
        }
        return try await send(path: "/memories", method: "POST", body: Body(content: text, source: "siri"))
    }

    func createTask(_ text: String) async throws -> TaskDTO {
        struct Body: Encodable {
            let text: String
        }
        return try await send(path: "/tasks", method: "POST", body: Body(text: text))
    }

    func search(_ text: String) async throws -> [SearchResultDTO] {
        let query = "?q=\(text.addingPercentEncoding(withAllowedCharacters: .urlQueryAllowed) ?? text)&limit=5"
        return try await send(path: "/search\(query)", method: "GET", body: Optional<Data>.none)
    }

    func tasks() async throws -> [TaskDTO] {
        try await send(path: "/tasks", method: "GET", body: Optional<Data>.none)
    }

    private func send<Request: Encodable, Response: Decodable>(path: String, method: String, body: Request?) async throws -> Response {
        let data = try body.map { try encoder.encode($0) }
        return try await send(path: path, method: method, body: data)
    }

    private func send<Response: Decodable>(path: String, method: String, body: Data?) async throws -> Response {
        let base = defaults.string(forKey: "apiBaseURL") ?? "http://localhost:8080/api/v1"
        guard let url = URL(string: base + path) else {
            throw URLError(.badURL)
        }
        var request = URLRequest(url: url)
        request.httpMethod = method
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        if let token = defaults.string(forKey: "accessToken") {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }
        request.httpBody = body
        let (data, _) = try await URLSession.shared.data(for: request)
        let envelope = try decoder.decode(APIClient.Envelope<Response>.self, from: data)
        if let value = envelope.data {
            return value
        }
        throw envelope.error ?? URLError(.badServerResponse)
    }
}
