import SwiftUI

struct TasksView: View {
    @Environment(SessionStore.self) private var session
    @State private var viewModel = TasksViewModel()

    var body: some View {
        NavigationStack {
            VStack(spacing: 14) {
                quickCreate
                statusPicker
                List(viewModel.visibleTasks) { task in
                    TaskRow(task: task) {
                        Task { await viewModel.complete(task, session: session) }
                    }
                    .listRowBackground(Color.clear)
                }
                .scrollContentBackground(.hidden)
            }
            .padding(.top, 8)
            .background(SecondBrainTheme.background.ignoresSafeArea())
            .navigationTitle("Tasks")
            .task { await viewModel.load(session: session) }
        }
    }

    private var quickCreate: some View {
        HStack(spacing: 10) {
            TextField("Create task...", text: $viewModel.quickTask)
                .textFieldStyle(.roundedBorder)
            Button {
                Task { await viewModel.create(session: session) }
            } label: {
                Image(systemName: "plus")
            }
            .buttonStyle(.borderedProminent)
            .tint(SecondBrainTheme.accent)
            .accessibilityLabel("Create task")
        }
        .padding(.horizontal, 20)
    }

    private var statusPicker: some View {
        Picker("Status", selection: $viewModel.selectedStatus) {
            ForEach(TaskStatus.allCases) { status in
                Text(status.rawValue.capitalized).tag(status)
            }
        }
        .pickerStyle(.segmented)
        .padding(.horizontal, 20)
    }
}
