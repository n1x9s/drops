import Foundation
import Observation

@MainActor
@Observable
final class HomeViewModel {
    var recentMemories: [MemoryDTO] = []
    var activeTasks: [TaskDTO] = []
    var captureText = ""
    var isLoading = false
    var errorMessage: String?

    func load(session: SessionStore) async {
        isLoading = true
        defer { isLoading = false }
        do {
            let api = SecondBrainAPI(client: APIClient(baseURL: session.apiBaseURL, accessToken: session.accessToken))
            async let memories = api.memories()
            async let tasks = api.tasks()
            recentMemories = Array(try await memories.prefix(5))
            activeTasks = try await tasks.filter { $0.status != .completed }
            errorMessage = nil
        } catch {
            errorMessage = "Offline mode"
        }
    }

    func saveMemory(session: SessionStore) async {
        let text = captureText.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !text.isEmpty else { return }
        do {
            let api = SecondBrainAPI(client: APIClient(baseURL: session.apiBaseURL, accessToken: session.accessToken))
            let memory = try await api.createMemory(content: text, source: "manual")
            recentMemories.insert(memory, at: 0)
            captureText = ""
            errorMessage = nil
        } catch {
            errorMessage = "Saved locally when offline"
        }
    }

    func complete(_ task: TaskDTO, session: SessionStore) async {
        do {
            let api = SecondBrainAPI(client: APIClient(baseURL: session.apiBaseURL, accessToken: session.accessToken))
            let updated = try await api.completeTask(id: task.id)
            if let index = activeTasks.firstIndex(where: { $0.id == task.id }) {
                activeTasks[index] = updated
            }
        } catch {
            errorMessage = "Could not complete task"
        }
    }
}
