package tests

import (
	"testing"
	"vendepass/internal/dao"
	"vendepass/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestCreateSession(t *testing.T) {
	sessions := dao.GetSessionDAO()
	defer sessions.DeleteAll()

	session := models.Session{}

	sessions.Insert(&session)

	insertedSession, err := sessions.FindById(session.ID)

	assert.Equal(t, insertedSession.ID, session.ID, "session not equal")
	assert.NoError(t, err, "expected no error, got %v", err)
}

func TestDeleteSession(t *testing.T) {
	sessions := dao.GetSessionDAO()
	defer sessions.DeleteAll()

	session := models.Session{}

	sessions.Insert(&session)

	sessions.Delete(session)

	_, err := sessions.FindById(session.ID)

	assert.Error(t, err, "expected error, got none")
}

func TestFindSessionById(t *testing.T) {
	sessions := dao.GetSessionDAO()
	defer sessions.DeleteAll()

	session := models.Session{}

	sessions.Insert(&session)

	foundSession, err := sessions.FindById(session.ID)

	assert.Equal(t, foundSession.ID, session.ID, "session not equal")
	assert.NoError(t, err, "expected no error, got %v", err)
}

func TestFindAllSessions(t *testing.T) {
	sessions := dao.GetSessionDAO()
	defer sessions.DeleteAll()

	session1 := models.Session{}
	session2 := models.Session{}

	sessions.Insert(&session1)
	sessions.Insert(&session2)

	allSessions := sessions.FindAll()

	assert.Len(t, allSessions, 2, "expected 2 sessions, got %d", len(allSessions))
}
