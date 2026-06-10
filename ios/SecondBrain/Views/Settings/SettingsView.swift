import SwiftUI

struct SettingsView: View {
    @Environment(SessionStore.self) private var session
    @State private var viewModel = SettingsViewModel()

    var body: some View {
        NavigationStack {
            Form {
                Section("API") {
                    TextField("Base URL", text: $viewModel.apiBaseURL)
                        .textInputAutocapitalization(.never)
                        .keyboardType(.URL)
                    Button {
                        viewModel.save(session: session)
                    } label: {
                        Label("Save", systemImage: "checkmark")
                    }
                }

                Section("Gemini") {
                    SecureField("API key", text: $viewModel.geminiAPIKey)
                }

                Section("Telegram") {
                    SecureField("Bot token", text: $viewModel.telegramBotToken)
                    TextField("Chat ID", text: $viewModel.telegramChatID)
                }

                Section("Linear") {
                    SecureField("API key", text: $viewModel.linearAPIKey)
                    TextField("Team ID", text: $viewModel.linearTeamID)
                }

                Section("System") {
                    Toggle("Siri", isOn: $viewModel.siriEnabled)
                    Toggle("Notifications", isOn: $viewModel.notificationsEnabled)
                }
            }
            .scrollContentBackground(.hidden)
            .background(SecondBrainTheme.background.ignoresSafeArea())
            .navigationTitle("Settings")
            .task { viewModel.load(session: session) }
        }
    }
}
