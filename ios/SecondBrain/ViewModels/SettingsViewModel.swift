import Foundation
import Observation

@MainActor
@Observable
final class SettingsViewModel {
    var apiBaseURL = ""
    var geminiAPIKey = ""
    var telegramBotToken = ""
    var telegramChatID = ""
    var linearAPIKey = ""
    var linearTeamID = ""
    var notificationsEnabled = true
    var siriEnabled = true

    func load(session: SessionStore) {
        apiBaseURL = session.apiBaseURL.absoluteString
    }

    func save(session: SessionStore) {
        if let url = URL(string: apiBaseURL) {
            session.updateBaseURL(url)
        }
    }
}
