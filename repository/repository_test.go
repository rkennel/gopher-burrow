package repository

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

const BLANK_UUID = "00000000-0000-0000-0000-000000000000"

type DBTestSuite struct {
	suite.Suite
	Repository GiJoeRepository
	DBContext  RepositoryContext
}

func TestRunDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}

func (suite *DBTestSuite) SetupSuite() {

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	if err != nil {
		suite.FailNow("Error creating database instance")
	}

	db.AutoMigrate(&GiJoe{})

	dbContext := RepositoryContext{
		DB: db,
	}

	repository := GiJoeRepository{
		BasicRepository: NewBasicRepository[GiJoe](&dbContext),
	}
	suite.Repository = repository
	suite.DBContext = dbContext
}

func (suite *DBTestSuite) TearDownSuite() {
	db, _ := suite.DBContext.DB.DB()
	db.Close()
}

var flint = GiJoe{
	firstName: "Dashiell",
	lastName:  "Fairborne",
	codeName:  "Flint",
	jobTitle:  "Warrant Officer",
}

var ladyJaye = GiJoe{
	firstName: "Alison",
	lastName:  "Hart-Burnett",
	codeName:  "Lady Jaye",
	jobTitle:  "Covert Operations",
}

func (suite *DBTestSuite) TestWhenEntityIsInsertedBasicFieldsArePopulated() {

	insertedJoe, err := suite.Repository.Save(&flint)
	defer suite.Repository.Delete(insertedJoe.ID)

	suite.Nil(err)
	suite.NotEqual(BLANK_UUID, insertedJoe.ID.String(), "Expected ID to be generated")
	suite.NotNil(insertedJoe.CreatedAt, "Expected CreatedAt to be populated")
	suite.NotNil(insertedJoe.UpdatedAt, "Expected UpdatedAt to be populated")
}

func (suite *DBTestSuite) TestWhenEntityIsUpdatedUpdatedTimestampIsUpdated() {

	insertedJoe, err := suite.Repository.Save(&flint)
	defer suite.Repository.Delete(insertedJoe.ID)

	insertedJoe.jobTitle = "Tiger Force Warrant Officer"
	updatedJoe, err := suite.Repository.Save(insertedJoe)

	suite.Nil(err)
	suite.NotEqual(updatedJoe.UpdatedAt, updatedJoe.CreatedAt, "Updated Timestamp should be different from Created Timestamp")
}

func (suite *DBTestSuite) TestWhenEntityIsDeletedItIsNoLongerInDB() {

	insertedJoe, err := suite.Repository.Save(&flint)
	err = suite.Repository.Delete(insertedJoe.ID)
	suite.Nil(err)

	deletedJoe, err := suite.Repository.FindByID(insertedJoe.ID)
	suite.Nil(deletedJoe, "Expected deleted entity to be nil")
	suite.NotNil(err, "Expected error to be populated")
}

func (suite *DBTestSuite) TestWhenEntityIsInsertedItCanBeRetrievedById() {

	insertedJoe, err := suite.Repository.Save(&flint)
	defer suite.Repository.Delete(insertedJoe.ID)

	retrievedJoe, err := suite.Repository.FindByID(insertedJoe.ID)

	suite.Nil(err)
	suite.Equal(insertedJoe.ID, retrievedJoe.ID, "Expected retrieved entity to match inserted entity")
}

func (suite *DBTestSuite) TestWhenMultipleEntitiesAreInsertedICanRetrieveAll() {

	insertedJoe, err := suite.Repository.Save(&flint)
	defer suite.Repository.Delete(insertedJoe.ID)

	insertedJane, err := suite.Repository.Save(&ladyJaye)
	defer suite.Repository.Delete(insertedJane.ID)

	retrievedJoes, err := suite.Repository.FindAll()

	suite.Nil(err)
	suite.Equal(2, len(retrievedJoes), "Expected retrieved entities to match inserted entities")
}

func (suite *DBTestSuite) TestWhenEntityExistsThenExistsReturnsTrue() {
	insertedJoe, _ := suite.Repository.Save(&flint)
	defer suite.Repository.Delete(insertedJoe.ID)

	exists, err := suite.Repository.Exists(insertedJoe.ID)
	suite.Nil(err)
	suite.True(exists, "Expected Exists to return true")
}

func (suite *DBTestSuite) TestWhenEntityDoesNotExistThenExistsReturnsFalse() {
	exists, err := suite.Repository.Exists(uuid.New())
	suite.Equal(gorm.ErrRecordNotFound, err, "Expected Exists to return gorm.ErrRecordNotFound")
	suite.False(exists, "Expected Exists to return false")
}

func (suite *DBTestSuite) TestCountReturnsTotalNumberOfRecordsInTable() {

	insertedJoe, err := suite.Repository.Save(&flint)
	defer suite.Repository.Delete(insertedJoe.ID)

	insertedJane, err := suite.Repository.Save(&ladyJaye)
	defer suite.Repository.Delete(insertedJane.ID)

	count, err := suite.Repository.Count()

	suite.Nil(err)
	suite.Equal(int64(2), count, "Expected count of records to be 2")
}

/*
**************************************************************************
Entity and Repository for Test Purposes
**************************************************************************
*/
type GiJoe struct {
	BasicFields
	firstName string
	lastName  string
	codeName  string
	jobTitle  string
}

func (gijoe GiJoe) GetId() uuid.UUID {
	return gijoe.ID
}

type GiJoeRepository struct {
	BasicRepository[GiJoe]
}
