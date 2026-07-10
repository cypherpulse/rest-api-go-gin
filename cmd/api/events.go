package main

import (
	"net/http"
	"rest-api-in-gin/internal/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *application) createEvent(c *gin.Context) {
	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := app.models.Events.Insert(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}
	c.JSON(http.StatusCreated, event)
}

func (app *application) getEvents(c *gin.Context) {
	events, err := app.models.Events.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}
	c.JSON(http.StatusOK, events)
}

func (app *application) getEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := app.models.Events.Get(id)

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

func (app *application) updateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	exisitingEvent, err := app.models.Events.Get(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if exisitingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	updatedEvent := database.Event{}

	if err := c.ShouldBindJSON(&updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedEvent.Id = id

	if err := app.models.Events.Update(&updatedEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}
	c.JSON(http.StatusOK, updatedEvent)

}

func (app *application) deleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

// GetAttendeesForEvent returns all attendees for a given event
//
//	@Summary		Returns all attendees for a given event
//	@Description	Returns all attendees for a given event
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	[]database.User
//	@Router			/api/v1/events/{id}/attendees [get]
func (app *application) getAttendeesForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event id"})
		return
	}

	users, err := app.models.Attendees.GetAttendeesByEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to to retreive attendees for events"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (app *application) addAttendeeToEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event Id"})
		return
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user Id"})
		return
	}

	event, err := app.models.Events.Get(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
	}

	userToAdd, err := app.models.Users.Get(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}
	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(event.Id, userToAdd.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendee"})
		return
	}
	if existingAttendee == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User is not an attendee of this event"})
		return
	}

	attendee := database.Attendee{
		EventId: event.Id,
		UserId:  userToAdd.Id,
	}

	_, err = app.models.Attendees.Insert(&attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attendee"})
		return
	}

	c.JSON(http.StatusCreated, attendee)
}

func (app *application) getEventAttendees(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Event ID"})
		return
	}

	users, err := app.models.Attendees.GetAttendeesByEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive attendees for event"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetEventsByAttendee returns all events for a given attendee
//
//	@Summary		Returns all events for a given attendee
//	@Description	Returns all events for a given attendee
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Attendee ID"
//	@Success		200	{object}	[]database.Event
//	@Router			/api/v1/attendees/{id}/events [get]
func (app *application) getEventsByAttendee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendee id"})
		return
	}
	events, err := app.models.Attendees.GetEventsByAttendee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// DeleteAttendeeFromEvent deletes an attendee from an event
// @Summary		Deletes an attendee from an event
// @Description	Deletes an attendee from an event
// @Tags			attendees
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Event ID"
// @Param			userId	path		int	true	"User ID"
// @Success		204
// @Router			/api/v1/events/{id}/attendees/{userId} [delete]
// @Security		BearerAuth
func (app *application) deleteAttendeeFromEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event id"})
		return
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	user := app.GetUserFromContext(c)
	if event.OwnerId != user.Id {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "You are not authorized to delete an attendeeFromEvent"})
	}

	err = app.models.Attendees.Delete(userId, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendee"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
