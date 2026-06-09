# Second Brain iOS

Native iOS source tree for the Second Brain app.

Target setup:

- iOS 26+
- Swift 6
- SwiftUI
- Observation
- SwiftData
- WidgetKit extension
- App Intents
- Local Notifications
- Background modes for sync as needed

Recommended Xcode setup:

1. Create an iOS App target named `SecondBrain`.
2. Add files under `ios/SecondBrain/SecondBrain` to the app target.
3. Add files under `ios/SecondBrain/Widgets` to a Widget Extension target.
4. Enable App Groups for sharing auth/session state between app, intents, and widget.
5. Link your whisper.cpp wrapper as the concrete engine used by `WhisperProvider`.

The app is intentionally local-first. Save/search/task flows can work against SwiftData cache while network sync and AI enrichment run opportunistically.
