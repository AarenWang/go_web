package client

import (
	"center/pb"
	"encoding/json"
	"io"
	"os"
	"strings"
)

type Service pb.Service

func (a *Service) Equals(b *Service) bool {
	if b == nil {
		return false
	}
	if a == b {
		return true
	}
	return a.Addr == b.Addr && a.Port == b.Port && a.Type == b.Type && a.Scheme == b.Scheme && a.Id == b.Id
}

type serviceDb struct {
	serviceType pb.ServiceType
	//CurrentService *Service
	List     []*Service
	filePath string
}

//func (s *serviceManager) AddAndRemoveDuplication(service *Service) {
//	list := excludeItem(s.List, service)
//	s.List = append([]*Service{service}, list...)
//}
//
//func excludeItem(list []*Service, item *Service) []*Service {
//	if len(list) > 0 {
//		for idx, item := range list {
//			if item.Equals(item) {
//				list = append(list[:idx], list[idx+1:]...)
//				break
//			}
//		}
//	}
//	return list
//}

func (s *serviceDb) LocalRead() error {
	jsonFile, err := os.Open(s.filePath)
	defer jsonFile.Close()
	if err != nil {
		return err
	}
	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	list := make([]*Service, 0)
	err = json.Unmarshal(jsonData, &list)
	if err != nil {
		return err
	}
	s.List = list
	return nil
}

func (s *serviceDb) Save() error {
	filePtr, err := os.Create(s.filePath)
	defer filePtr.Close()
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(s.List)
	if err != nil {
		return err
	}
	return nil
}

func NewServiceDb(serviceType pb.ServiceType) *serviceDb {
	fileName := strings.ToLower(pb.ServiceType_name[int32(serviceType)])
	fileName += ".json"
	m := &serviceDb{
		serviceType: serviceType,
		filePath:    fileName,
		List:        make([]*Service, 0),
	}
	m.LocalRead()
	return m
}
