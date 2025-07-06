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

func (api *FrontendPageAPI) ListFrontendPages(c *gin.Context) {
	list := &frontendv1alpha1.FrontendPageList{}
	err := api.K8sClient.List(c.Request.Context(), list, client.InNamespace(api.Namespace))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
