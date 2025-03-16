package transformer

import (
    "encoding/json"
    "fmt"
    "os"
)

type KeyConfig struct {
    Fields     string `json:"fields"`
    Resolvable *bool  `json:"resolvable,omitempty"`
}

type TypeConfig struct {
    Keys     []KeyConfig `json:"keys"`
    External []string    `json:"external,omitempty"`
}

type DefaultConfig struct {
    Key *KeyConfig `json:"key,omitempty"`
}

type Config struct {
    Types    map[string]TypeConfig `json:"types"`
    Defaults *DefaultConfig        `json:"defaults,omitempty"`
}

func loadConfig(path string) (Config, error) {
    defaultResolvable := true
    config := Config{
        Types: make(map[string]TypeConfig),
        Defaults: &DefaultConfig{
            Key: &KeyConfig{
                Fields:     "id",
                Resolvable: &defaultResolvable,
            },
        },
    }

    if path == "" {
        return config, nil
    }

    data, err := os.ReadFile(path)
    if err != nil {
        return config, fmt.Errorf("設定ファイルの読み込みエラー: %w", err)
    }

    if err := json.Unmarshal(data, &config); err != nil {
        return config, fmt.Errorf("JSON解析エラー: %w", err)
    }

    return config, nil
}
