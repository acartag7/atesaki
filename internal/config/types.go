package config

// Fixed v0 numbers from docs/contract-boundaries.md B8 that this package
// enforces. Values with a B1 key are defaults; the rest have no config key.
const (
	APIVersion = "atesaki/v1alpha1"

	defaultAccessTTL      int64 = 600     // 10 min
	defaultConsentTTL     int64 = 300     // 5 min
	defaultCodeTTL        int64 = 120     // 2 min
	defaultMaxDuration    int64 = 28800   // 8 h
	hardMaxDuration       int64 = 2592000 // 30 days
	defaultPairingCodeTTL int64 = 600     // 10 min

	purposeMaxBytes  = 512
	configMaxBytes   = 1 << 20  // 1 MiB
	secretFileMax    = 64 << 10 // 64 KiB
	scopeMaxEntries  = 128
	groupsMaxEntries = 128

	defaultCallbackPath = "/oauth/callback"
	defaultMCPEndpoint  = "/mcp"
	defaultLivePath     = "/healthz"
	defaultReadyPath    = "/readyz"
	defaultGroupsClaim  = "groups"

	prmPrefix    = "/.well-known/oauth-protected-resource"
	asMetadata   = "/.well-known/oauth-authorization-server"
	oidcMetadata = "/.well-known/openid-configuration"
)

// Config is a fully validated configuration: exactly one Gateway and at least
// one Route, every cross-reference resolved.
type Config struct {
	Gateway *Gateway
	Routes  []*Route
}

type Gateway struct {
	Name              string
	ExternalBaseURL   string // canonical origin, the issuer
	Listen            string
	ListenHost        string
	ListenPort        int
	SigningKeyRef     Ref
	AllowedHosts      []string
	AllowedOrigins    []string
	TrustProxyHeaders bool
	TrustedProxies    []string
	Identity          Identity
	Clients           Clients
	Egress            map[string]EgressProfile
	StorePath         string
	AuditPath         string
	Health            Health
	Tokens            Tokens
	MachineClients    []MachineClient
}

type Identity struct {
	Provider        string // entra | oidc | header | console
	Registration    string // dedicated | shared
	RedirectURI     string
	CallbackPath    string
	ClientID        string
	ClientSecretRef *Ref
	PublicClient    bool
	TenantID        string
	Issuer          string
	GroupsClaim     string
	EgressProfile   string
	Assertion       *Assertion
	GroupsToScopes  map[string][]string
	PairingCodeTTL  int64
}

// Assertion is the B4 signed-assertion contract for identity.provider: header.
type Assertion struct {
	Header       string
	Issuer       string
	Audience     string
	Alg          string // ES256 | RS256, pinned
	Keys         AssertionKeys
	SubjectClaim string
	GroupsClaim  string
}

// AssertionKeys is the B4 union: {jwksUrl, egressProfile} or {jwksRef}.
type AssertionKeys struct {
	JWKSURL       string
	EgressProfile string
	JWKSRef       *Ref
}

type Clients struct {
	KnownCIMD              []string
	AllowLoopbackRedirects bool
	RedirectAllowlist      []string
}

type EgressProfile struct {
	Proxy       string // fromEnv | none | URL
	CABundleRef *Ref
}

type Health struct {
	LivePath  string
	ReadyPath string
}

type Tokens struct {
	AccessTTL  int64
	ConsentTTL int64
	CodeTTL    int64
}

type MachineClient struct {
	ID          string
	SecretRef   Ref
	Purpose     string
	MaxDuration int64
	Routes      []MachineRoute
}

type MachineRoute struct {
	Path   string
	Scopes []string
}

type Route struct {
	Name          string
	Path          string
	MCPEndpoint   string
	Upstream      Upstream
	Catalog       []string
	DefaultScopes []string
	RequireScope  string
	MaxDuration   int64
	Rules         []PolicyRule
	Approvers     []string

	// Derived (B3): audience and metadata locations, byte-exact.
	Audience    string
	MCPPath     string
	MetadataURL string
}

type Upstream struct {
	URL           string
	EgressProfile string
	Credential    Credential
}

type Credential struct {
	Type     string // none | static-header
	Header   string
	Scheme   string
	ValueRef *Ref
}

// PolicyRule is one G7 rule: AND-only conditions, effect allow or deny.
type PolicyRule struct {
	Effect         string
	ScopesSubsetOf []string
	DurationAtMost *int64
	SubjectInGroup string
	ClientIn       []string
}
