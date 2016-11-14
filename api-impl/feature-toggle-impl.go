package feature_toggle_impl

import (
	"fmt"

	api "github.com/peterrosell/feature-toggle-service/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/peterrosell/feature-toggle-service/storage"
	"github.com/pkg/errors"
	"github.com/peterrosell/feature-toggle-service/featuretree"
)

type FeatureToggleServiceServer struct {
	fs storage.FeatureToggleStore
	tree *featuretree.ToggleRuleTree
}
func (s *FeatureToggleServiceServer) GetFeaturesForProperties(ctx context.Context, req *api.GetFeaturesByPropertiesRequest) (*api.GetFeaturesByPropertiesResponse, error){
	if s.tree == nil {
		return nil, errors.New("Feature toggle service not initialized.")
	}
	return &api.GetFeaturesByPropertiesResponse{Features:s.tree.FindFeatures(req.Properties)}, nil
}

func (s *FeatureToggleServiceServer) CreateToggleRule(ctx context.Context, req *api.CreateToggleRuleRequest) (*api.CreateToggleRuleResponse, error) {
	fmt.Printf("CreateToggleRule: %v\n", req.ToggleRule)

	feature, err := s.fs.ReadFeatureByName(req.ToggleRule.Name)
	if err != nil {
		return nil, errors.New( "Unknown feature")
	}

	propAsSlice := []string{}
	for k,v := range req.ToggleRule.Properties {
		propAsSlice = append(propAsSlice, k, v)
	}

	ruleId, err := s.fs.CreateToggleRule(*storage.NewToggleRule((*feature).Id, req.ToggleRule.Enabled, propAsSlice...))
	if err != nil {
		return nil, err
	}
	response := new(api.CreateToggleRuleResponse)
	response.Id = *ruleId

	return response, nil
}

func (s *FeatureToggleServiceServer) ReadToggleRule(ctx context.Context, req *api.ReadToggleRuleRequest) (*api.ReadToggleRuleResponse, error) {
	fmt.Printf("ReadToggleRule: id=%s\n", req.Id)

	toggleRule := new(api.ToggleRule)
	toggleRule.Id = req.Id
	toggleRule.Name = "extra speed"
	response := new(api.ReadToggleRuleResponse)
	response.ToggleRule = toggleRule

	return response, nil
}

func (s *FeatureToggleServiceServer) DeleteToggleRule(ctx context.Context, req *api.DeleteToggleRuleRequest) (*api.DeleteToggleRuleResponse, error) {
	fmt.Printf("DeleteToggleRule: id=%s\n", req.Id)
	return new(api.DeleteToggleRuleResponse), nil
}

func (s *FeatureToggleServiceServer) SearchToggleRule(ctx context.Context, req *api.SearchToggleRuleRequest) (*api.SearchToggleRuleResponse, error) {
	fmt.Printf("SearchToggleRule: %s\n", req)
	toggleRule := new(api.ToggleRule)
	toggleRule.Name = "Smörrebröd"

	//_, err := s.fs.SearchToggleRule(req.Name, make(Filter))
	//if err != nil {
	//	return nil, errors.New( "failed to search for rules")
	//}

	response := new(api.SearchToggleRuleResponse)
	//response.ToggleRules = []*api.ToggleRule{ToApiToggleRules(rules)}

	return response, nil
}
/*
func ToApiToggleRules(rules *[]storage.ToggleRule) []api.ToggleRule{
	toggleRules := make([]api.ToggleRule, len(rules))
	for i := 0; i < len(rules); i++ {
		toggleRules[i] = ToApiToggleRule( rules[0])
	}
}

func ToApiToggleRule(rule storage.ToggleRule) api.ToggleRule {
	return api.ToggleRule{Id:rule.Id,Name:rule.FeatureId}
}
*/

func (s *FeatureToggleServiceServer) CreateFeature(ctx context.Context, req *api.CreateFeatureRequest) (*api.CreateFeatureResponse, error) {
	fmt.Printf("CreateFeature: %v\n", req.Feature)
	fmt.Printf("CreateFeature: id=%s\n", req.Feature.Name)

	featureId, err := s.fs.CreateFeature(*storage.NewFeature(req.Feature.Name, req.Feature.Enabled, req.Feature.Description))
	if err != nil {
		return nil, err
	}
	response := new(api.CreateFeatureResponse)
	response.Id = *featureId

	return response, nil
}

