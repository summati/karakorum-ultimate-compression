package license

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	EvalStart string `json:"eval_start"`
	Key       string `json:"license_key"`
}

func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".kuc_config.json"
	}
	return filepath.Join(home, ".kuc_config.json")
}

func loadConfig() (*Config, error) {
	path := getConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func saveConfig(cfg *Config) error {
	path := getConfigPath()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func isValidKey(key string) bool {
	// Simple offline check: require it starts with "KUC-" and is at least 15 chars long.
	key = strings.TrimSpace(key)
	return strings.HasPrefix(key, "KUC-") && len(key) >= 15
}

// Check Enforcement verifies if the user is allowed to run the program.
func CheckEnforcement() error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("could not read license config: %w", err)
	}

	// 1. If they have a valid key, they are allowed to use it forever.
	if cfg.Key != "" && isValidKey(cfg.Key) {
		return nil
	}

	// 2. If no key, check if this is the very first run.
	if cfg.EvalStart == "" {
		cfg.EvalStart = time.Now().Format(time.RFC3339)
		saveConfig(cfg)
		fmt.Println("🚀 First run! Your 14-day evaluation period has started.")
		return nil
	}

	// 3. Check if 14 days have passed.
	evalStart, err := time.Parse(time.RFC3339, cfg.EvalStart)
	if err != nil {
		// Reset if corrupted
		cfg.EvalStart = time.Now().Format(time.RFC3339)
		saveConfig(cfg)
		return nil
	}

	daysPassed := time.Since(evalStart).Hours() / 24
	if daysPassed > 14 {
		return errors.New("Evaluation period expired (14 days). Please purchase a license and run 'kuc auth <key>'.")
	}

	daysLeft := 14 - int(daysPassed)
	fmt.Printf("⏳ Evaluation period: %d days remaining. Purchase a license at https://summati.gumroad.com/l/karakorum\n", daysLeft)
	return nil
}

func Authenticate(key string) error {
	if !isValidKey(key) {
		return errors.New("invalid license key format (must start with KUC- and be 15+ characters)")
	}

	cfg, err := loadConfig()
	if err != nil {
		cfg = &Config{}
	}

	cfg.Key = strings.TrimSpace(key)
	if err := saveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save license key: %w", err)
	}
	return nil
}
