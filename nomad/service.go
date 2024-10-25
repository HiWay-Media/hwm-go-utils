package nomad 

import (
	"encoding/json"
    "fmt"
    
	"go.uber.org/zap"
	"github.com/go-resty/resty/v2"
)

// the nomad api client type Interface
type IService interface {
	GetDefinition(jobid, region string) (*JobDefinition, error)
	AllocationStats(allocID, region string) (*ResourceUsage, error)
	GetAllocations(clientID, region string) (*NomadAllocations, error)
}

type service struct {
	client *resty.Client
	logger *zap.SugaredLogger
}

func NewService(baseUrl string, logger *zap.SugaredLogger) IService {
	client := resty.New()
	client.SetBaseURL(baseUrl)
	return &service{
		client: client,
		logger: logger,
	}
}

func (s *service) GetDefinition(jobid, region string) (*JobDefinition, error) {
	params := map[string]string{
		"region": region,
	}
	resp, err := s.client.
		R().
		SetQueryParams(params).
		Get("/job/" + jobid)

	if err != nil {
		s.logger.Errorf("Error getting job definition: %v", err)
		return nil, err
	}
	//s.logger.Debugf("Job definition: %v", resp)
	if resp.IsError() {
		s.logger.Errorf("GetDefinition error %s", resp.String())
		return nil, fmt.Errorf("GetDefinition error %s", resp.String())
	}
	//
	var obj JobDefinition
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

func (s *service) AllocationStats(allocID, region string) (*ResourceUsage, error) {
	// 	/v1/client/allocation/:alloc_id/stats
	params := map[string]string{
		"region": region,
	}
	resp, err := s.client.
		R().
		SetQueryParams(params).
		Get("/client/allocation/" + allocID + "/stats")

	if err != nil {
		s.logger.Errorf("Error getting allocation stats: %v", err)
		return nil, err
	}
	//s.logger.Debugf("Allocation stats: %v", resp)
	if resp.IsError() {
		return nil, fmt.Errorf("AllocationStats error %s", resp.String())
	}
	//
	var obj ResourceUsage
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

func (s *service) GetAllocations(clientID, region string) (*NomadAllocations, error) {
	params := map[string]string{
		"region": region,
	}
	resp, err := s.client.
		R().
		SetQueryParams(params).
		Get("/node/" + clientID + "/allocations") 
	if err != nil {
		s.logger.Errorf("Error getting allocations: %v", err)
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("GetAllocations error %s", resp.String())
	}
	//
	var obj NomadAllocations
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}
