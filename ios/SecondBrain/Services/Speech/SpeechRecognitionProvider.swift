import Foundation

enum SpeechLanguage: String, CaseIterable, Identifiable {
    case auto
    case english = "en"
    case russian = "ru"

    var id: String { rawValue }
}

protocol SpeechRecognitionProvider: Sendable {
    func transcribe(audioFileURL: URL, language: SpeechLanguage) async throws -> String
}

enum SpeechRecognitionError: Error {
    case modelMissing
    case unsupportedFormat
    case engineUnavailable
}
