package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	os "shifter/openshift"

	"github.com/gin-gonic/gin"
)

func (server *Server) SOSGetDeploymentConfigs(ctx *gin.Context) {

	// TODO: Handle Incorrect Payload. Where Shifter{} block or incorrect Shifter{} block is provided in body

	// Parse REST Request JSON Body
	var sOSDeploymentConfigs SOSDeploymentConfigs
	err := json.NewDecoder(ctx.Request.Body).Decode(&sOSDeploymentConfigs)
	if err != nil {
		fmt.Printf("error %s", err)
		ctx.JSON(501, gin.H{"error": err})
	}

	// Create OpenShift Client
	var openshift os.Openshift
	openshift.Endpoint = sOSDeploymentConfigs.Shifter.ClusterConfig.BaseUrl
	openshift.AuthToken = sOSDeploymentConfigs.Shifter.ClusterConfig.BearerToken
	openshift.Username = sOSDeploymentConfigs.Shifter.ClusterConfig.Username
	openshift.Password = sOSDeploymentConfigs.Shifter.ClusterConfig.Password

	deploymentconfigs, err := openshift.GetAllDeploymentConfigs("default")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	// Add Projects to the Response
	sOSDeploymentConfigs.DeploymentConfigs = *deploymentconfigs

	// Return JSON API Response
	ctx.JSON(http.StatusOK, sOSDeploymentConfigs)
}
