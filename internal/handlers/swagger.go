package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Bixor-Engine/backend/docs"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag"
)

// FilterSwaggerSpec filters the Swagger spec by tags and updates title/description/security
func FilterSwaggerSpec(tags []string, title, description string, securityDefs map[string]interface{}) (map[string]interface{}, error) {
	// Get the full Swagger spec as JSON string
	swaggerJSON, err := swag.ReadDoc(docs.SwaggerInfo.InstanceName())
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var fullSpec map[string]interface{}
	if err := json.Unmarshal([]byte(swaggerJSON), &fullSpec); err != nil {
		return nil, err
	}

	// Create a map of allowed tags for quick lookup
	allowedTags := make(map[string]bool)
	for _, tag := range tags {
		allowedTags[tag] = true
	}

	// Filter paths and collect referenced model names
	paths, ok := fullSpec["paths"].(map[string]interface{})
	if !ok {
		return fullSpec, nil
	}

	filteredPaths := make(map[string]interface{})
	referencedModels := make(map[string]bool)

	for path, pathItem := range paths {
		pathItemMap, ok := pathItem.(map[string]interface{})
		if !ok {
			continue
		}

		filteredPathItem := make(map[string]interface{})
		for method, methodItem := range pathItemMap {
			methodItemMap, ok := methodItem.(map[string]interface{})
			if !ok {
				filteredPathItem[method] = methodItem
				continue
			}

			// Check if this method has any of the allowed tags
			methodTags, ok := methodItemMap["tags"].([]interface{})
			if !ok {
				continue
			}

			hasAllowedTag := false
			for _, tagInterface := range methodTags {
				tag, ok := tagInterface.(string)
				if ok && allowedTags[tag] {
					hasAllowedTag = true
					break
				}
			}

			if hasAllowedTag {
				filteredPathItem[method] = methodItem

				// Extract model references from responses and parameters
				extractModelReferences(methodItemMap, referencedModels)
			}
		}

		if len(filteredPathItem) > 0 {
			filteredPaths[path] = filteredPathItem
		}
	}

	// Filter definitions to only include referenced models
	definitions, ok := fullSpec["definitions"].(map[string]interface{})
	filteredDefinitions := make(map[string]interface{})
	if ok {
		for modelName, modelDef := range definitions {
			// Check if this model is referenced
			if referencedModels[modelName] {
				filteredDefinitions[modelName] = modelDef
			}
		}
	}

	// Create filtered spec
	filteredSpec := make(map[string]interface{})
	for k, v := range fullSpec {
		if k == "paths" {
			filteredSpec[k] = filteredPaths
		} else if k == "definitions" {
			filteredSpec[k] = filteredDefinitions
		} else if k == "info" {
			// Update info section with custom title and description
			info, ok := v.(map[string]interface{})
			if ok {
				infoCopy := make(map[string]interface{})
				for infoKey, infoValue := range info {
					if infoKey == "title" {
						infoCopy[infoKey] = title
					} else if infoKey == "description" {
						infoCopy[infoKey] = description
					} else {
						infoCopy[infoKey] = infoValue
					}
				}
				filteredSpec[k] = infoCopy
			} else {
				filteredSpec[k] = v
			}
		} else if k == "securityDefinitions" {
			// Replace security definitions with custom ones
			if securityDefs != nil {
				filteredSpec[k] = securityDefs
			} else {
				filteredSpec[k] = v
			}
		} else {
			filteredSpec[k] = v
		}
	}

	return filteredSpec, nil
}

// extractModelReferences extracts model references from Swagger spec items
func extractModelReferences(item map[string]interface{}, referencedModels map[string]bool) {
	// Extract from responses
	if responses, ok := item["responses"].(map[string]interface{}); ok {
		for _, response := range responses {
			if responseMap, ok := response.(map[string]interface{}); ok {
				if schema, ok := responseMap["schema"].(map[string]interface{}); ok {
					extractSchemaReferences(schema, referencedModels)
				}
			}
		}
	}

	// Extract from parameters
	if parameters, ok := item["parameters"].([]interface{}); ok {
		for _, param := range parameters {
			if paramMap, ok := param.(map[string]interface{}); ok {
				if schema, ok := paramMap["schema"].(map[string]interface{}); ok {
					extractSchemaReferences(schema, referencedModels)
				}
			}
		}
	}
}

// extractSchemaReferences recursively extracts model references from schema
func extractSchemaReferences(schema map[string]interface{}, referencedModels map[string]bool) {
	// Check for direct $ref
	if ref, ok := schema["$ref"].(string); ok {
		// Extract model name from ref (e.g., "#/definitions/models.Coin" -> "models.Coin")
		if len(ref) > 14 && ref[:14] == "#/definitions/" {
			modelName := ref[14:]
			referencedModels[modelName] = true
		}
	}

	// Check for items (arrays)
	if items, ok := schema["items"].(map[string]interface{}); ok {
		extractSchemaReferences(items, referencedModels)
	}

	// Check for properties (objects)
	if properties, ok := schema["properties"].(map[string]interface{}); ok {
		for _, prop := range properties {
			if propMap, ok := prop.(map[string]interface{}); ok {
				extractSchemaReferences(propMap, referencedModels)
			}
		}
	}
}

type SwaggerHandler struct{}

func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

// GetPublicSwaggerSpec returns filtered Swagger spec for public APIs
func (h *SwaggerHandler) GetPublicSwaggerSpec(c *gin.Context) {
	tags := []string{
		"System",
		"Currency",
	}

	// Public API has no security definitions
	securityDefs := make(map[string]interface{})

	spec, err := FilterSwaggerSpec(tags, "Public API", "Public API endpoints for system health, status, and currency information. No authentication required.", securityDefs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate Swagger spec",
		})
		return
	}

	c.JSON(http.StatusOK, spec)
}

// GetPrivateSwaggerSpec returns filtered Swagger spec for private APIs
func (h *SwaggerHandler) GetPrivateSwaggerSpec(c *gin.Context) {
	tags := []string{
		"Authorization",
	}

	// Private API uses BackendSecret and BearerAuth
	securityDefs := map[string]interface{}{
		"BackendSecret": map[string]interface{}{
			"type":        "apiKey",
			"name":        "X-Backend-Secret",
			"in":          "header",
			"description": "Backend secret for API authentication (required for protected routes)",
		},
		"BearerAuth": map[string]interface{}{
			"type":        "apiKey",
			"name":        "Authorization",
			"in":          "header",
			"description": "Type \"Bearer\" followed by a space and JWT token",
		},
	}

	spec, err := FilterSwaggerSpec(tags, "Private API", "Private API endpoints for authentication, registration, and user management. Backend secret required.", securityDefs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate Swagger spec",
		})
		return
	}

	c.JSON(http.StatusOK, spec)
}

// GetPersonalSwaggerSpec returns filtered Swagger spec for personal APIs
func (h *SwaggerHandler) GetPersonalSwaggerSpec(c *gin.Context) {
	// For now, return empty spec since we don't have personal APIs yet
	// This will be populated when personal APIs are added
	tags := []string{
		"Personal", // Future tag
	}

	// Personal API uses PersonalToken
	securityDefs := map[string]interface{}{
		"PersonalToken": map[string]interface{}{
			"type":        "apiKey",
			"name":        "X-Personal-Token",
			"in":          "header",
			"description": "Personal API token for user-specific operations",
		},
	}

	spec, err := FilterSwaggerSpec(tags, "Personal API", "Personal API endpoints for user-specific data and operations. Personal token required.", securityDefs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate Swagger spec",
		})
		return
	}

	c.JSON(http.StatusOK, spec)
}
