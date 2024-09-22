package tests

import (
	"sync"
	"testing"
	"time"
	"vendepass/internal/dao"
	"vendepass/internal/models"

	"github.com/google/uuid"
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

	sessions.Delete(&session)

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

func TestUpdateSession(t *testing.T) {
	sessions := dao.GetSessionDAO()
	defer sessions.DeleteAll()

	session1 := &models.Session{ClientID: uuid.New(), LastTimeActive: time.Now()}
	sessions.Insert(session1)

	session2 := &models.Session{ID: session1.ID, LastTimeActive: time.Now().Add(30 * time.Minute)}
	session2.ClientID = uuid.New()

	assert.Equal(t, session1.ID, session2.ID, "expected session IDs to be equal, received %v != %v", session1.ID, session2.ID)
	assert.NotEqual(t, session1.ClientID, session2.ClientID, "expected client IDs to be different, received: %v == %v", session1.ClientID, session2.ClientID)
	assert.NotEqual(t, session1.LastTimeActive, session2.LastTimeActive, "expected time to be different, received: %v == %v", session1.LastTimeActive, session2.LastTimeActive)

	err := sessions.Update(session2)
	assert.NoError(t, err, "expected no errors, got %v", err)

	session1, err = sessions.FindById(session1.ID)
	assert.NoError(t, err, "expected no errors, got %v", err)

	assert.Equal(t, session1.ID, session2.ID, "expected session IDs to be equal, received: %v != %v", session1.ID, session2.ID)
	assert.Equal(t, session1.ClientID, session2.ClientID, "expected client IDs to be equal, received: %v != %v", session1.ClientID, session2.ClientID)
	assert.Equal(t, session1.LastTimeActive, session2.LastTimeActive, "expected time to be equal, received: %v != %v", session1.LastTimeActive, session2.LastTimeActive)
}

func TestConcurrentUpdates(t *testing.T) {
	sessions := dao.GetSessionDAO()
	defer sessions.DeleteAll()

	session1 := &models.Session{ClientID: uuid.New(), LastTimeActive: time.Now()}
	sessions.Insert(session1)

	var wg sync.WaitGroup
	var mu sync.Mutex
	startCh := make(chan struct{})
	numGoroutines := 100000
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(Id uuid.UUID) {
			defer wg.Done()

			<-startCh

			mu.Lock()
			session1.ClientID = Id
			t.Logf("%v: %v", time.Now(), session1.ClientID)
			sessions.Update(session1)
			mu.Unlock()
		}(uuid.New())
	}

	close(startCh)

	wg.Wait()

	finalState, _ := sessions.FindById(session1.ID)
	t.Logf("Final state: %s", finalState.ClientID)
}
