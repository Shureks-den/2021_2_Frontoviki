package repository

/*
import (
	"errors"
	"fmt"
	"log"
	"sync"
	"yula/internal/models"
	"yula/internal/pkg/session"

	"github.com/tarantool/go-tarantool"
)

type SessionRepository struct {
	pool          []*tarantool.Connection
	m             sync.RWMutex
	roundRobinCur uint32
}

func NewSessionRepository() session.SessionRepository {
	opts := tarantool.Opts{User: "guest"}
	conn, err := tarantool.Connect("158.255.163.135:3301", opts)

	if err != nil {
		log.Fatalf("Connection refused")
	}

	var pool []*tarantool.Connection
	pool = append(pool, conn)

	return &SessionRepository{
		pool:          pool,
		roundRobinCur: 0,
	}
}

func (sr *SessionRepository) AddNewConnectionToPool() error {
	conn, err := tarantool.Connect("localhost:3301", tarantool.Opts{ //  192.168.1.9
		User: "admin",
		Pass: "pass",
	})

	if err != nil {
		return errors.New("connect error")
	}

	sr.pool = append(sr.pool, conn)
	return nil
}

func (sr *SessionRepository) Set(sess *models.Session) error {

	sr.m.Lock()
	conn := sr.pool[sr.roundRobinCur]
	resp, err := conn.Insert("sessions", []interface{}{sess.Value, sess.UserId, sess.ExpiresAt.String()})
	sr.roundRobinCur = (sr.roundRobinCur + 1) % uint32(len(sr.pool))
	sr.m.Unlock()

	if err != nil {
		return errors.New("insert error")
	}

	fmt.Println(resp)
	return nil
}

func (sr *SessionRepository) Delete(sess *models.Session) error {
	sr.m.Lock()
	conn := sr.pool[sr.roundRobinCur]
	resp, err := conn.Delete("sessions", "primary", []interface{}{sess.UserId})
	sr.roundRobinCur = (sr.roundRobinCur + 1) % uint32(len(sr.pool))
	sr.m.Unlock()

	if err != nil {
		return errors.New("delete error")
	}

	fmt.Println(resp)
	return errors.New("error")
}

func (sr *SessionRepository) GetByValue(value string) (*models.Session, error) {
	sr.m.RLock()
	conn := sr.pool[sr.roundRobinCur]
	resp, err := conn.Select("sessions", "secondary", 0, 1, tarantool.IterEq, []interface{}{value})
	sr.roundRobinCur = (sr.roundRobinCur + 1) % uint32(len(sr.pool))
	sr.m.RUnlock()

	if err != nil {
		return nil, errors.New("select error")
	}

	fmt.Println(resp)
	return nil, nil
}
*/
