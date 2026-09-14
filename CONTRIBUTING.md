# Contributing to Karakorum Ultimate Compression 🏔️

First off, thanks for taking the time to contribute! 🎉

## Ways to Contribute

- 🐛 **Bug reports** — Open an issue with steps to reproduce
- 💡 **New algorithms** — Add a new `internal/algo/` implementation
- ⚡ **Performance** — Benchmarks and optimizations welcome
- 📖 **Docs** — Improve README or add examples

## Adding a New Algorithm

1. Create `internal/algo/youralgo.go`
2. Implement the `Algorithm` interface:
   ```go
   type Algorithm interface {
       Name() string
       Compress(src []byte) ([]byte, error)
       Decompress(src []byte, origLen uint32) ([]byte, error)
   }
   ```
3. Add your `ID` constant in `algo.go`
4. Register it in `internal/compressor/chunk.go`
5. Add benchmark results to the PR

## Development Setup

```bash
git clone https://github.com/summati/karakorum-ultimate-compression
cd karakorum-ultimate-compression
go build -o kuc ./cmd/kuc/
./kuc bench
```

## Running Tests

```bash
go test ./...
```

## Pull Request Guidelines

- Keep PRs focused — one feature per PR
- Add benchmark numbers in the PR description
- Ensure `./kuc bench` still passes all ✅ lossless checks
- Follow existing code style

## Code Style

- Standard Go formatting: `gofmt -w .`
- All public functions must have doc comments
- Errors must be wrapped with context: `fmt.Errorf("context: %w", err)`
