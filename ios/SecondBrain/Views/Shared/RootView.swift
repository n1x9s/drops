import SwiftUI

struct RootView: View {
    @Environment(SessionStore.self) private var session

    var body: some View {
        TabView {
            HomeView()
                .tabItem { Label("Home", systemImage: "circle.grid.2x2") }

            MemoriesView()
                .tabItem { Label("Memories", systemImage: "text.book.closed") }

            TasksView()
                .tabItem { Label("Tasks", systemImage: "checklist") }

            SearchView()
                .tabItem { Label("Search", systemImage: "magnifyingglass") }

            SettingsView()
                .tabItem { Label("Settings", systemImage: "gearshape") }
        }
        .tint(SecondBrainTheme.primary)
        .background(SecondBrainTheme.background)
        .task {
            await session.bootstrapDevelopmentSession()
        }
    }
}
