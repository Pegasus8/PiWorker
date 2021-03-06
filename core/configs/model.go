package configs

import (
	"sync"
)

// Configs is the struct used to store all PiWorker configurations.
type Configs struct {
	Behavior   Behavior   `json:"behavior"`
	Security   Security   `json:"security"`
	Backups    Backups    `json:"backups"`
	APIConfigs APIConfigs `json:"api-configs"`
	Updates    Updates    `json:"updates"`
	WebUI      WebUI      `json:"webui"`
	Users      []User     `json:"users"`

	path         string
	sync.RWMutex `json:"-"`
}

// Behavior is the struct used to store Behavior configs of PiWorker.
type Behavior struct {
	LoopSleep int64 `json:"loop-sleep(ms)"`
}

// Security is the struct used to store configs related with the security of PiWorker.
type Security struct {
	DeniedIPs          []string `json:"denied-ips"`
	LocalNetworkAccess bool     `json:"local-network-access"`
}

// Backups is the struct used to store configs related with the backups of PiWorker (not implemented yet).
type Backups struct {
	BackupData        bool   `json:"backup-data"`
	BackupConfigs     bool   `json:"backup-configs"`
	DataBackupPath    string `json:"data-backup-path"`
	ConfigsBackupPath string `json:"configs-backup-path"`
	Freq              int16  `json:"frequency(hs)"`
}

// APIConfigs is the struct used to store configs related with the different APIs of PiWorker.
type APIConfigs struct {
	// APIs States
	NewTaskAPI     bool `json:"new-task-api"`
	EditTaskAPI    bool `json:"edit-task-api"`
	DeleteTaskAPI  bool `json:"delete-task-api"`
	GetAllTasksAPI bool `json:"get-all-tasks-api"`
	StatisticsAPI  bool `json:"statistics-api"`
	LogsAPI        bool `json:"logs-api"`
	TypesCompatAPI bool `json:"types-compat-api"`

	// Authentication
	RequireToken  bool   `json:"require-token"`
	SigningKey    string `json:"signing-key"`
	TokenDuration int64  `json:"token-duration(hs)"`
}

// Updates is the struct used to store configs related with the self-update of PiWorker (not implemented yet).
type Updates struct {
	DailyCheck     bool `json:"daily-check"`
	AutoDownload   bool `json:"auto-download"` // Only if daily check is active
	BugsPrevention bool `json:"bugs-prevention"`
}

// WebUI is the struct used to store configs related with the WebUI of PiWorker.
type WebUI struct {
	Enabled       bool   `json:"enabled"`
	ListeningPort string `json:"listening-port"`
}

// User is used to store each user's credentials.
type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password-hash"`
	Admin        bool   `json:"admin"`
}

// Sync writes the configs into the proper file, overwriting the previous content of it.
func (c *Configs) Sync() error {
	return writeToFile(c.path, c, true)
}

// unsafeSync does the same that `Configs.Sync()` with the difference that if does not lock the `RWMutex`, letting the
// responsibility of doing it to the developer. If the `RWMutex` is not locked, a **race condition** can be caused.
func (c *Configs) unsafeSync() error {
	return writeToFile(c.path, c, false)
}
