import SwiftUI

struct SearchView: View {
    @Environment(SessionStore.self) private var session
    @State private var viewModel = SearchViewModel()

    var body: some View {
        NavigationStack {
            List(viewModel.results) { result in
                QuietPanel {
                    VStack(alignment: .leading, spacing: 8) {
                        HStack {
                            Text(result.title)
                                .font(.headline)
                                .foregroundStyle(SecondBrainTheme.text)
                            Spacer()
                            Text(result.type)
                                .font(.caption)
                                .foregroundStyle(SecondBrainTheme.secondaryText)
                        }
                        Text(result.snippet)
                            .font(.subheadline)
                            .foregroundStyle(SecondBrainTheme.secondaryText)
                            .lineLimit(3)
                        ProgressView(value: min(max(result.score, 0), 1))
                            .tint(SecondBrainTheme.accent)
                    }
                }
                .listRowBackground(Color.clear)
            }
            .scrollContentBackground(.hidden)
            .background(SecondBrainTheme.background.ignoresSafeArea())
            .navigationTitle("Search")
            .searchable(text: $viewModel.query, prompt: "What did I say about Kafka?")
            .onSubmit(of: .search) { Task { await viewModel.search(session: session) } }
        }
    }
}
