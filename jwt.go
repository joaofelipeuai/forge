package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// JWT represents a JSON Web Token
type JWT struct {
	Header    JWTHeader    `json:"header"`
	Payload   JWTPayload   `json:"payload"`
	Signature string       `json:"signature"`
	Raw       string       `json:"raw"`
}

// JWTHeader represents the JWT header
type JWTHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

// JWTPayload represents the JWT payload
type JWTPayload struct {
	Issuer         string                 `json:"iss,omitempty"`
	Subject        string                 `json:"sub,omitempty"`
	Audience       string                 `json:"aud,omitempty"`
	ExpirationTime int64                  `json:"exp,omitempty"`
	NotBefore      int64                  `json:"nbf,omitempty"`
	IssuedAt       int64                  `json:"iat,omitempty"`
	JWTID          string                 `json:"jti,omitempty"`
	Claims         map[string]interface{} `json:"-"`
}

// JWTConfig represents JWT configuration
type JWTConfig struct {
	Secret     string
	Issuer     string
	Expiration time.Duration
	Algorithm  string
}

// NewJWTConfig creates a new JWT configuration
func NewJWTConfig(secret string) *JWTConfig {
	return &JWTConfig{
		Secret:     secret,
		Issuer:     "velocity-framework",
		Expiration: 24 * time.Hour,
		Algorithm:  "HS256",
	}
}

// GenerateToken generates a new JWT token
func (config *JWTConfig) GenerateToken(claims map[string]interface{}) (string, error) {
	now := time.Now()
	
	header := JWTHeader{
		Algorithm: config.Algorithm,
		Type:      "JWT",
	}
	
	payload := JWTPayload{
		Issuer:         config.Issuer,
		IssuedAt:       now.Unix(),
		ExpirationTime: now.Add(config.Expiration).Unix(),
		Claims:         claims,
	}
	
	// Merge custom claims into payload
	payloadMap := make(map[string]interface{})
	payloadBytes, _ := json.Marshal(payload)
	json.Unmarshal(payloadBytes, &payloadMap)
	
	for key, value := range claims {
		payloadMap[key] = value
	}
	
	// Encode header
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerBytes)
	
	// Encode payload
	payloadBytes, err = json.Marshal(payloadMap)
	if err != nil {
		return "", err
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)
	
	// Create signature
	message := headerEncoded + "." + payloadEncoded
	signature := config.createSignature(message)
	
	return message + "." + signature, nil
}

// ValidateToken validates a JWT token
func (config *JWTConfig) ValidateToken(tokenString string) (*JWT, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}
	
	headerEncoded, payloadEncoded, signatureEncoded := parts[0], parts[1], parts[2]
	
	// Verify signature
	message := headerEncoded + "." + payloadEncoded
	expectedSignature := config.createSignature(message)
	if !hmac.Equal([]byte(signatureEncoded), []byte(expectedSignature)) {
		return nil, errors.New("invalid signature")
	}
	
	// Decode header
	headerBytes, err := base64.RawURLEncoding.DecodeString(headerEncoded)
	if err != nil {
		return nil, err
	}
	var header JWTHeader
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, err
	}
	
	// Decode payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return nil, err
	}
	var payloadMap map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payloadMap); err != nil {
		return nil, err
	}
	
	// Check expiration
	if exp, ok := payloadMap["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, errors.New("token expired")
		}
	}
	
	// Create JWT struct
	jwt := &JWT{
		Header:    header,
		Signature: signatureEncoded,
		Raw:       tokenString,
	}
	
	// Fill payload
	if iss, ok := payloadMap["iss"].(string); ok {
		jwt.Payload.Issuer = iss
	}
	if sub, ok := payloadMap["sub"].(string); ok {
		jwt.Payload.Subject = sub
	}
	if aud, ok := payloadMap["aud"].(string); ok {
		jwt.Payload.Audience = aud
	}
	if exp, ok := payloadMap["exp"].(float64); ok {
		jwt.Payload.ExpirationTime = int64(exp)
	}
	if nbf, ok := payloadMap["nbf"].(float64); ok {
		jwt.Payload.NotBefore = int64(nbf)
	}
	if iat, ok := payloadMap["iat"].(float64); ok {
		jwt.Payload.IssuedAt = int64(iat)
	}
	if jti, ok := payloadMap["jti"].(string); ok {
		jwt.Payload.JWTID = jti
	}
	
	// Store custom claims
	jwt.Payload.Claims = make(map[string]interface{})
	for key, value := range payloadMap {
		if key != "iss" && key != "sub" && key != "aud" && key != "exp" && 
		   key != "nbf" && key != "iat" && key != "jti" {
			jwt.Payload.Claims[key] = value
		}
	}
	
	return jwt, nil
}

// createSignature creates HMAC signature
func (config *JWTConfig) createSignature(message string) string {
	h := hmac.New(sha256.New, []byte(config.Secret))
	h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// JWT Authentication Middleware
func JWTAuth(config *JWTConfig) MiddlewareFunc {
	return func(c *Context) error {
		// Get token from Authorization header
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(401, map[string]string{"error": "Missing authorization header"})
		}
		
		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(401, map[string]string{"error": "Invalid authorization header format"})
		}
		
		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		// Validate token
		jwt, err := config.ValidateToken(tokenString)
		if err != nil {
			return c.JSON(401, map[string]string{"error": fmt.Sprintf("Invalid token: %v", err)})
		}
		
		// Store JWT in context
		c.Set("jwt", jwt)
		c.Set("user_id", jwt.Payload.Subject)
		
		// Store custom claims
		for key, value := range jwt.Payload.Claims {
			c.Set(key, value)
		}
		
		return c.Next()
	}
}

// Optional JWT middleware (doesn't fail if token is missing)
func JWTOptional(config *JWTConfig) MiddlewareFunc {
	return func(c *Context) error {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if jwt, err := config.ValidateToken(tokenString); err == nil {
				c.Set("jwt", jwt)
				c.Set("user_id", jwt.Payload.Subject)
				for key, value := range jwt.Payload.Claims {
					c.Set(key, value)
				}
			}
		}
		return c.Next()
	}
}

// Helper function to get JWT from context
func GetJWT(c *Context) *JWT {
	if jwt := c.Get("jwt"); jwt != nil {
		return jwt.(*JWT)
	}
	return nil
}

// Helper function to get user ID from context
func GetUserID(c *Context) string {
	if userID := c.Get("user_id"); userID != nil {
		return userID.(string)
	}
	return ""
}