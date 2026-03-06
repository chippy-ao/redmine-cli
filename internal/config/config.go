package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Profile はRedmineサーバーへの接続情報を保持する。
type Profile struct {
	URL    string `yaml:"url"`
	APIKey string `yaml:"api_key"`
}

// Config はCLI全体の設定を保持する。
type Config struct {
	DefaultProfile string             `yaml:"default_profile"`
	Profiles       map[string]Profile `yaml:"profiles"`
}

// DefaultConfigPath は設定ファイルのデフォルトパスを返す。
func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".config", "redmine-cli", "config.yaml")
}

// LoadConfig は指定パスからYAML設定ファイルを読み込む。
// ファイルが存在しない場合は空のConfigを返す（エラーにしない）。
func LoadConfig(path string) (*Config, error) {
	cfg := &Config{
		Profiles: make(map[string]Profile),
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("設定ファイルのパースに失敗しました: %w", err)
	}

	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]Profile)
	}

	// URLの末尾スラッシュを除去
	for name, p := range cfg.Profiles {
		p.URL = strings.TrimRight(p.URL, "/")
		cfg.Profiles[name] = p
	}

	return cfg, nil
}

// SaveConfig は設定をYAMLファイルとして保存する。
// ディレクトリは0700、ファイルは0600のパーミッションで作成する。
func SaveConfig(path string, cfg *Config) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("設定ディレクトリの作成に失敗しました: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("設定のシリアライズに失敗しました: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("設定ファイルの書き込みに失敗しました: %w", err)
	}

	return nil
}

// GetProfile は指定された名前のプロファイルを返す。
// nameが空の場合はDefaultProfileを使用する。
func (c *Config) GetProfile(name string) (*Profile, error) {
	if name == "" {
		if c.DefaultProfile == "" {
			return nil, fmt.Errorf("プロファイルが設定されていません。redmine-cli config add で設定してください。")
		}
		name = c.DefaultProfile
	}

	p, ok := c.Profiles[name]
	if !ok {
		return nil, fmt.Errorf("プロファイル %q が見つかりません。redmine-cli config list で確認してください。", name)
	}

	return &p, nil
}
