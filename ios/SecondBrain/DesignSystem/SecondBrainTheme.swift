import SwiftUI

enum SecondBrainTheme {
    static let background = Color(red: 0.035, green: 0.039, blue: 0.047)
    static let elevated = Color(red: 0.075, green: 0.082, blue: 0.096)
    static let surface = Color(red: 0.105, green: 0.114, blue: 0.132)
    static let primary = Color(red: 0.79, green: 0.88, blue: 1.0)
    static let accent = Color(red: 0.44, green: 0.72, blue: 0.64)
    static let warning = Color(red: 0.94, green: 0.66, blue: 0.34)
    static let text = Color(red: 0.93, green: 0.94, blue: 0.96)
    static let secondaryText = Color(red: 0.62, green: 0.66, blue: 0.72)
    static let hairline = Color.white.opacity(0.08)

    static let spring = Animation.spring(response: 0.36, dampingFraction: 0.82)
}

struct QuietPanel<Content: View>: View {
    let content: Content

    init(@ViewBuilder content: () -> Content) {
        self.content = content()
    }

    var body: some View {
        content
            .padding(16)
            .background(SecondBrainTheme.elevated, in: RoundedRectangle(cornerRadius: 8, style: .continuous))
            .overlay(
                RoundedRectangle(cornerRadius: 8, style: .continuous)
                    .stroke(SecondBrainTheme.hairline, lineWidth: 1)
            )
    }
}
