package controllers

import (
	"platform_api/collections"
	"platform_api/configs"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProcessController struct{
	ProcessCollection collections.ProcessCollection
}

func NewProcessController(client *mongo.Client) *ProcessController {
	return &ProcessController{ProcessCollection: *collections.NewProcessCollection(client)}
}

var processCollection *mongo.Collection = configs.OpenCollection(configs.Client, "process_engine")

// GetAllProcesses godoc
//
//	@Summary		Retrieves all processes
//	@Description	Get all the processes from the process engine
//	@Tags			processes
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.Process
//	@Failure		500	{object}	models.HTTPError
//	@Router			/processes [get]
func (t ProcessController) GetAllProcesses(c *gin.Context) {
	process, statusCode, err := t.ProcessCollection.GetAllProcesses()
	if err != nil {
		handleError(
			c,
			statusCode,
			"Failed to retrieve process",
			err,
		)
		return
	}

	c.JSON(statusCode, *process)
}

// GetProcessByCorID godoc
//
//	@Summary		Retrieves a process by Correlation ID
//	@Description	Retrieve a list of processes by their Correlation ID
//	@Tags			processes
//	@Accept			json
//	@Produce		json
//	@Param			corId	path		string	true	"Correlation ID"
//	@Success		200		{object}	models.Process
//	@Failure		400		{object}	models.HTTPError
//	@Failure		404		{object}	models.HTTPError
//	@Router			/processes/{corId} [get]
func (t ProcessController) GetProcessByCorID(c *gin.Context) {
	corId := c.Param("corId")
	
	process, statusCode, err := t.ProcessCollection.GetAllProcessByCorID(corId)
	if err != nil {
		handleError(
			c,
			statusCode,
			"Failed to retrieve process",
			err,
		)
		return
	}

	c.JSON(statusCode, *process)
}

// GetProcessByCreatorName godoc
//
//	@Summary		Retrieves processes by Creator Name
//	@Description	Retrieve a list of processes filtered by Creator Name
//	@Tags			processes
//	@Accept			json
//	@Produce		json
//	@Param			creatorName	path		string	true	"Creator's Name"
//	@Success		200			{array}		models.Process
//	@Failure		400			{object}	models.HTTPError
//	@Failure		404			{object}	models.HTTPError
//	@Router			/processes/byCreator/{creatorName} [get]
func (t ProcessController) GetProcessByCreatorName(c *gin.Context) {
	creatorName := c.Param("creatorName")
	
	process, statusCode, err := t.ProcessCollection.GetProcessByCreatorName(creatorName)
	if err != nil {
		handleError(
			c,
			statusCode,
			"Failed to retrieve process",
			err,
		)
	}

	c.JSON(statusCode, *process)
}

// GetProcessStatusByCorId godoc
//
//	@Summary		Retrieves the status of a process by Correlation ID
//	@Description	Get the most recent status of a specific process by Correlation ID
//	@Tags			processes
//	@Accept			json
//	@Produce		json
//	@Param			corId	path		string	true	"Correlation ID"
//	@Success		200		{object}	models.Process
//	@Failure		400		{object}	models.HTTPError
//	@Failure		404		{object}	models.HTTPError
//	@Failure		500		{object}	models.HTTPError
//	@Router			/processes/status/{corId} [get]
func (t ProcessController) GetProcessStatusByCorId(c *gin.Context) {
	corId := c.Param("corId")
	
	process, statusCode, err := t.ProcessCollection.GetLatestStatusByCorId(corId)
	if err != nil {
		handleError(
			c,
			statusCode,
			"Failed to retrieve process",
			err,
		)
	}

	c.JSON(statusCode, *process)
}

