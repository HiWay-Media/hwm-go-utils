package nomad 

import (
	"encoding/json"
    "fmt"
    "time"
	"go.uber.org/zap"
	"github.com/go-resty/resty/v2"
)

// the nomad api client type Interface
type IService interface {
	GetDefinition(jobid, region string) (*JobDefinition, error)
	AllocationStats(allocID, region string) (*ResourceUsage, error)
	GetAllocations(clientID, region string) (*NomadAllocations, error)
	DeleteJob(jobid, region string, purge bool) error
	ScaleJob(jobid string, count int, region string) error
	RunJob(definition JobDefinition, region string) error
	RestartJob(jobid, region string) error
}

type service struct {
	client *resty.Client
	logger *zap.SugaredLogger
}

type Options struct {
	BaseUrl  string
	LogLevel string
	Logger   *zap.SugaredLogger
}

func NewService(options Options) IService {
	client := resty.New()
	client.SetBaseURL(options.BaseUrl)
	if options.LogLevel == "debug" {
		client.SetDebug(true)
	}
	if options.Logger == nil {
		options.Logger = zap.NewNop().Sugar()
	}
	return &service{
		client: client,
		logger: options.Logger,
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


func (s *service) RestartJob(jobid, region string) error {
	err := s.ScaleJob(jobid, 0, region)
	if err != nil {
		s.logger.Errorf("Error scaling job while stopping: %v", err)
		return err
	}

	go func() {
		time.Sleep(time.Second * 1)
		err = s.ScaleJob(jobid, 1, region)
		if err != nil {
			s.logger.Errorf("Error restarting job while starting: %v", err)
		}
	}()

	return nil
}

func (s *service) DeleteJob(jobid, region string, purge bool) error {
	params := map[string]string{
		"region": region,
		"purge":  strconv.FormatBool(purge),
	}

	resp, err := s.client.
		R().
		SetQueryParams(params).
		Delete("v1/job/" + jobid)

	if err != nil {
		s.logger.Errorf("Error stopping job: %v", err)
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("StopJob error %s", resp.String())
	}

	return nil
}

func (s *service) ScaleJob(jobid string, count int, region string) error {
	params := map[string]string{
		"region": region,
	}

	request := ScaleRequest{
		Count: count,
		Target: ScaleGroupRequest{
			Group: "restreamer",
		},
	}

	resp, err := s.client.
		R().
		SetQueryParams(params).
		SetBody(request).
		SetHeader("Content-Type", "application/json").
		Post("v1/job/" + jobid + "scale")

	if err != nil {
		s.logger.Errorf("Error starting job: %v", err)
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("StartJob error %s", resp.String())
	}

	return nil
}

func (s *service) RunJob(definition JobDefinition, region string) error {
	params := map[string]string{
		"region": region,
	}

	request := RunJobRequest{
		Job:    definition,
		Format: "json",
	}

	resp, err := s.client.
		R().
		SetQueryParams(params).
		SetBody(request).
		SetHeader("content-type", "application/json").
		Post("v1/jobs")

	if err != nil {
		s.logger.Errorf("Error starting job: %v", err)
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("StartJob error %s", resp.String())
	}

	return nil
}
