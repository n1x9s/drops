import AppIntents
import Foundation

struct RememberThisIntent: AppIntent {
    static let title: LocalizedStringResource = "Remember This"
    static let description = IntentDescription("Save a thought to Second Brain.")
    static let openAppWhenRun = false

    @Parameter(title: "Thought")
    var text: String

    func perform() async throws -> some IntentResult & ProvidesDialog {
        let memory = try await IntentAPIClient.shared.saveMemory(text)
        return .result(dialog: IntentDialog(stringLiteral: "Remembered: \(memory.summary.isEmpty ? memory.content : memory.summary)"))
    }
}

struct FindMemoryIntent: AppIntent {
    static let title: LocalizedStringResource = "Find Memory"
    static let description = IntentDescription("Search memories and tasks in Second Brain.")
    static let openAppWhenRun = false

    @Parameter(title: "Question")
    var query: String

    func perform() async throws -> some IntentResult & ReturnsValue<String> & ProvidesDialog {
        let results = try await IntentAPIClient.shared.search(query)
        let summary = results.prefix(3).map(\.title).joined(separator: "\n")
        return .result(value: summary, dialog: IntentDialog(stringLiteral: summary.isEmpty ? "No matching memories." : summary))
    }
}
