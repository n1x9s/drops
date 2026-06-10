import Foundation
import Observation

@MainActor
@Observable
final class SearchViewModel {
    var query = ""
    var results: [SearchResultDTO] = []
    var isLoading = false

    func search(session: SessionStore) async {
        let trimmed = query.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !trimmed.isEmpty else {
            results = []
            return
        }
        isLoading = true
        defer { isLoading = false }
        do {
            let api = SecondBrainAPI(client: APIClient(baseURL: session.apiBaseURL, accessToken: session.accessToken))
            results = try await api.search(.init(text: trimmed, limit: 25))
        } catch {
            results = []
        }
    }
}
