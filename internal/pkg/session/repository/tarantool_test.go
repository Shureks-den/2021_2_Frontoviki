package repository

import (
	"testing"
	"time"
	"yula/internal/config"
	"yula/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	cfg := config.TarantoolConfig{
		TarantoolServerAddress: "localhost:3301",
		TarantoolOpts: config.TarantoolOptions{
			User: "admin",
			Pass: "pass",
		},
	}
	NewSessionRepository(&cfg)
}

// func TestAddToPool(t *testing.T) {
// 	cfg := config.TarantoolConfig{
// 		TarantoolServerAddress: "localhost:3031",
// 		TarantoolOpts: config.TarantoolOptions{
// 			User: "admin",
// 			Pass: "pass",
// 		},
// 	}
// 	sr := NewSessionRepository(&cfg)
// 	sr.AddNewConnectionToPool(&cfg)

// }

func TestSetAndDelete(t *testing.T) {
	cfg := config.TarantoolConfig{
		TarantoolServerAddress: "localhost:3301",
		TarantoolOpts: config.TarantoolOptions{
			User: "admin",
			Pass: "pass",
		},
	}
	sr := NewSessionRepository(&cfg)
	sess := models.Session{
		Value:     "13144441",
		UserId:    134132434,
		ExpiresAt: time.Now().Add(1000 * time.Second),
	}
	sr.Set(&sess)
	sr.Delete(&sess)

	time.Sleep(time.Second * 1)

	newSess, _ := sr.GetByValue(sess.Value)
	assert.Nil(t, newSess)

}

func TestSetGetDelete(t *testing.T) {
	cfg := config.TarantoolConfig{
		TarantoolServerAddress: "localhost:3301",
		TarantoolOpts: config.TarantoolOptions{
			User: "admin",
			Pass: "pass",
		},
	}
	sr := NewSessionRepository(&cfg)
	sess := models.Session{
		Value:     "134441",
		UserId:    1341314,
		ExpiresAt: time.Now().Add(1000 * time.Second),
	}
	sr.Set(&sess)
	sr.GetByValue(sess.Value)
	sr.Delete(&sess)

	newSess, _ := sr.GetByValue(sess.Value)
	assert.Nil(t, newSess)

}

func TestExpire(t *testing.T) {
	cfg := config.TarantoolConfig{
		TarantoolServerAddress: "localhost:3301",
		TarantoolOpts: config.TarantoolOptions{
			User: "admin",
			Pass: "pass",
		},
	}
	sr := NewSessionRepository(&cfg)
	sess := models.Session{
		Value:     "1343441",
		UserId:    13413134,
		ExpiresAt: time.Now().Add(-10 * time.Second),
	}
	sr.Set(&sess)
	time.Sleep(1 * time.Second)
	newSess, _ := sr.GetByValue(sess.Value)
	assert.Nil(t, newSess)
}