func (s *FeatureToggleServiceServer) ReadFeature(ctx context.Context, req *api.ReadFeatureRequest) (*api.ReadFeatureResponse, error) {
	fmt.Printf("ReadFeature: id=%s\n", req.Id)

	feature := new(api.Feature)
	feature.Id = req.Id
	response := new(api.ReadFeatureResponse)
	response.Feature = feature

	return response, nil
}

func (s *FeatureToggleServiceServer) DeleteFeature(ctx context.Context, req *api.DeleteFeatureRequest) (*api.DeleteFeatureResponse, error) {
	fmt.Printf("DeleteFeature: id=%s\n", req.Id)
	return new(api.DeleteFeatureResponse), nil
}

func (s *FeatureToggleServiceServer) SearchFeature(ctx context.Context, req *api.SearchFeatureRequest) (*api.SearchFeatureResponse, error) {
	/*
	if( req.Filter != nil) {
		fmt.Printf("SearchFeature: %s\n", req.Filter.Name)
	} else {
		fmt.Printf("SearchFeature filter=nil: %v\n", req)
	}
	*/
	fmt.Printf("SearchFeature name=: %v\n", req.Name)
	feature := new(api.Feature)
	feature.Name = "Smörrebröd"

	response := new(api.SearchFeatureResponse)
	response.Features = []*api.Feature{feature}

	return response, nil
}

func (s *FeatureToggleServiceServer) CreateProperty(ctx context.Context, req *api.CreatePropertyRequest) (*api.CreatePropertyResponse, error) {
	fmt.Printf("CreateProperty: %v\n", req.Property)
	fmt.Printf("CreateProperty: id=%s\n", req.Property.Name)

	propertyId, err := s.fs.CreateProperty(*storage.NewProperty(req.Property.Name, req.Property.Description))
	if err != nil {
		return nil, err
	}

	response := new(api.CreatePropertyResponse)
	response.Name = *propertyId

	return response, nil
}

func (s *FeatureToggleServiceServer) ReadProperty(ctx context.Context, req *api.ReadPropertyRequest) (*api.ReadPropertyResponse, error) {
	fmt.Printf("ReadProperty: id=%s\n", req.Name)

	property := new(api.Property)
	property.Name = req.Name
	response := new(api.ReadPropertyResponse)
	response.Property = property

	return response, nil
}

func (s *FeatureToggleServiceServer) DeleteProperty(ctx context.Context, req *api.DeletePropertyRequest) (*api.DeletePropertyResponse, error) {
	fmt.Printf("DeleteProperty: id=%s\n", req.Name)
	return new(api.DeletePropertyResponse), nil
}

func (s *FeatureToggleServiceServer) SearchProperty(ctx context.Context, req *api.SearchPropertyRequest) (*api.SearchPropertyResponse, error) {
	fmt.Printf("SearchFeature: %s\n", req.Name)
	property := new(api.Property)
	property.Name = "usertype"

	response := new(api.SearchPropertyResponse)
	response.Properties = []*api.Property{property}

	return response, nil
}

func newFeatureToggleServiceServer() *FeatureToggleServiceServer {
	s := new(FeatureToggleServiceServer)
	s.fs = storage.NewFeatureToggleStoreImpl()
	s.fs.Open()

	toggleRules, err := s.fs.GetEnabledToggleRules()
	if( toggleRules != nil) {
		//properties := s.fs.ReadAllPropertyNames();
		propertyNames := make([]string,2)
		propertyNames[0] = "usertype"
		propertyNames[1] = "userid"
		tree := featuretree.NewFeatureTree(propertyNames)

		for _, rule := range *toggleRules {
			tree.AddFeature(rule)
		}
		s.tree = tree
	} else {
		fmt.Printf("Failed to init feature toggle service, %v\n", err)
	}
	return s
}

func RegisterFeatureToggleService(s *grpc.Server) {
	api.RegisterFeatureToggleServiceServer(s, newFeatureToggleServiceServer())
}
