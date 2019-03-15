package server

import (
	"bufio"
	"github.com/markusleevip/taodb/log"
	"github.com/markusleevip/taodb/util"
	"io"
	"net"
)

type Server struct{

}

func (s *Server) readKey(r *bufio.Reader) (string,error){
	klen ,err := util.ReadLen(r)
	if err !=nil{
		return "",err
	}
	k :=make([]byte,klen)
	_,err= io.ReadFull(r,k)
	if err !=nil{
		return "",err
	}
	return string(k),nil
}

func (s *Server) readAll(r *bufio.Reader) (string, []byte,error){
	klen, err := util.ReadLen(r)
	if err!=nil{
		return "",nil,err
	}
	vlen ,err:= util.ReadLen(r)
	if err!=nil{
		return "",nil,err
	}
	key :=make([]byte,klen)
	_,err = io.ReadFull(r,key)
	if err!=nil{
		return "",nil,err
	}
	value := make([]byte,vlen)
	_,err = io.ReadFull(r,value)
	if err!=nil{
		return "",nil,err
	}
	return string(key[:]),value,nil
}

func (s *Server) get(conn net.Conn, r *bufio.Reader) error{
	key,err := s.readKey(r)
	if err !=nil{
		return err
	}
	log.Info("get key=%s",key)
	value,err:=db.Get(key)
	return util.SendData(value,nil,conn)
}

func (s *Server) set(conn net.Conn, r *bufio.Reader) error{
	key,value,err := s.readAll(r)
	if err !=nil{
		return err
	}
	log.Info("get key=%s",key)
	err=db.Set(key,value)
	return util.SendData(value,err,conn)
}