import Foundation
import Observation

@MainActor
@Observable
final class SessionStore {
    var accessToken: String?
    var refreshToken: String?
    var apiBaseURL = URL(string: "http://localhost:8080/api/v1")!
    var isAuthenticated: Bool { accessToken?.isEmpty == false }

    private let defaults = UserDefaults(suiteName: "group.secondbrain.app") ?? .standard

    init() {
        accessToken = defaults.string(forKey: "accessToken")
        refreshToken = defaults.string(forKey: "refreshToken")
        if let raw = defaults.string(forKey: "apiBaseURL"), let url = URL(string: raw) {
            apiBaseURL = url
        }
    }

    func update(accessToken: String, refreshToken: String) {
        self.accessToken = accessToken
        self.refreshToken = refreshToken
        defaults.set(accessToken, forKey: "accessToken")
        defaults.set(refreshToken, forKey: "refreshToken")
    }

    func updateBaseURL(_ url: URL) {
        apiBaseURL = url
        defaults.set(url.absoluteString, forKey: "apiBaseURL")
    }

    func clear() {
        accessToken = nil
        refreshToken = nil
        defaults.removeObject(forKey: "accessToken")
        defaults.removeObject(forKey: "refreshToken")
    }

    func bootstrapDevelopmentSession() async {
        guard accessToken == nil else { return }
        let baseURL = apiBaseURL
        let email = "dev@secondbrain.local"
        let password = "password123"
        do {
            let session = try await Self.authenticate(baseURL: baseURL, path: "/auth/register", body: AuthRequest(email: email, name: "Developer", password: password))
            update(accessToken: session.accessToken, refreshToken: session.refreshToken)
        } catch {
            do {
                let session = try await Self.authenticate(baseURL: baseURL, path: "/auth/login", body: LoginRequest(email: email, password: password))
                update(accessToken: session.accessToken, refreshToken: session.refreshToken)
            } catch {
                // The app remains local-first when the backend is unavailable.
            }
        }
    }

    private static func authenticate<Request: Encodable>(baseURL: URL, path: String, body: Request) async throws -> AuthSession {
        var request = URLRequest(url: baseURL.appending(path: path))
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.httpBody = try JSONEncoder().encode(body)
        let (data, _) = try await URLSession.shared.data(for: request)
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        let envelope = try decoder.decode(AuthEnvelope.self, from: data)
        guard envelope.success, let session = envelope.data else {
            throw URLError(.userAuthenticationRequired)
        }
        return session
    }
}

private struct AuthRequest: Encodable {
    let email: String
    let name: String
    let password: String
}

private struct LoginRequest: Encodable {
    let email: String
    let password: String
}

private struct AuthEnvelope: Decodable {
    let success: Bool
    let data: AuthSession?
}

private struct AuthSession: Decodable {
    let accessToken: String
    let refreshToken: String

    enum CodingKeys: String, CodingKey {
        case accessToken = "access_token"
        case refreshToken = "refresh_token"
    }
}
