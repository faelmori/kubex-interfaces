package settings

const StructureFlag = "gospider_config_structure"
const CertificatesFlag = "gospider_certificates"
const RedisPasswordFlag = "gospider_redis_password"
const RefreshSecretFlag = "gospider_refresh_secret"
const DBPasswordFlag = "gospider_db_password"
const CacheSetupFlag = "gospider_cache_setup"
const ServicesSetupFlag = "gospider_services_setup"
const VaultSetupFlag = "gospider_vault_setup"
const DepsSetupFlag = "gospider_deps_setup"

const KeyringService = "kubex"
const KeyringKey = "gospider_tls_pass"

const DefaultKubexDir = "$HOME/.kubex"
const DefaultCacheDir = "$HOME/.cache/kubex/gospider"
const DefaultVaultDir = "$HOME/.kubex/.vault"

const DefaultGoSpiderDir = "$HOME/.kubex/gospider"
const DefaultGoSpiderVaultDir = "$HOME/.kubex/gospider/.vault"
const DefaultGoSpiderConfigDir = "$HOME/.kubex/gospider/config"
const DefaultGoSpiderConfigPath = "$HOME/.kubex/gospider/config/config.json"
const DefaultGoSpiderKeyPath = "$HOME/.kubex/gospider/gospider-key.pem"
const DefaultGoSpiderCertPath = "$HOME/.kubex/gospider/gospider-cert.pem"
