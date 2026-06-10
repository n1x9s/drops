import Foundation

struct GeminiProvider: LLMProvider {
    let apiKey: String
    let model: String

    init(apiKey: String, model: String = "gemini-2.5-flash") {
        self.apiKey = apiKey
        self.model = model
    }

    func enrichMemory(_ text: String) async throws -> MemoryEnrichment {
        guard !apiKey.isEmpty else {
            return LocalMemoryEnricher().enrich(text)
        }
        let url = URL(string: "https://generativelanguage.googleapis.com/v1beta/models/\(model):generateContent?key=\(apiKey)")!
        let prompt = """
        Return JSON only: {"summary":"one sentence","category":"Work|Learning|Personal|Projects|Meetings|Ideas","tags":["lowercase"]}.
        Text: \(text)
        """
        let body = GeminiRequest(contents: [.init(parts: [.init(text: prompt)])])
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.httpBody = try JSONEncoder().encode(body)
        let (data, _) = try await URLSession.shared.data(for: request)
        let response = try JSONDecoder().decode(GeminiResponse.self, from: data)
        guard let raw = response.candidates.first?.content.parts.first?.text,
              let json = raw.data(using: .utf8) else {
            return LocalMemoryEnricher().enrich(text)
        }
        return (try? JSONDecoder().decode(MemoryEnrichment.self, from: json)) ?? LocalMemoryEnricher().enrich(text)
    }
}

private struct GeminiRequest: Encodable {
    let contents: [Content]

    struct Content: Encodable {
        let parts: [Part]
    }

    struct Part: Encodable {
        let text: String
    }
}

private struct GeminiResponse: Decodable {
    let candidates: [Candidate]

    struct Candidate: Decodable {
        let content: Content
    }

    struct Content: Decodable {
        let parts: [Part]
    }

    struct Part: Decodable {
        let text: String
    }
}
