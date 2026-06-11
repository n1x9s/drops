import SwiftUI
import WidgetKit

struct SecondBrainEntry: TimelineEntry {
    let date: Date
    let todayCount: Int
    let inboxCount: Int
    let latestMemory: String
}

struct SecondBrainTimelineProvider: TimelineProvider {
    func placeholder(in context: Context) -> SecondBrainEntry {
        SecondBrainEntry(date: .now, todayCount: 3, inboxCount: 7, latestMemory: "Kafka rebalance diagram")
    }

    func getSnapshot(in context: Context, completion: @escaping (SecondBrainEntry) -> Void) {
        completion(readEntry())
    }

    func getTimeline(in context: Context, completion: @escaping (Timeline<SecondBrainEntry>) -> Void) {
        completion(Timeline(entries: [readEntry()], policy: .after(.now.addingTimeInterval(900))))
    }

    private func readEntry() -> SecondBrainEntry {
        let defaults = UserDefaults(suiteName: "group.secondbrain.app") ?? .standard
        return SecondBrainEntry(
            date: .now,
            todayCount: defaults.integer(forKey: "widget.todayCount"),
            inboxCount: defaults.integer(forKey: "widget.inboxCount"),
            latestMemory: defaults.string(forKey: "widget.latestMemory") ?? "Ready to capture"
        )
    }
}

struct SecondBrainWidgetView: View {
    let entry: SecondBrainEntry

    var body: some View {
        ZStack {
            SecondBrainTheme.background
            VStack(alignment: .leading, spacing: 10) {
                HStack {
                    Image(systemName: "brain.head.profile")
                    Text("Second Brain")
                        .font(.headline)
                    Spacer()
                }
                .foregroundStyle(SecondBrainTheme.text)

                Text(entry.latestMemory)
                    .font(.subheadline)
                    .foregroundStyle(SecondBrainTheme.secondaryText)
                    .lineLimit(2)

                HStack {
                    Label("\(entry.todayCount)", systemImage: "calendar")
                    Label("\(entry.inboxCount)", systemImage: "tray")
                }
                .font(.caption.weight(.semibold))
                .foregroundStyle(SecondBrainTheme.accent)
            }
            .padding()
        }
    }
}

struct SecondBrainWidget: Widget {
    let kind = "SecondBrainWidget"

    var body: some WidgetConfiguration {
        StaticConfiguration(kind: kind, provider: SecondBrainTimelineProvider()) { entry in
            SecondBrainWidgetView(entry: entry)
        }
        .configurationDisplayName("Second Brain")
        .description("Today’s tasks and latest memory.")
        .supportedFamilies([.systemSmall, .systemMedium])
    }
}
