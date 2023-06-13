package client

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type IService interface {
	SetAuthorization(token string)
	Get(route string, params map[string]string) (*resty.Response, error)
	Delete(route string, params map[string]string) (*resty.Response, error)
	Put(route string, body any, params map[string]string) (*resty.Response, error)
	Post(route string, body any, params map[string]string) (*resty.Response, error)
}

type service struct {
	client *resty.Client
	logger *zap.SugaredLogger
	token  string
}

func NewService(logger *zap.SugaredLogger, baseURL string) IService {
	return &service{logger: logger, client: resty.New().SetBaseURL(baseURL)}
}

func (s *service) send(method string, route string, params map[string]string, body any) (*resty.Response, error) {
	request := s.client.R()
	request.Method = method
	request.URL = route

	for k, v := range params {
		request.QueryParam.Add(k, v)
	}

	if s.token != "" {
		request.SetAuthToken(s.token)
	}

	if body != nil {
		request.SetHeader("Content-Type", "application/json")
		request.SetBody(body)
	}

	response, err := request.Send()

	if err != nil {
		s.logger.Errorf("issue while sending api: %v", err)
		return nil, err
	}

	s.logger.Debugf("request: %s", request.Body)
	s.logger.Infof("sent %s %s", request.Method, request.URL)
	s.logger.Debugf("response: %s", string(response.Body()))

	return response, nil
}

func (s *service) SetAuthorization(token string) {
	s.token = token
}

func (s *service) Get(route string, params map[string]string) (*resty.Response, error) {
	return s.send("GET", route, params, nil)
}

func (s *service) Delete(route string, params map[string]string) (*resty.Response, error) {
	return s.send("DELETE", route, params, nil)
}

func (s *service) Put(route string, body any, params map[string]string) (*resty.Response, error) {
	return s.send("PUT", route, params, body)
}

func (s *service) Post(route string, body any, params map[string]string) (*resty.Response, error) {
	return s.send("POST", route, params, body)
}
