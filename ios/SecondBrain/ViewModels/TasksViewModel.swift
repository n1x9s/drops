import Foundation
import Observation

@MainActor
@Observable
final class TasksViewModel {
    var tasks: [TaskDTO] = []
    var selectedStatus: TaskStatus = .inbox
    var quickTask = ""

    var visibleTasks: [TaskDTO] {
        tasks.filter { $0.status == selectedStatus }
    }

    func load(session: SessionStore) async {
        do {
            let api = SecondBrainAPI(client: APIClient(baseURL: session.apiBaseURL, accessToken: session.accessToken))
            tasks = try await api.tasks()
        } catch {
            tasks = []
        }
    }

    func create(session: SessionStore) async {
        let text = quickTask.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !text.isEmpty else { return }
        do {
            let api = SecondBrainAPI(client: APIClient(baseURL: session.apiBaseURL, accessToken: session.accessToken))
            let task = try await api.createTask(text: text)
            tasks.insert(task, at: 0)
            quickTask = ""
        } catch {
            quickTask = text
        }
    }

    func complete(_ task: TaskDTO, session: SessionStore) async {
        do {
            let api = SecondBrainAPI(client: APIClient(baseURL: session.apiBaseURL, accessToken: session.accessToken))
            let updated = try await api.completeTask(id: task.id)
            if let index = tasks.firstIndex(where: { $0.id == task.id }) {
                tasks[index] = updated
            }
        } catch {}
    }
}
