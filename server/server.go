package server

import (
	"bufio"
	"encoding/json"
	"github.com/markusleevip/taodb/log"
	"github.com/markusleevip/taodb/util"
	"io"
	"net"
)

type Server struct {
}

func (s *Server) readKey(r *bufio.Reader) (string, error) {
	kLen, err := util.ReadLen(r)
	if err != nil {
		return "", err
	}
	k := make([]byte, kLen)
	_, err = io.ReadFull(r, k)
	if err != nil {
		return "", err
	}
	return string(k), nil
}

func (s *Server) readAll(r *bufio.Reader) (string, []byte, error) {
	kLen, err := util.ReadLen(r)
	if err != nil {
		return "", nil, err
	}
	vLen, err := util.ReadLen(r)
	if err != nil {
		return "", nil, err
	}
	key := make([]byte, kLen)
	_, err = io.ReadFull(r, key)
	if err != nil {
		return "", nil, err
	}
	value := make([]byte, vLen)
	_, err = io.ReadFull(r, value)
	if err != nil {
		return "", nil, err
	}
	return string(key[:]), value, nil
}

func (s *Server) get(conn net.Conn, r *bufio.Reader) error {
	key, err := s.readKey(r)
	if err != nil {
		return err
	}
	value, err := db.Get(key)
	log.Info("get key=%s", key)
	return util.SendData(value, nil, conn)
}

func (s *Server) set(conn net.Conn, r *bufio.Reader) error {
	key, value, err := s.readAll(r)
	if err != nil {
		return err
	}
	log.Info("set key=%s", key)
	err = db.Set(key, value)
	return util.SendData([]byte(key), err, conn)
}

func (s *Server) del(conn net.Conn, r *bufio.Reader) error {
	key, err := s.readKey(r)
	if err != nil {
		return err
	}
	log.Info("del key=%s", key)
	err = db.Del(key)
	if err != nil {
		log.Error("del err:=%v", err)
	}
	return util.SendData([]byte(key), err, conn)
}

func (s *Server) prefix(conn net.Conn, r *bufio.Reader) error {
	key, err := s.readKey(r)
	if err != nil {
		return err
	}
	value, _ := db.Iterator(key)
	ctx, err := json.Marshal(value)
	if err!=nil{
		log.Error("prefix error:",err)
	}
	return util.SendData(ctx, err, conn)
}


func (s *Server) prefixOnlyKey(conn net.Conn, r *bufio.Reader) error {
	key, err := s.readKey(r)
	if err != nil {
		return err
	}
	value, _ := db.IteratorOnlyKey(key)
	ctx, err := json.Marshal(value)
	if err!=nil{
		log.Error("prefix error:",err)
	}
	return util.SendData(ctx, err, conn)
}

