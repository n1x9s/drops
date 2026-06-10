import Foundation

actor APIClient {
    struct Envelope<T: Decodable>: Decodable {
        let success: Bool
        let data: T?
        let error: APIErrorPayload?
    }

    struct APIErrorPayload: Decodable, Error {
        let code: String
        let message: String
    }

    private let session: URLSession
    private let baseURL: URL
    private let accessToken: String?
    private let decoder: JSONDecoder
    private let encoder: JSONEncoder

    init(baseURL: URL, accessToken: String?, session: URLSession = .shared) {
        self.baseURL = baseURL
        self.accessToken = accessToken
        self.session = session
        decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        encoder = JSONEncoder()
        encoder.dateEncodingStrategy = .iso8601
    }

    func get<T: Decodable>(_ path: String, query: [URLQueryItem] = []) async throws -> T {
        try await send(path: path, method: "GET", query: query, body: Optional<Data>.none)
    }

    func post<Request: Encodable, Response: Decodable>(_ path: String, body: Request) async throws -> Response {
        let data = try encoder.encode(body)
        return try await send(path: path, method: "POST", body: data)
    }

    func patch<Request: Encodable, Response: Decodable>(_ path: String, body: Request) async throws -> Response {
        let data = try encoder.encode(body)
        return try await send(path: path, method: "PATCH", body: data)
    }

    func put<Request: Encodable, Response: Decodable>(_ path: String, body: Request) async throws -> Response {
        let data = try encoder.encode(body)
        return try await send(path: path, method: "PUT", body: data)
    }

    func delete<Response: Decodable>(_ path: String) async throws -> Response {
        try await send(path: path, method: "DELETE", body: Optional<Data>.none)
    }

    private func send<T: Decodable>(path: String, method: String, query: [URLQueryItem] = [], body: Data?) async throws -> T {
        var components = URLComponents(url: baseURL.appending(path: path), resolvingAgainstBaseURL: false)!
        if !query.isEmpty {
            components.queryItems = query
        }
        var request = URLRequest(url: components.url!)
        request.httpMethod = method
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        if let token = accessToken {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }
        request.httpBody = body

        let (data, response) = try await session.data(for: request)
        guard let http = response as? HTTPURLResponse else {
            throw URLError(.badServerResponse)
        }
        let envelope = try decoder.decode(Envelope<T>.self, from: data)
        if (200..<300).contains(http.statusCode), let value = envelope.data {
            return value
        }
        throw envelope.error ?? APIErrorPayload(code: "http_error", message: "HTTP \(http.statusCode)")
    }
}
