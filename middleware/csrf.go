package middleware

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/random"
)

type (
	// CSRFConfig defines the config for CSRF middleware.
	CSRFConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// TokenLength is the length of the generated token.
		TokenLength uint8 `yaml:"token_length"`
		// Optional. Default value 32.

		// TokenLookup is a string in the form of "<source>:<name>" or "<source>:<name>,<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:X-CSRF-Token".
		// Possible values:
		// - "header:<name>" or "header:<name>:<cut-prefix>"
		// - "query:<name>"
		// - "form:<name>"
		// Multiple sources example:
		// - "header:X-CSRF-Token,query:csrf"
		TokenLookup string `yaml:"token_lookup"`

		HeaderName string `yaml:"header_name"`

		// Context key to store generated CSRF token into context.
		// Optional. Default value "csrf".
		ContextKey string `yaml:"context_key"`

		// Name of the CSRF cookie. This cookie will store CSRF token.
		// Optional. Default value "csrf".
		CookieName string `yaml:"cookie_name"`

		// Domain of the CSRF cookie.
		// Optional. Default value none.
		CookieDomain string `yaml:"cookie_domain"`

		// Path of the CSRF cookie.
		// Optional. Default value none.
		CookiePath string `yaml:"cookie_path"`

		// Max age (in seconds) of the CSRF cookie.
		// Optional. Default value 86400 (24hr).
		CookieMaxAge int `yaml:"cookie_max_age"`

		// Indicates if CSRF cookie is secure.
		// Optional. Default value false.
		CookieSecure bool `yaml:"cookie_secure"`

		// Indicates if CSRF cookie is HTTP only.
		// Optional. Default value false.
		CookieHTTPOnly bool `yaml:"cookie_http_only"`

		// Indicates SameSite mode of the CSRF cookie.
		// Optional. Default value SameSiteDefaultMode.
		CookieSameSite http.SameSite `yaml:"cookie_same_site"`

		// ErrorHandler defines a function which is executed for returning custom errors.
		ErrorHandler CSRFErrorHandler
	}

	// CSRFErrorHandler is a function which is executed for creating custom errors.
	CSRFErrorHandler func(err error, c echo.Context) error
)

var errHeaderExtractorValueMissing = errors.New("missing value in request header")
var errHeaderExtractorValueInvalid = errors.New("invalid value in request header")
var errQueryExtractorValueMissing = errors.New("missing value in the query string")
var errParamExtractorValueMissing = errors.New("missing value in path params")
var errCookieExtractorValueMissing = errors.New("missing value in cookies")
var errFormExtractorValueMissing = errors.New("missing value in the form")

// ErrCSRFInvalid is returned when CSRF check fails
var ErrCSRFInvalid = echo.NewHTTPError(http.StatusForbidden, "invalid csrf token")

// ErrCSRFMissing is returned when CSRF token is missing
var clientTokens = make(map[string]string, 100)

var (
	// DefaultCSRFConfig is the default CSRF middleware config.
	DefaultCSRFConfig = CSRFConfig{
		Skipper:        middleware.DefaultSkipper,
		TokenLength:    32,
		TokenLookup:    "header:" + echo.HeaderXCSRFToken,
		HeaderName:     echo.HeaderXCSRFToken,
		ContextKey:     "csrf",
		CookieName:     "_csrf",
		CookieMaxAge:   86400,
		CookieSameSite: http.SameSiteDefaultMode,
	}
)

// CSRF returns a Cross-Site Request Forgery (CSRF) middleware.
// See: https://en.wikipedia.org/wiki/Cross-site_request_forgery
func CSRF() echo.MiddlewareFunc {
	c := DefaultCSRFConfig
	return CSRFWithConfig(c)
}

// CSRFWithConfig returns a CSRF middleware with config.
// See `CSRF()`.
func CSRFWithConfig(config CSRFConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultCSRFConfig.Skipper
	}
	if config.TokenLength == 0 {
		config.TokenLength = DefaultCSRFConfig.TokenLength
	}
	if config.TokenLookup == "" {
		config.TokenLookup = DefaultCSRFConfig.TokenLookup
	}
	if config.ContextKey == "" {
		config.ContextKey = DefaultCSRFConfig.ContextKey
	}
	if config.CookieName == "" {
		config.CookieName = DefaultCSRFConfig.CookieName
	}
	if config.CookieMaxAge == 0 {
		config.CookieMaxAge = DefaultCSRFConfig.CookieMaxAge
	}
	if config.CookieSameSite == http.SameSiteNoneMode {
		config.CookieSecure = true
	}

	// extractors, cErr := middleware.CreateExtractors(config.TokenLookup)
	// if cErr != nil {
	// 	panic(cErr)
	// }

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			token := ""
			if k := c.Request().Header.Get(echo.HeaderXCSRFToken); k == "" {
				token = random.String(config.TokenLength) // Generate token
			} else {
				token = k // Reuse token
			}
			fmt.Println("token: ", token)

			switch c.Request().Method {
			case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
			default:
				// Validate token only for requests which are not defined as 'safe' by RFC7231
				var lastExtractorErr error
				var lastTokenErr error
				// fmt.Println("lastExtractorErr: ", lastExtractorErr, lastTokenErr)
			outer:
				// for _, extractor := range extractors {
				// clientTokens, err := extractor(c)
				// if err != nil {
				// 	lastExtractorErr = err
				// 	continue
				// }

				for _, clientToken := range clientTokens {
					fmt.Println("clientToken: ", clientToken)
					if validateCSRFToken(token, clientToken) {
						lastTokenErr = nil
						lastExtractorErr = nil
						break outer
					}
					lastTokenErr = ErrCSRFInvalid
				}
				// }
				if len(clientTokens) == 0 {
					lastExtractorErr = echo.NewHTTPError(http.StatusBadRequest, "invalid csrf token in request header")
				}
				var finalErr error
				if lastTokenErr != nil {
					finalErr = lastTokenErr
				} else if lastExtractorErr != nil {
					// ugly part to preserve backwards compatible errors. someone could rely on them
					if lastExtractorErr == errQueryExtractorValueMissing {
						lastExtractorErr = echo.NewHTTPError(http.StatusBadRequest, "missing csrf token in the query string")
					} else if lastExtractorErr == errFormExtractorValueMissing {
						lastExtractorErr = echo.NewHTTPError(http.StatusBadRequest, "missing csrf token in the form parameter")
					} else if lastExtractorErr == errHeaderExtractorValueMissing {
						lastExtractorErr = echo.NewHTTPError(http.StatusBadRequest, "missing csrf token in request header")
					} else {
						lastExtractorErr = echo.NewHTTPError(http.StatusBadRequest, lastExtractorErr.Error())
					}
					finalErr = lastExtractorErr
				}

				if finalErr != nil {
					if config.ErrorHandler != nil {
						return config.ErrorHandler(finalErr, c)
					}
					return finalErr
				}
			}

			// Store token in the context
			c.Set(config.ContextKey, token)
			clientTokens[token] = "token"

			// Protect clients from caching the response
			c.Response().Header().Add(echo.HeaderVary, echo.HeaderCookie)
			c.Response().Header().Add("X-CSRF-Token", token)

			return next(c)
		}
	}
}

func validateCSRFToken(token, clientToken string) bool {
	return subtle.ConstantTimeCompare([]byte(token), []byte(clientToken)) == 1
}