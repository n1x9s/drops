import Testing
@testable import SecondBrain

struct LocalMemoryEnricherTests {
    @Test
    func detectsLearningCategory() {
        let result = LocalMemoryEnricher().enrich("Study Go runtime netpoll before Sunday")
        #expect(result.category == .learning)
        #expect(result.tags.contains("study"))
    }
}
