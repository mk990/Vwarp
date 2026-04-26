package masque

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Diniboy1123/usque/api"
	"github.com/Diniboy1123/usque/config"
)

// AutoRegisterOrLoad loads an existing MASQUE config from configPath, or registers a new
// WARP device and saves the resulting config if none exists. Returns the config on success.
func AutoRegisterOrLoad(ctx context.Context, configPath, deviceName string) (*config.Config, error) {
	if configPath == "" {
		configPath = GetDefaultConfigPath()
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Try loading an existing config first.
	if _, err := os.Stat(configPath); err == nil {
		if err := config.LoadConfig(configPath); err == nil {
			cfg := config.AppConfig
			if cfg.PrivateKey != "" && cfg.EndpointV4 != "" && cfg.ID != "" {
				return &cfg, nil
			}
		}
		// Config invalid — remove and re-register.
		os.Remove(configPath)
	}

	if deviceName == "" {
		deviceName = "vwarp"
	}

	accountData, err := api.Register("PC", "en_US", "", true)
	if err != nil {
		return nil, fmt.Errorf("failed to register device: %w", err)
	}

	privKey, pubKey, err := generateEcKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}

	updatedAccountData, apiErr, err := api.EnrollKey(accountData, pubKey, deviceName)
	if err != nil {
		if apiErr != nil {
			return nil, fmt.Errorf("failed to enroll key: %w (API errors: %s)", err, apiErr.ErrorsAsString("; "))
		}
		return nil, fmt.Errorf("failed to enroll key: %w", err)
	}

	if len(updatedAccountData.Config.Peers) == 0 ||
		updatedAccountData.Config.Peers[0].Endpoint.V4 == "" ||
		updatedAccountData.Config.Peers[0].PublicKey == "" ||
		updatedAccountData.ID == "" {
		return nil, fmt.Errorf("registration returned incomplete data")
	}

	cfg := &config.Config{
		PrivateKey:     base64.StdEncoding.EncodeToString(privKey),
		EndpointV4:     stripPortSuffix(updatedAccountData.Config.Peers[0].Endpoint.V4),
		EndpointV6:     stripPortSuffix(updatedAccountData.Config.Peers[0].Endpoint.V6),
		EndpointPubKey: updatedAccountData.Config.Peers[0].PublicKey,
		License:        updatedAccountData.Account.License,
		ID:             updatedAccountData.ID,
		AccessToken:    accountData.Token,
		IPv4:           updatedAccountData.Config.Interface.Addresses.V4,
		IPv6:           updatedAccountData.Config.Interface.Addresses.V6,
	}

	if err := saveConfigFile(configPath, cfg); err != nil {
		return nil, fmt.Errorf("failed to save config: %w", err)
	}

	return cfg, nil
}

// LoadMasqueConfig loads and validates a MASQUE config from disk.
func LoadMasqueConfig(configPath string) (*config.Config, error) {
	if err := config.LoadConfig(configPath); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	cfg := config.AppConfig
	if cfg.PrivateKey == "" || cfg.ID == "" {
		return nil, fmt.Errorf("config at %s is incomplete", configPath)
	}
	return &cfg, nil
}
