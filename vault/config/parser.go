package config

import (
        "fmt"
        "encoding/json"
        "os"
        "io/ioutil"
        "strings"
)

type VaultAuthMethod string

const (
        VaultAuthMethodAWS = VaultAuthMethod("aws")

        VaultAuthMethodToken = VaultAuthMethod("token")

        DefaultConfigFilePath string = "/etc/docker-credential-vault-login/config.json"

        EnvConfigFilePath string = "DOCKER_CREDS_CONFIG_FILE"
)

type CredHelperConfig struct {
        Method   VaultAuthMethod `json:"vault_auth_method"`
        Role     string          `json:"vault_role"`
        Secret   string          `json:"vault_secret_path"`
        ServerID string          `json:"vault_iam_server_id_header_value"`
        Path     string          `json:"-"`
}

func (c *CredHelperConfig) validate() error {
        var errors []string

        method := c.Method

        switch method {
        case "":
                errors = append(errors, `No Vault authentication method ("vault_auth_method") is provided`)
        case VaultAuthMethodAWS:
                if c.Role == "" {
                        errors = append(errors, fmt.Sprintf("%s %s", `No Vault role ("vault_role") is`,
                                "provided (required when the AWS authentication method is chosen)"))
                }
        case VaultAuthMethodToken:
                if v := os.Getenv("VAULT_TOKEN"); v == "" {
                        errors = append(errors, fmt.Sprintf("VAULT_TOKEN environment variable is not set"))
                }
        default:
                errors = append(errors, fmt.Sprintf("%s %s %q (must be either %q or %q)",
                        "Unrecognized Vault authentication method",
                        `("vault_auth_method") value`, method,
                        VaultAuthMethodAWS, VaultAuthMethodToken))
        }

        if c.Secret == "" {
                errors = append(errors, fmt.Sprintf("%s %s", "No path to the location of",
                        `your secret in Vault ("vault_secret_path") is provided`))
        }

        if len(errors) > 0 {
                return fmt.Errorf("Configuration file %s has the following errors:\n* %s", 
                        c.Path, strings.Join(errors, "\n* "))
        }
        return nil
}

func GetCredHelperConfig() (*CredHelperConfig, error) {
        cfg, err := parseConfig()
        if err != nil {
                return nil, err
        }

        if err = cfg.validate(); err != nil {
                return nil, err
        }
        return cfg, nil
}

func parseConfig() (*CredHelperConfig, error) {
        var path = DefaultConfigFilePath

        if v := os.Getenv(EnvConfigFilePath); v != "" {
                path = v
        }

        data, err := ioutil.ReadFile(path)
        if err != nil {
                return nil, err
        }

        var cfg = new(CredHelperConfig)
        if err = json.Unmarshal(data, cfg); err != nil {
                return cfg, err
        }
        cfg.Path = path
        return cfg, nil
}
