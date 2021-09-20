package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"log"
	"net/http"

	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/lachee/athenapdf/weaver/converter"
	"gopkg.in/alexcesaro/statsd.v2"
)

var (
	// ErrAuthorization should be returned when the authorization key is invalid.
	ErrAuthorization = errors.New("invalid authorization key provided")
	// ErrParams should be returned when there is missing parameters
	ErrParams = errors.New("missing or invalid query parameters")
	// ErrSignature should be returned when the HMAC computed does not match the one given
	ErrSignature = errors.New("invalid signature")
	// ErrInternalServer should be returned when a private error is returned
	// from a handler.
	ErrInternalServer = errors.New("PDF conversion failed due to an internal server error")
)

// ConfigMiddleware sets the config in the context.
func ConfigMiddleware(conf Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("config", conf)
	}
}

// WorkQueueMiddleware sets the work queue (write only) in the context.
func WorkQueueMiddleware(q chan<- converter.Work) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("queue", q)
	}
}

// SentryMiddleware sets the Sentry client (Raven) in the context.
func SentryMiddleware(r *raven.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("sentry", r)
	}
}

// StatsdMiddleware sets the Statsd client in the context.
func StatsdMiddleware(s *statsd.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("statsd", s)
	}
}

// ErrorMiddleware runs after all handlers have been executed, and it handles
// any errors returned from the handlers. It will return an internal server
// error with a predefined message if the last error type is not public.
// Otherwise, it will display the last error message it received, and the
// associated HTTP status code.
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		lastError := c.Errors.Last()
		statusCode := c.Writer.Status()

		if lastError != nil {
			// Log all errors
			log.Println("captured errors:")
			log.Printf("%+v\n", c.Errors)

			// Public errors
			if lastError.IsType(gin.ErrorTypePublic) {
				c.JSON(statusCode, gin.H{
					"error": lastError.Error(),
				})
				return
			}

			// Private errors
			c.JSON(500, gin.H{
				"error": ErrInternalServer.Error(),
			})
		}
	}
}

// AuthorizationMiddleware is a simple authorization middleware which matches
// an authentication key, provided via a query parameter, against a defined
// authentication key in the environment config.
func AuthorizationMiddleware(k string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Query("auth") != k {
			c.AbortWithError(http.StatusUnauthorized, ErrAuthorization).SetType(gin.ErrorTypePublic)
		}

		c.Next()
	}
}

// SignedMiddleware is a simple HMAC signing middleware which ensures the signed url (passed in the HMAC query)
// is correct and matches the key we have. Use this instead of AuthorizationMiddleware to abstract the key away an extra layer
// and make it unqiue per request.
func SignedMiddleware(k []byte) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Fetch URL
		url := c.Query("url")
		if url == "" {
			c.AbortWithError(http.StatusUnauthorized, ErrParams).SetType(gin.ErrorTypePublic)
		} else {
			// Fetch HMAC
			receivedMAC := c.Query("hmac")
			if receivedMAC == "" {
				c.AbortWithError(http.StatusUnauthorized, ErrParams).SetType(gin.ErrorTypePublic)
			} else {

				// Verify the HMAC matches the URL using the key
				mac := hmac.New(sha256.New, k)
				expectedMAC := mac.Sum([]byte(url))
				matches := hmac.Equal([]byte(receivedMAC), expectedMAC)
				if !matches {

					// Abort, invalid hmac
					c.AbortWithError(http.StatusUnauthorized, ErrSignature).SetType(gin.ErrorTypePublic)
				}
			}
		}

		c.Next()
	}
}
