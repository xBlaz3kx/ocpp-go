package types

import (
	"testing"
	"time"

	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type typesTestSuite struct {
	suite.Suite
}

func (suite *typesTestSuite) TestIdTagInfoValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{IdTagInfo{ExpiryDate: NewDateTime(time.Now()), ParentIdTag: "00000", Status: AuthorizationStatusAccepted}, true},
		{IdTagInfo{ExpiryDate: NewDateTime(time.Now()), Status: AuthorizationStatusAccepted}, true},
		{IdTagInfo{ParentIdTag: "00000", Status: AuthorizationStatusAccepted}, true},
		{IdTagInfo{Status: AuthorizationStatusAccepted}, true},
		{IdTagInfo{Status: AuthorizationStatusBlocked}, true},
		{IdTagInfo{Status: AuthorizationStatusExpired}, true},
		{IdTagInfo{Status: AuthorizationStatusInvalid}, true},
		{IdTagInfo{Status: AuthorizationStatusConcurrentTx}, true},
		{IdTagInfo{}, false},
		{IdTagInfo{Status: "invalidAuthorizationStatus"}, false},
		{IdTagInfo{ParentIdTag: ">20..................", Status: AuthorizationStatusAccepted}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *typesTestSuite) TestChargingSchedulePeriodValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{ChargingSchedulePeriod{StartPeriod: 0, Limit: 10.0, NumberPhases: tests.NewInt(3)}, true},
		{ChargingSchedulePeriod{StartPeriod: 0, Limit: 10.0}, true},
		{ChargingSchedulePeriod{StartPeriod: 0}, true},
		{ChargingSchedulePeriod{}, true},
		{ChargingSchedulePeriod{StartPeriod: -1, Limit: 10.0}, false},
		{ChargingSchedulePeriod{StartPeriod: 0, Limit: -1.0}, false},
		{ChargingSchedulePeriod{StartPeriod: 0, Limit: 10.0, NumberPhases: tests.NewInt(-1)}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *typesTestSuite) TestChargingScheduleValidation() {
	t := suite.T()
	periods := []ChargingSchedulePeriod{
		NewChargingSchedulePeriod(0, 10.0),
		NewChargingSchedulePeriod(100, 8.0),
	}
	var testTable = []tests.GenericTestEntry{
		{ChargingSchedule{Duration: tests.NewInt(0), StartSchedule: NewDateTime(time.Now()), ChargingRateUnit: ChargingRateUnitWatts, ChargingSchedulePeriod: periods, MinChargingRate: tests.NewFloat(1.0)}, true},
		{ChargingSchedule{Duration: tests.NewInt(0), ChargingRateUnit: ChargingRateUnitWatts, ChargingSchedulePeriod: periods, MinChargingRate: tests.NewFloat(1.0)}, true},
		{ChargingSchedule{ChargingRateUnit: ChargingRateUnitAmperes, ChargingSchedulePeriod: periods}, true},
		{ChargingSchedule{ChargingRateUnit: ChargingRateUnitWatts, ChargingSchedulePeriod: periods}, true},
		{ChargingSchedule{ChargingRateUnit: ChargingRateUnitWatts}, false},
		{ChargingSchedule{ChargingSchedulePeriod: periods}, false},
		{ChargingSchedule{ChargingRateUnit: ChargingRateUnitWatts, ChargingSchedulePeriod: []ChargingSchedulePeriod{}}, false},
		{ChargingSchedule{Duration: tests.NewInt(-1), ChargingRateUnit: ChargingRateUnitWatts, ChargingSchedulePeriod: periods}, false},
		{ChargingSchedule{ChargingRateUnit: ChargingRateUnitWatts, ChargingSchedulePeriod: periods, MinChargingRate: tests.NewFloat(-1.0)}, false},
		{ChargingSchedule{ChargingRateUnit: "invalidUnit", ChargingSchedulePeriod: periods}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *typesTestSuite) TestChargingProfileValidation() {
	t := suite.T()
	schedule := NewChargingSchedule(ChargingRateUnitWatts, NewChargingSchedulePeriod(0, 10.0), NewChargingSchedulePeriod(100, 8.0))
	var testTable = []tests.GenericTestEntry{
		{ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ChargingProfileKindAbsolute, RecurrencyKind: RecurrencyKindDaily, ValidFrom: NewDateTime(time.Now()), ValidTo: NewDateTime(time.Now().Add(8 * time.Hour)), ChargingSchedule: schedule}, true},
		{ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ChargingProfilePurposeTxDefaultProfile, ChargingProfileKind: ChargingProfileKindRecurring, RecurrencyKind: RecurrencyKindWeekly, ChargingSchedule: schedule}, true},
		{ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ChargingProfilePurposeTxProfile, ChargingProfileKind: ChargingProfileKindRelative, ChargingSchedule: schedule}, true},
		{ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ChargingProfileKindAbsolute, ChargingSchedule: schedule}, true},
		{ChargingProfile{StackLevel: 0, ChargingProfilePurpose: ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ChargingProfileKindAbsolute, ChargingSchedule: schedule}, true},
		{ChargingProfile{ChargingProfilePurpose: ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ChargingProfileKindAbsolute, ChargingSchedule: schedule}, true},
		{ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ChargingProfileKindAbsolute}, false},
		{ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfileKind: ChargingProfileKindAbsolute, ChargingSchedule: schedule}, false},
		{ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ChargingProfilePurposeChargePointMaxProfile, ChargingSchedule: schedule}, false},
		{ChargingProfile{ChargingProfileId: 1, StackLevel: -1, ChargingProfilePurpose: ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ChargingProfileKindAbsolute, ChargingSchedule: schedule}, false},
		{ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: "invalidPurpose", ChargingProfileKind: ChargingProfileKindAbsolute, ChargingSchedule: schedule}, false},
		{ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: "invalidKind", ChargingSchedule: schedule}, false},
		{ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ChargingProfileKindAbsolute, RecurrencyKind: "invalidRecurrencyKind", ChargingSchedule: schedule}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *typesTestSuite) TestSampledValueValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{SampledValue{Value: "42.0", Context: ReadingContextTransactionEnd, Format: ValueFormatRaw, Measurand: MeasurandPowerActiveImport, Phase: PhaseL1, Location: LocationOutlet, Unit: UnitOfMeasureWh}, true},
		{SampledValue{Value: "42.0", Context: ReadingContextTransactionEnd, Measurand: MeasurandPowerActiveImport, Phase: PhaseL1, Location: LocationOutlet}, true},
		{SampledValue{Value: "42.0", Context: ReadingContextSamplePeriodic}, true},
		{SampledValue{Value: "42.0", Measurand: MeasurandVoltage}, true},
		{SampledValue{Value: "42.0", Format: ValueFormatSignedData}, true},
		{SampledValue{Value: "42.0", Phase: PhaseL2N}, true},
		{SampledValue{Value: "42.0", Location: LocationBody}, true},
		{SampledValue{Value: "42.0", Unit: UnitOfMeasureKWh}, true},
		{SampledValue{Value: "42.0"}, true},
		{SampledValue{}, false},
		{SampledValue{Value: "42.0", Context: "invalidContext"}, false},
		{SampledValue{Value: "42.0", Format: "invalidFormat"}, false},
		{SampledValue{Value: "42.0", Measurand: "invalidMeasurand"}, false},
		{SampledValue{Value: "42.0", Phase: "invalidPhase"}, false},
		{SampledValue{Value: "42.0", Location: "invalidLocation"}, false},
		{SampledValue{Value: "42.0", Unit: "invalidUnit"}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *typesTestSuite) TestMeterValueValidation() {
	t := suite.T()
	validSample := SampledValue{Value: "42.0", Context: ReadingContextTransactionEnd, Measurand: MeasurandPowerActiveImport}
	var testTable = []tests.GenericTestEntry{
		{MeterValue{Timestamp: NewDateTime(time.Now()), SampledValue: []SampledValue{validSample}}, true},
		{MeterValue{Timestamp: NewDateTime(time.Now()), SampledValue: []SampledValue{validSample, {Value: "10.0"}}}, true},
		{MeterValue{SampledValue: []SampledValue{validSample}}, false},
		{MeterValue{Timestamp: NewDateTime(time.Now()), SampledValue: []SampledValue{}}, false},
		{MeterValue{Timestamp: NewDateTime(time.Now())}, false},
		{MeterValue{}, false},
		{MeterValue{Timestamp: NewDateTime(time.Now()), SampledValue: []SampledValue{{Value: "42.0", Context: "invalidContext"}}}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *typesTestSuite) TestStatusInfoValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{StatusInfo{ReasonCode: "okCode", AdditionalInfo: "someAdditionalInfo"}, true},
		{StatusInfo{ReasonCode: "okCode", AdditionalInfo: ""}, true},
		{StatusInfo{ReasonCode: "okCode"}, true},
		{StatusInfo{ReasonCode: ""}, false},
		{StatusInfo{}, false},
		{StatusInfo{ReasonCode: ">20.................."}, false},
		{StatusInfo{ReasonCode: "okCode", AdditionalInfo: tests.NewLongString(513)}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *typesTestSuite) TestCertificateHashDataValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{CertificateHashData{HashAlgorithm: SHA256, IssuerNameHash: "nameHash", IssuerKeyHash: "keyHash", SerialNumber: "serial"}, true},
		{CertificateHashData{HashAlgorithm: SHA384, IssuerNameHash: "nameHash", IssuerKeyHash: "keyHash", SerialNumber: "serial"}, true},
		{CertificateHashData{HashAlgorithm: SHA512, IssuerNameHash: "nameHash", IssuerKeyHash: "keyHash", SerialNumber: "serial"}, true},
		{CertificateHashData{HashAlgorithm: SHA256, IssuerNameHash: tests.NewLongString(128), IssuerKeyHash: "keyHash", SerialNumber: "serial"}, true},
		{CertificateHashData{HashAlgorithm: SHA256, IssuerNameHash: "nameHash", IssuerKeyHash: tests.NewLongString(128), SerialNumber: "serial"}, true},
		{CertificateHashData{HashAlgorithm: SHA256, IssuerNameHash: "nameHash", IssuerKeyHash: "keyHash", SerialNumber: tests.NewLongString(40)}, true},
		{CertificateHashData{}, false},
		{CertificateHashData{HashAlgorithm: "invalidAlgorithm", IssuerNameHash: "nameHash", IssuerKeyHash: "keyHash", SerialNumber: "serial"}, false},
		{CertificateHashData{HashAlgorithm: SHA256, IssuerKeyHash: "keyHash", SerialNumber: "serial"}, false},
		{CertificateHashData{HashAlgorithm: SHA256, IssuerNameHash: "nameHash", SerialNumber: "serial"}, false},
		{CertificateHashData{HashAlgorithm: SHA256, IssuerNameHash: "nameHash", IssuerKeyHash: "keyHash"}, false},
		{CertificateHashData{HashAlgorithm: SHA256, IssuerNameHash: tests.NewLongString(129), IssuerKeyHash: "keyHash", SerialNumber: "serial"}, false},
		{CertificateHashData{HashAlgorithm: SHA256, IssuerNameHash: "nameHash", IssuerKeyHash: tests.NewLongString(129), SerialNumber: "serial"}, false},
		{CertificateHashData{HashAlgorithm: SHA256, IssuerNameHash: "nameHash", IssuerKeyHash: "keyHash", SerialNumber: tests.NewLongString(41)}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func TestTypes(t *testing.T) {
	suite.Run(t, new(typesTestSuite))
}