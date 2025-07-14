package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sigs.k8s.io/controller-runtime/pkg/client"

	frontendv1alpha1 "k8s-controller-tmpl/pkg/apis/frontend/v1alpha1"
)

// FrontendPageAPI provides handlers for FrontendPage resources.
type FrontendPageAPI struct {
	K8sClient client.Client
	Namespace string
}

// --- Swagger-only structs for documentation ---
// FrontendPageDoc is a simplified version for Swagger docs
// @Description FrontendPage resource (Swagger only)
type FrontendPageDoc struct {
	Name     string `json:"name" example:"example-page"`
	Contents string `json:"contents" example:"<h1>Hello</h1>"`
	Image    string `json:"image" example:"nginx:latest"`
	Replicas int    `json:"replicas" example:"2"`
}

// FrontendPageListDoc is a list of FrontendPageDoc
// @Description List of FrontendPage resources (Swagger only)
type FrontendPageListDoc struct {
	Items []FrontendPageDoc `json:"items"`
}

// ListFrontendPages godoc
// @Summary List all FrontendPages
// @Description Get all FrontendPage resources
// @Tags frontendpages
// @Produce json
// @Success 200 {object} FrontendPageListDoc
// @Router /api/frontendpages [get]
func (api *FrontendPageAPI) ListFrontendPages(c *gin.Context) {
	list := &frontendv1alpha1.FrontendPageList{}
	err := api.K8sClient.List(c.Request.Context(), list, client.InNamespace(api.Namespace))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// GetFrontendPage godoc
// @Summary Get a FrontendPage
// @Description Get a FrontendPage by name
// @Tags frontendpages
// @Produce json
// @Param name path string true "FrontendPage name"
// @Success 200 {object} FrontendPageDoc
// @Failure 404 {object} map[string]string
// @Router /api/frontendpages/{name} [get]
func (api *FrontendPageAPI) GetFrontendPage(c *gin.Context) {
	name := c.Param("name")
	page := &frontendv1alpha1.FrontendPage{}
	err := api.K8sClient.Get(c.Request.Context(), client.ObjectKey{Namespace: api.Namespace, Name: name}, page)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FrontendPage not found"})
		return
	}
	c.JSON(http.StatusOK, page)
}

// CreateFrontendPage godoc
// @Summary Create a FrontendPage
// @Description Create a new FrontendPage
// @Tags frontendpages
// @Accept json
// @Produce json
// @Param body body FrontendPageDoc true "FrontendPage object"
// @Success 201 {object} FrontendPageDoc
// @Failure 400 {object} map[string]string
// @Router /api/frontendpages [post]
func (api *FrontendPageAPI) CreateFrontendPage(c *gin.Context) {
	var page frontendv1alpha1.FrontendPage
	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	page.Namespace = api.Namespace
	if err := api.K8sClient.Create(c.Request.Context(), &page); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, page)
}

// UpdateFrontendPage godoc
// @Summary Update a FrontendPage
// @Description Update an existing FrontendPage
// @Tags frontendpages
// @Accept json
// @Produce json
// @Param name path string true "FrontendPage name"
// @Param body body FrontendPageDoc true "FrontendPage object"
// @Success 200 {object} FrontendPageDoc
// @Failure 400 {object} map[string]string
// @Router /api/frontendpages/{name} [put]
func (api *FrontendPageAPI) UpdateFrontendPage(c *gin.Context) {
	name := c.Param("name")

	existing := &frontendv1alpha1.FrontendPage{}
	err := api.K8sClient.Get(c.Request.Context(), client.ObjectKey{Namespace: api.Namespace, Name: name}, existing)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FrontendPage not found"})
		return
	}

	var patch struct {
		Spec frontendv1alpha1.FrontendPageSpec `json:"spec"`
	}

	if err := c.ShouldBindJSON(&patch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	existing.Spec = patch.Spec

	if err := api.K8sClient.Update(c.Request.Context(), existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, existing)
}

// DeleteFrontendPage godoc
// @Summary Delete a FrontendPage
// @Description Delete a FrontendPage by name
// @Tags frontendpages
// @Param name path string true "FrontendPage name"
// @Success 204 {object} nil
// @Failure 404 {object} map[string]string
// @Router /api/frontendpages/{name} [delete]
func (api *FrontendPageAPI) DeleteFrontendPage(c *gin.Context) {
	name := c.Param("name")
	page := &frontendv1alpha1.FrontendPage{}
	err := api.K8sClient.Get(c.Request.Context(), client.ObjectKey{Namespace: api.Namespace, Name: name}, page)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "FrontendPage not found"})
		return
	}

	if err := api.K8sClient.Delete(c.Request.Context(), page); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
