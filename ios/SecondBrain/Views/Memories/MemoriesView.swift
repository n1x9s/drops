import SwiftUI

struct MemoriesView: View {
    @Environment(SessionStore.self) private var session
    @State private var viewModel = MemoriesViewModel()

    var body: some View {
        NavigationStack {
            VStack(spacing: 14) {
                categoryStrip
                List(viewModel.memories) { memory in
                    MemoryRow(memory: memory)
                        .listRowBackground(Color.clear)
                }
                .scrollContentBackground(.hidden)
            }
            .padding(.top, 8)
            .background(SecondBrainTheme.background.ignoresSafeArea())
            .navigationTitle("Memories")
            .searchable(text: $viewModel.query, prompt: "Search memories")
            .onSubmit(of: .search) { Task { await viewModel.load(session: session) } }
            .task { await viewModel.load(session: session) }
        }
    }

    private var categoryStrip: some View {
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(spacing: 8) {
                Button("All") {
                    viewModel.selectedCategory = nil
                    Task { await viewModel.load(session: session) }
                }
                .buttonStyle(.bordered)

                ForEach(MemoryCategory.allCases) { category in
                    Button(category.rawValue) {
                        viewModel.selectedCategory = category
                        Task { await viewModel.load(session: session) }
                    }
                    .buttonStyle(.borderedProminent)
                    .tint(viewModel.selectedCategory == category ? SecondBrainTheme.accent : SecondBrainTheme.surface)
                }
            }
            .padding(.horizontal, 20)
        }
    }
}
