import Foundation

struct WhisperProvider: SpeechRecognitionProvider {
    let modelURL: URL?

    init(modelURL: URL? = Bundle.main.url(forResource: "ggml-medium", withExtension: "bin")) {
        self.modelURL = modelURL
    }

    func transcribe(audioFileURL: URL, language: SpeechLanguage) async throws -> String {
        guard modelURL != nil else {
            throw SpeechRecognitionError.modelMissing
        }

        // Production builds should call the whisper.cpp C/C++ wrapper here.
        // The provider boundary keeps the app local-first and testable.
        return try await WhisperRuntime.shared.transcribe(audioURL: audioFileURL, language: language)
    }
}

actor WhisperRuntime {
    static let shared = WhisperRuntime()

    func transcribe(audioURL: URL, language: SpeechLanguage) async throws -> String {
        _ = audioURL
        _ = language
        throw SpeechRecognitionError.engineUnavailable
    }
}
