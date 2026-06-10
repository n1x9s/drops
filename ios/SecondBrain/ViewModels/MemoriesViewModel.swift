import Foundation
import Observation

@MainActor
@Observable
final class MemoriesViewModel {
    var query = ""
    var selectedCategory: MemoryCategory?
    var memories: [MemoryDTO] = []
    var isSearching = false

    func load(session: SessionStore) async {
        do {
            let api = SecondBrainAPI(client: APIClient(baseURL: session.apiBaseURL, accessToken: session.accessToken))
            memories = try await api.memories(query: query, category: selectedCategory)
        } catch {
            memories = []
        }
    }

    func categoryBinding(_ category: MemoryCategory) -> Bool {
        selectedCategory == category
    }
}
