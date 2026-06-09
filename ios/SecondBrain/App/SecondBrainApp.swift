import SwiftData
import SwiftUI

@main
struct SecondBrainApp: App {
    @State private var session = SessionStore()

    var body: some Scene {
        WindowGroup {
            RootView()
                .environment(session)
                .preferredColorScheme(.dark)
        }
        .modelContainer(for: [CachedMemory.self, CachedTask.self, PendingCapture.self])
    }
}
