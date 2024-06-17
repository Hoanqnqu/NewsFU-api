package service

import outport "news-api/application/port/out"

type DummyService struct {
	dummyPort outport.Dummy
}

func NewDummyService(dummyPort outport.Dummy) *DummyService {
	return &DummyService{dummyPort: dummyPort}
}

func (d *DummyService) Get() string {
	return d.dummyPort.GetDummy()
}
