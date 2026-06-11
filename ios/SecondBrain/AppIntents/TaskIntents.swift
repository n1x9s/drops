import AppIntents
import Foundation

struct CreateTaskIntent: AppIntent {
    static let title: LocalizedStringResource = "Create Task"
    static let description = IntentDescription("Create a task from natural language.")
    static let openAppWhenRun = false

    @Parameter(title: "Task")
    var text: String

    func perform() async throws -> some IntentResult & ProvidesDialog {
        let task = try await IntentAPIClient.shared.createTask(text)
        return .result(dialog: IntentDialog(stringLiteral: "Created: \(task.title)"))
    }
}

struct WhatAreMyTasksIntent: AppIntent {
    static let title: LocalizedStringResource = "What Are My Tasks"
    static let description = IntentDescription("Read active tasks from Second Brain.")
    static let openAppWhenRun = false

    func perform() async throws -> some IntentResult & ReturnsValue<String> & ProvidesDialog {
        let tasks = try await IntentAPIClient.shared.tasks().filter { $0.status != .completed }
        let summary = tasks.prefix(5).map(\.title).joined(separator: "\n")
        return .result(value: summary, dialog: IntentDialog(stringLiteral: summary.isEmpty ? "No active tasks." : summary))
    }
}

struct SecondBrainShortcuts: AppShortcutsProvider {
    static var appShortcuts: [AppShortcut] {
        AppShortcut(
            intent: RememberThisIntent(),
            phrases: [
                "Remember this in \(.applicationName)",
                "Save this thought in \(.applicationName)"
            ],
            shortTitle: "Remember",
            systemImageName: "waveform"
        )

        AppShortcut(
            intent: CreateTaskIntent(),
            phrases: [
                "Create task in \(.applicationName)",
                "Add task to \(.applicationName)"
            ],
            shortTitle: "Create Task",
            systemImageName: "checklist"
        )

        AppShortcut(
            intent: FindMemoryIntent(),
            phrases: [
                "Find memory in \(.applicationName)",
                "Search \(.applicationName)"
            ],
            shortTitle: "Find",
            systemImageName: "magnifyingglass"
        )

        AppShortcut(
            intent: WhatAreMyTasksIntent(),
            phrases: [
                "What are my tasks in \(.applicationName)"
            ],
            shortTitle: "Tasks",
            systemImageName: "tray"
        )
    }
}
