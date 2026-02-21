package config

const (
	DefaultBrowserType    = "chromium"
	DefaultBrowserMode    = "dedicated"
	DefaultChromeChannel  = "chromium"
	DefaultHost           = "localhost"
	DefaultDebugPort      = 9222
	DefaultViewportWidth  = 1080
	DefaultViewportHeight = 600
	DefaultUserAgent      = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/116.0.0.0 Safari/537.36"
)

type ProxyConfig struct {
	Server   string `json:"server,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	IP       string `json:"ip,omitempty"`
}

type BrowserConfig struct {
	BrowserType           string            `json:"browser_type"`
	Headless              bool              `json:"headless"`
	BrowserMode           string            `json:"browser_mode"`
	UseManagedBrowser     bool              `json:"use_managed_browser"`
	CDPURL                string            `json:"cdp_url,omitempty"`
	BrowserContextID      string            `json:"browser_context_id,omitempty"`
	TargetID              string            `json:"target_id,omitempty"`
	CDPCleanupOnClose     bool              `json:"cdp_cleanup_on_close"`
	CreateIsolatedContext bool              `json:"create_isolated_context"`
	UsePersistentContext  bool              `json:"use_persistent_context"`
	UserDataDir           string            `json:"user_data_dir,omitempty"`
	ChromeChannel         string            `json:"chrome_channel"`
	Channel               string            `json:"channel"`
	Proxy                 string            `json:"proxy,omitempty"`
	ProxyConfig           *ProxyConfig      `json:"proxy_config,omitempty"`
	ViewportWidth         int               `json:"viewport_width"`
	ViewportHeight        int               `json:"viewport_height"`
	AcceptDownloads       bool              `json:"accept_downloads"`
	DownloadsPath         string            `json:"downloads_path,omitempty"`
	StorageState          any               `json:"storage_state,omitempty"`
	IgnoreHTTPSErrors     bool              `json:"ignore_https_errors"`
	JavaScriptEnabled     bool              `json:"java_script_enabled"`
	SleepOnClose          bool              `json:"sleep_on_close"`
	Verbose               bool              `json:"verbose"`
	Cookies               []map[string]any  `json:"cookies,omitempty"`
	Headers               map[string]string `json:"headers,omitempty"`
	UserAgent             string            `json:"user_agent"`
	UserAgentMode         string            `json:"user_agent_mode,omitempty"`
	UserAgentGeneratorCfg map[string]any    `json:"user_agent_generator_config,omitempty"`
	TextMode              bool              `json:"text_mode"`
	LightMode             bool              `json:"light_mode"`
	ExtraArgs             []string          `json:"extra_args,omitempty"`
	DebuggingPort         int               `json:"debugging_port"`
	Host                  string            `json:"host"`
	EnableStealth         bool              `json:"enable_stealth"`
	InitScripts           []string          `json:"init_scripts,omitempty"`
}

func DefaultBrowserConfig() BrowserConfig {
	return BrowserConfig{
		BrowserType:       DefaultBrowserType,
		Headless:          true,
		BrowserMode:       DefaultBrowserMode,
		UseManagedBrowser: false,
		ChromeChannel:     DefaultChromeChannel,
		Channel:           DefaultChromeChannel,
		ViewportWidth:     DefaultViewportWidth,
		ViewportHeight:    DefaultViewportHeight,
		AcceptDownloads:   false,
		IgnoreHTTPSErrors: true,
		JavaScriptEnabled: true,
		SleepOnClose:      false,
		Verbose:           true,
		Cookies:           []map[string]any{},
		Headers:           map[string]string{},
		UserAgent:         DefaultUserAgent,
		TextMode:          false,
		LightMode:         false,
		ExtraArgs:         []string{},
		DebuggingPort:     DefaultDebugPort,
		Host:              DefaultHost,
		EnableStealth:     false,
		InitScripts:       []string{},
	}
}
