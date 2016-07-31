package swagger

// SwaggerVersion show the current swagger version
const SwaggerVersion = "2.0"

type Swagger struct {
	SwaggerVer string `json:"swagger"`
	Info    Info   `json:"info"`
	Host    string `json:"host,omitempty"`   //"petstore.swagger.io"
	BasePath string `json:"basePath,omitempty"`  //"/v2"
	Schemes []string `json:"schemes"` //"http"
	Tags  []Tag `json:"tags"`
//	Paths []paths `json:"paths"`
}

type Info struct {
	// Required. The title of the application.
	Title string `json:"title,omitempty"`
	// A short description of the application. GFM syntax can be used for rich text representation.
	Description string `json:"description,omitempty"`
	// The Terms of Service for the API.
	TermsOfService string `json:"termsOfService,omitempty"`
	// The contact information for the exposed API.
	Contact Contact `json:"contact"`
	// The license information for the exposed API.
	License License `json:"license"`
	// Required Provides the version of the application API (not to be confused with the specification version).
	Version string `json:"version,omitempty"`
}

type Contact struct {
	// The identifying name of the contact person/organization.
	Name string `json:"name,omitempty"`
	// The URL pointing to the contact information. MUST be in the format of a URL.
	URL string `json:"url,omitempty"`
	// The email address of the contact person/organization. MUST be in the format of an email address.
	Email string `json:"email,omitempty"`
}

type License struct {
	// Required. The license name used for the API.
	Name string `json:"name,omitempty"`
	// A URL to the license used for the API. MUST be in the format of a URL.
	URL string `json:"url,omitempty"`
}

// Tag allows adding meta data to a single tag
type Tag struct {
	Name         string                `json:"name,omitempty"`
	Description  string                `json:"description,omitempty"`
	ExternalDocs ExternalDocumentation `json:"externalDocs"`
}

// ExternalDocumentation allows referencing an external resource for extended documentation.
type ExternalDocumentation struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

///////////////////////////
// ResourceListing list the resource
type ResourceListing struct {
	APIVersion     string `json:"apiVersion"`
	SwaggerVersion string `json:"swaggerVersion"` // e.g 1.2
						      // BasePath       string `json:"basePath"`  obsolete in 1.1
	APIs []APIRef    `json:"apis"`
	Info Information `json:"info"`
}

// APIRef description the api path and description
type APIRef struct {
	Path        string `json:"path"` // relative or absolute, must start with /
	Description string `json:"description"`
}

// Information show the API Information
type Information struct {
	Title             string `json:"title,omitempty"`
	Description       string `json:"description,omitempty"`
	Contact           string `json:"contact,omitempty"`
	TermsOfServiceURL string `json:"termsOfServiceUrl,omitempty"`
	License           string `json:"license,omitempty"`
	LicenseURL        string `json:"licenseUrl,omitempty"`
}

// APIDeclaration see https://github.com/wordnik/swagger-core/blob/scala_2.10-1.3-RC3/schemas/api-declaration-schema.json
type APIDeclaration struct {
	APIVersion     string           `json:"apiVersion"`
	SwaggerVersion string           `json:"swaggerVersion"`
	BasePath       string           `json:"basePath"`
	ResourcePath   string           `json:"resourcePath"` // must start with /
	Consumes       []string         `json:"consumes,omitempty"`
	Produces       []string         `json:"produces,omitempty"`
	APIs           []API            `json:"apis,omitempty"`
	Models         map[string]Model `json:"models,omitempty"`
}

// API show tha API struct
type API struct {
	Path        string      `json:"path"` // relative or absolute, must start with /
	Description string      `json:"description"`
	Operations  []Operation `json:"operations,omitempty"`
}

// Operation desc the Operation
type Operation struct {
	HTTPMethod string `json:"httpMethod"`
	Nickname   string `json:"nickname"`
	Type       string `json:"type"` // in 1.1 = DataType
					// ResponseClass    string            `json:"responseClass"` obsolete in 1.2
	Summary          string            `json:"summary,omitempty"`
	Notes            string            `json:"notes,omitempty"`
	Parameters       []Parameter       `json:"parameters,omitempty"`
	ResponseMessages []ResponseMessage `json:"responseMessages,omitempty"` // optional
	Consumes         []string          `json:"consumes,omitempty"`
	Produces         []string          `json:"produces,omitempty"`
	Authorizations   []Authorization   `json:"authorizations,omitempty"`
	Protocols        []Protocol        `json:"protocols,omitempty"`
}

// Protocol support which Protocol
type Protocol struct {
}

// ResponseMessage Show the
type ResponseMessage struct {
	Code          int    `json:"code"`
	Message       string `json:"message"`
	ResponseModel string `json:"responseModel"`
}

// Parameter desc the request parameters
type Parameter struct {
	ParamType     string `json:"paramType"` // path,query,body,header,form
	Name          string `json:"name"`
	Description   string `json:"description"`
	DataType      string `json:"dataType"` // 1.2 needed?
	Type          string `json:"type"`     // integer
	Format        string `json:"format"`   // int64
	AllowMultiple bool   `json:"allowMultiple"`
	Required      bool   `json:"required"`
	Default       string `json:"defaultValue"`
	Minimum       int    `json:"minimum"`
	Maximum       int    `json:"maximum"`
}

// ErrorResponse desc response
type ErrorResponse struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

// Model define the data model
type Model struct {
	ID         string                   `json:"id"`
	Required   []string                 `json:"required,omitempty"`
	Properties map[string]ModelProperty `json:"properties"`
}

// ModelProperty define the properties
type ModelProperty struct {
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Items       map[string]string `json:"items,omitempty"`
	Format      string            `json:"format"`
}

// Authorization see https://github.com/wordnik/swagger-core/wiki/authorizations
type Authorization struct {
	LocalOAuth OAuth  `json:"local-oauth"`
	APIKey     APIKey `json:"apiKey"`
}

// OAuth see https://github.com/wordnik/swagger-core/wiki/authorizations
type OAuth struct {
	Type       string               `json:"type"`   // e.g. oauth2
	Scopes     []string             `json:"scopes"` // e.g. PUBLIC
	GrantTypes map[string]GrantType `json:"grantTypes"`
}

// GrantType see https://github.com/wordnik/swagger-core/wiki/authorizations
type GrantType struct {
	LoginEndpoint        Endpoint `json:"loginEndpoint"`
	TokenName            string   `json:"tokenName"` // e.g. access_code
	TokenRequestEndpoint Endpoint `json:"tokenRequestEndpoint"`
	TokenEndpoint        Endpoint `json:"tokenEndpoint"`
}

// Endpoint see https://github.com/wordnik/swagger-core/wiki/authorizations
type Endpoint struct {
	URL              string `json:"url"`
	ClientIDName     string `json:"clientIdName"`
	ClientSecretName string `json:"clientSecretName"`
	TokenName        string `json:"tokenName"`
}

// APIKey see https://github.com/wordnik/swagger-core/wiki/authorizations
type APIKey struct {
	Type   string `json:"type"`   // e.g. apiKey
	PassAs string `json:"passAs"` // e.g. header
}