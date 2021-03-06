package amboy

import (
	"errors"
	"strings"
	"testing"

	"github.com/mongodb/amboy/dependency"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BaseCheckSuite struct {
	base    *JobBase
	require *require.Assertions
	suite.Suite
}

func TestBaseCheckSuite(t *testing.T) {
	suite.Run(t, new(BaseCheckSuite))
}

func (s *BaseCheckSuite) SetupSuite() {
	s.require = s.Require()
}

func (s *BaseCheckSuite) SetupTest() {
	s.base = &JobBase{dep: dependency.NewAlways()}
}

func (s *BaseCheckSuite) TestInitialValuesOfBaseObject() {
	s.False(s.base.IsComplete)
	s.Len(s.base.Errors, 0)
}

func (s *BaseCheckSuite) TestAddErrorWithNilObjectDoesNotChangeErrorState() {
	for i := 0; i < 100; i++ {
		s.base.AddError(nil)
		s.NoError(s.base.Error())
		s.Len(s.base.Errors, 0)
		s.False(s.base.HasErrors())
	}
}

func (s *BaseCheckSuite) TestAddErrorsPersistsErrorsInJob() {
	for i := 1; i <= 100; i++ {
		s.base.AddError(errors.New("foo"))
		s.Error(s.base.Error())
		s.Len(s.base.Errors, i)
		s.True(s.base.HasErrors())
		s.Len(strings.Split(s.base.Error().Error(), "\n"), i)
	}
}

func (s *BaseCheckSuite) TestIdIsAccessorForTaskIDAttribute() {
	s.Equal(s.base.TaskID, s.base.ID())
	s.base.TaskID = "foo"
	s.Equal("foo", s.base.ID())
	s.Equal(s.base.TaskID, s.base.ID())
}

func (s *BaseCheckSuite) TestDependencyAccessorIsCorrect() {
	s.Equal(s.base.dep, s.base.Dependency())
	s.base.SetDependency(dependency.NewAlways())
	s.Equal(dependency.AlwaysRun, s.base.Dependency().Type().Name)
}

func (s *BaseCheckSuite) TestSetDependencyAcceptsAndPersistsChangesToDependencyType() {
	s.Equal(dependency.AlwaysRun, s.base.dep.Type().Name)
	localDep := dependency.NewLocalFileInstance()
	s.NotEqual(localDep.Type().Name, dependency.AlwaysRun)
	s.base.SetDependency(localDep)
	s.Equal(dependency.LocalFileRelationship, s.base.dep.Type().Name)
}

func (s *BaseCheckSuite) TestMarkCompleteHelperSetsCompleteState() {
	s.False(s.base.IsComplete)
	s.False(s.base.Completed())
	s.base.MarkComplete()

	s.True(s.base.IsComplete)
	s.True(s.base.Completed())
}

func (s *BaseCheckSuite) TestRoundTripAbilityThroughImportAndExport() {
	s.base.JobType = JobType{
		Format:  JSON,
		Version: 42,
	}

	s.Equal(42, s.base.JobType.Version)
	out, err := s.base.Export()

	s.NoError(err)
	s.NotNil(out)

	s.base.JobType.Version = 21
	s.Equal(21, s.base.JobType.Version)

	err = s.base.Import(out)
	s.NoError(err)
	s.Equal(42, s.base.JobType.Version)

}
