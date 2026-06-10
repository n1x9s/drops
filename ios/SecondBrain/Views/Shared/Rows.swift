import SwiftUI

struct MemoryRow: View {
    let memory: MemoryDTO

    var body: some View {
        QuietPanel {
            VStack(alignment: .leading, spacing: 8) {
                HStack {
                    Text(memory.category.rawValue)
                        .font(.caption.weight(.medium))
                        .foregroundStyle(SecondBrainTheme.accent)
                    Spacer()
                    Text(memory.createdAt, style: .date)
                        .font(.caption)
                        .foregroundStyle(SecondBrainTheme.secondaryText)
                }
                Text(memory.summary.isEmpty ? memory.content : memory.summary)
                    .font(.body)
                    .foregroundStyle(SecondBrainTheme.text)
                    .lineLimit(3)
                if !memory.tags.isEmpty {
                    HStack {
                        ForEach(memory.tags.prefix(4), id: \.self) { tag in
                            Text("#\(tag.name)")
                                .font(.caption)
                                .foregroundStyle(SecondBrainTheme.secondaryText)
                        }
                    }
                }
            }
        }
    }
}

struct TaskRow: View {
    let task: TaskDTO
    let onComplete: () -> Void

    var body: some View {
        QuietPanel {
            HStack(alignment: .top, spacing: 12) {
                Button(action: onComplete) {
                    Image(systemName: task.status == .completed ? "checkmark.circle.fill" : "circle")
                }
                .buttonStyle(.plain)
                .foregroundStyle(task.status == .completed ? SecondBrainTheme.accent : SecondBrainTheme.secondaryText)
                .accessibilityLabel("Complete task")

                VStack(alignment: .leading, spacing: 6) {
                    Text(task.title)
                        .font(.body.weight(.medium))
                        .foregroundStyle(SecondBrainTheme.text)
                    if let due = task.dueAt {
                        Text(due, style: .date)
                            .font(.caption)
                            .foregroundStyle(SecondBrainTheme.warning)
                    }
                }
                Spacer()
                Text(task.priority.rawValue)
                    .font(.caption.weight(.semibold))
                    .foregroundStyle(SecondBrainTheme.secondaryText)
            }
        }
    }
}
