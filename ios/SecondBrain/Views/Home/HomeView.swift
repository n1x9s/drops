import SwiftUI

struct HomeView: View {
    @Environment(SessionStore.self) private var session
    @State private var viewModel = HomeViewModel()
    @FocusState private var isCaptureFocused: Bool

    var body: some View {
        NavigationStack {
            ScrollView {
                VStack(alignment: .leading, spacing: 18) {
                    capture
                    today
                    recent
                    activeTasks
                }
                .padding(20)
            }
            .background(SecondBrainTheme.background.ignoresSafeArea())
            .navigationTitle("Second Brain")
            .toolbar {
                Button {
                    isCaptureFocused = true
                } label: {
                    Image(systemName: "waveform")
                }
                .accessibilityLabel("Capture")
            }
            .task { await viewModel.load(session: session) }
        }
    }

    private var capture: some View {
        QuietPanel {
            VStack(alignment: .leading, spacing: 12) {
                Text("Capture")
                    .font(.headline)
                    .foregroundStyle(SecondBrainTheme.text)

                TextField("Remember this...", text: $viewModel.captureText, axis: .vertical)
                    .textFieldStyle(.plain)
                    .lineLimit(2...5)
                    .focused($isCaptureFocused)
                    .foregroundStyle(SecondBrainTheme.text)

                HStack {
                    Button {
                        Task { await viewModel.saveMemory(session: session) }
                    } label: {
                        Label("Save", systemImage: "tray.and.arrow.down")
                    }
                    .buttonStyle(.borderedProminent)
                    .tint(SecondBrainTheme.accent)

                    Spacer()

                    if let error = viewModel.errorMessage {
                        Text(error)
                            .font(.caption)
                            .foregroundStyle(SecondBrainTheme.secondaryText)
                    }
                }
            }
        }
    }

    private var today: some View {
        HStack(spacing: 12) {
            metric("Today", "\(viewModel.activeTasks.filter { $0.status == .today }.count)", "calendar")
            metric("Inbox", "\(viewModel.activeTasks.filter { $0.status == .inbox }.count)", "tray")
            metric("Recent", "\(viewModel.recentMemories.count)", "clock")
        }
    }

    private func metric(_ title: String, _ value: String, _ symbol: String) -> some View {
        QuietPanel {
            VStack(alignment: .leading, spacing: 10) {
                Image(systemName: symbol)
                    .foregroundStyle(SecondBrainTheme.primary)
                Text(value)
                    .font(.system(size: 26, weight: .semibold, design: .rounded))
                    .foregroundStyle(SecondBrainTheme.text)
                Text(title)
                    .font(.caption)
                    .foregroundStyle(SecondBrainTheme.secondaryText)
            }
            .frame(maxWidth: .infinity, alignment: .leading)
        }
    }

    private var recent: some View {
        VStack(alignment: .leading, spacing: 10) {
            Text("Recent Memories")
                .font(.headline)
                .foregroundStyle(SecondBrainTheme.text)
            ForEach(viewModel.recentMemories) { memory in
                MemoryRow(memory: memory)
            }
        }
    }

    private var activeTasks: some View {
        VStack(alignment: .leading, spacing: 10) {
            Text("Active Tasks")
                .font(.headline)
                .foregroundStyle(SecondBrainTheme.text)
            ForEach(viewModel.activeTasks.prefix(5)) { task in
                TaskRow(task: task) {
                    Task { await viewModel.complete(task, session: session) }
                }
            }
        }
    }
}
