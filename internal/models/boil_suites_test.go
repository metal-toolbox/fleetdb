// Code generated by SQLBoiler 4.15.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("AocMacAddresses", testAocMacAddresses)
	t.Run("Attributes", testAttributes)
	t.Run("AttributesFirmwareSets", testAttributesFirmwareSets)
	t.Run("BiosConfigComponents", testBiosConfigComponents)
	t.Run("BiosConfigSets", testBiosConfigSets)
	t.Run("BiosConfigSettings", testBiosConfigSettings)
	t.Run("BMCMacAddresses", testBMCMacAddresses)
	t.Run("BomInfos", testBomInfos)
	t.Run("ComponentFirmwareSets", testComponentFirmwareSets)
	t.Run("ComponentFirmwareSetMaps", testComponentFirmwareSetMaps)
	t.Run("ComponentFirmwareVersions", testComponentFirmwareVersions)
	t.Run("EventHistories", testEventHistories)
	t.Run("ServerComponentTypes", testServerComponentTypes)
	t.Run("ServerComponents", testServerComponents)
	t.Run("ServerCredentialTypes", testServerCredentialTypes)
	t.Run("ServerCredentials", testServerCredentials)
	t.Run("ServerSkus", testServerSkus)
	t.Run("ServerSkuAuxDevices", testServerSkuAuxDevices)
	t.Run("ServerSkuDisks", testServerSkuDisks)
	t.Run("ServerSkuMemories", testServerSkuMemories)
	t.Run("ServerSkuNics", testServerSkuNics)
	t.Run("Servers", testServers)
	t.Run("VersionedAttributes", testVersionedAttributes)
}

func TestSoftDelete(t *testing.T) {
	t.Run("Servers", testServersSoftDelete)
}

func TestQuerySoftDeleteAll(t *testing.T) {
	t.Run("Servers", testServersQuerySoftDeleteAll)
}

func TestSliceSoftDeleteAll(t *testing.T) {
	t.Run("Servers", testServersSliceSoftDeleteAll)
}

func TestDelete(t *testing.T) {
	t.Run("AocMacAddresses", testAocMacAddressesDelete)
	t.Run("Attributes", testAttributesDelete)
	t.Run("AttributesFirmwareSets", testAttributesFirmwareSetsDelete)
	t.Run("BiosConfigComponents", testBiosConfigComponentsDelete)
	t.Run("BiosConfigSets", testBiosConfigSetsDelete)
	t.Run("BiosConfigSettings", testBiosConfigSettingsDelete)
	t.Run("BMCMacAddresses", testBMCMacAddressesDelete)
	t.Run("BomInfos", testBomInfosDelete)
	t.Run("ComponentFirmwareSets", testComponentFirmwareSetsDelete)
	t.Run("ComponentFirmwareSetMaps", testComponentFirmwareSetMapsDelete)
	t.Run("ComponentFirmwareVersions", testComponentFirmwareVersionsDelete)
	t.Run("EventHistories", testEventHistoriesDelete)
	t.Run("ServerComponentTypes", testServerComponentTypesDelete)
	t.Run("ServerComponents", testServerComponentsDelete)
	t.Run("ServerCredentialTypes", testServerCredentialTypesDelete)
	t.Run("ServerCredentials", testServerCredentialsDelete)
	t.Run("ServerSkus", testServerSkusDelete)
	t.Run("ServerSkuAuxDevices", testServerSkuAuxDevicesDelete)
	t.Run("ServerSkuDisks", testServerSkuDisksDelete)
	t.Run("ServerSkuMemories", testServerSkuMemoriesDelete)
	t.Run("ServerSkuNics", testServerSkuNicsDelete)
	t.Run("Servers", testServersDelete)
	t.Run("VersionedAttributes", testVersionedAttributesDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("AocMacAddresses", testAocMacAddressesQueryDeleteAll)
	t.Run("Attributes", testAttributesQueryDeleteAll)
	t.Run("AttributesFirmwareSets", testAttributesFirmwareSetsQueryDeleteAll)
	t.Run("BiosConfigComponents", testBiosConfigComponentsQueryDeleteAll)
	t.Run("BiosConfigSets", testBiosConfigSetsQueryDeleteAll)
	t.Run("BiosConfigSettings", testBiosConfigSettingsQueryDeleteAll)
	t.Run("BMCMacAddresses", testBMCMacAddressesQueryDeleteAll)
	t.Run("BomInfos", testBomInfosQueryDeleteAll)
	t.Run("ComponentFirmwareSets", testComponentFirmwareSetsQueryDeleteAll)
	t.Run("ComponentFirmwareSetMaps", testComponentFirmwareSetMapsQueryDeleteAll)
	t.Run("ComponentFirmwareVersions", testComponentFirmwareVersionsQueryDeleteAll)
	t.Run("EventHistories", testEventHistoriesQueryDeleteAll)
	t.Run("ServerComponentTypes", testServerComponentTypesQueryDeleteAll)
	t.Run("ServerComponents", testServerComponentsQueryDeleteAll)
	t.Run("ServerCredentialTypes", testServerCredentialTypesQueryDeleteAll)
	t.Run("ServerCredentials", testServerCredentialsQueryDeleteAll)
	t.Run("ServerSkus", testServerSkusQueryDeleteAll)
	t.Run("ServerSkuAuxDevices", testServerSkuAuxDevicesQueryDeleteAll)
	t.Run("ServerSkuDisks", testServerSkuDisksQueryDeleteAll)
	t.Run("ServerSkuMemories", testServerSkuMemoriesQueryDeleteAll)
	t.Run("ServerSkuNics", testServerSkuNicsQueryDeleteAll)
	t.Run("Servers", testServersQueryDeleteAll)
	t.Run("VersionedAttributes", testVersionedAttributesQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("AocMacAddresses", testAocMacAddressesSliceDeleteAll)
	t.Run("Attributes", testAttributesSliceDeleteAll)
	t.Run("AttributesFirmwareSets", testAttributesFirmwareSetsSliceDeleteAll)
	t.Run("BiosConfigComponents", testBiosConfigComponentsSliceDeleteAll)
	t.Run("BiosConfigSets", testBiosConfigSetsSliceDeleteAll)
	t.Run("BiosConfigSettings", testBiosConfigSettingsSliceDeleteAll)
	t.Run("BMCMacAddresses", testBMCMacAddressesSliceDeleteAll)
	t.Run("BomInfos", testBomInfosSliceDeleteAll)
	t.Run("ComponentFirmwareSets", testComponentFirmwareSetsSliceDeleteAll)
	t.Run("ComponentFirmwareSetMaps", testComponentFirmwareSetMapsSliceDeleteAll)
	t.Run("ComponentFirmwareVersions", testComponentFirmwareVersionsSliceDeleteAll)
	t.Run("EventHistories", testEventHistoriesSliceDeleteAll)
	t.Run("ServerComponentTypes", testServerComponentTypesSliceDeleteAll)
	t.Run("ServerComponents", testServerComponentsSliceDeleteAll)
	t.Run("ServerCredentialTypes", testServerCredentialTypesSliceDeleteAll)
	t.Run("ServerCredentials", testServerCredentialsSliceDeleteAll)
	t.Run("ServerSkus", testServerSkusSliceDeleteAll)
	t.Run("ServerSkuAuxDevices", testServerSkuAuxDevicesSliceDeleteAll)
	t.Run("ServerSkuDisks", testServerSkuDisksSliceDeleteAll)
	t.Run("ServerSkuMemories", testServerSkuMemoriesSliceDeleteAll)
	t.Run("ServerSkuNics", testServerSkuNicsSliceDeleteAll)
	t.Run("Servers", testServersSliceDeleteAll)
	t.Run("VersionedAttributes", testVersionedAttributesSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("AocMacAddresses", testAocMacAddressesExists)
	t.Run("Attributes", testAttributesExists)
	t.Run("AttributesFirmwareSets", testAttributesFirmwareSetsExists)
	t.Run("BiosConfigComponents", testBiosConfigComponentsExists)
	t.Run("BiosConfigSets", testBiosConfigSetsExists)
	t.Run("BiosConfigSettings", testBiosConfigSettingsExists)
	t.Run("BMCMacAddresses", testBMCMacAddressesExists)
	t.Run("BomInfos", testBomInfosExists)
	t.Run("ComponentFirmwareSets", testComponentFirmwareSetsExists)
	t.Run("ComponentFirmwareSetMaps", testComponentFirmwareSetMapsExists)
	t.Run("ComponentFirmwareVersions", testComponentFirmwareVersionsExists)
	t.Run("EventHistories", testEventHistoriesExists)
	t.Run("ServerComponentTypes", testServerComponentTypesExists)
	t.Run("ServerComponents", testServerComponentsExists)
	t.Run("ServerCredentialTypes", testServerCredentialTypesExists)
	t.Run("ServerCredentials", testServerCredentialsExists)
	t.Run("ServerSkus", testServerSkusExists)
	t.Run("ServerSkuAuxDevices", testServerSkuAuxDevicesExists)
	t.Run("ServerSkuDisks", testServerSkuDisksExists)
	t.Run("ServerSkuMemories", testServerSkuMemoriesExists)
	t.Run("ServerSkuNics", testServerSkuNicsExists)
	t.Run("Servers", testServersExists)
	t.Run("VersionedAttributes", testVersionedAttributesExists)
}

func TestFind(t *testing.T) {
	t.Run("AocMacAddresses", testAocMacAddressesFind)
	t.Run("Attributes", testAttributesFind)
	t.Run("AttributesFirmwareSets", testAttributesFirmwareSetsFind)
	t.Run("BiosConfigComponents", testBiosConfigComponentsFind)
	t.Run("BiosConfigSets", testBiosConfigSetsFind)
	t.Run("BiosConfigSettings", testBiosConfigSettingsFind)
	t.Run("BMCMacAddresses", testBMCMacAddressesFind)
	t.Run("BomInfos", testBomInfosFind)
	t.Run("ComponentFirmwareSets", testComponentFirmwareSetsFind)
	t.Run("ComponentFirmwareSetMaps", testComponentFirmwareSetMapsFind)
	t.Run("ComponentFirmwareVersions", testComponentFirmwareVersionsFind)
	t.Run("EventHistories", testEventHistoriesFind)
	t.Run("ServerComponentTypes", testServerComponentTypesFind)
	t.Run("ServerComponents", testServerComponentsFind)
	t.Run("ServerCredentialTypes", testServerCredentialTypesFind)
	t.Run("ServerCredentials", testServerCredentialsFind)
	t.Run("ServerSkus", testServerSkusFind)
	t.Run("ServerSkuAuxDevices", testServerSkuAuxDevicesFind)
	t.Run("ServerSkuDisks", testServerSkuDisksFind)
	t.Run("ServerSkuMemories", testServerSkuMemoriesFind)
	t.Run("ServerSkuNics", testServerSkuNicsFind)
	t.Run("Servers", testServersFind)
	t.Run("VersionedAttributes", testVersionedAttributesFind)
}

func TestBind(t *testing.T) {
	t.Run("AocMacAddresses", testAocMacAddressesBind)
	t.Run("Attributes", testAttributesBind)
	t.Run("AttributesFirmwareSets", testAttributesFirmwareSetsBind)
	t.Run("BiosConfigComponents", testBiosConfigComponentsBind)
	t.Run("BiosConfigSets", testBiosConfigSetsBind)
	t.Run("BiosConfigSettings", testBiosConfigSettingsBind)
	t.Run("BMCMacAddresses", testBMCMacAddressesBind)
	t.Run("BomInfos", testBomInfosBind)
	t.Run("ComponentFirmwareSets", testComponentFirmwareSetsBind)
	t.Run("ComponentFirmwareSetMaps", testComponentFirmwareSetMapsBind)
	t.Run("ComponentFirmwareVersions", testComponentFirmwareVersionsBind)
	t.Run("EventHistories", testEventHistoriesBind)
	t.Run("ServerComponentTypes", testServerComponentTypesBind)
	t.Run("ServerComponents", testServerComponentsBind)
	t.Run("ServerCredentialTypes", testServerCredentialTypesBind)
	t.Run("ServerCredentials", testServerCredentialsBind)
	t.Run("ServerSkus", testServerSkusBind)
	t.Run("ServerSkuAuxDevices", testServerSkuAuxDevicesBind)
	t.Run("ServerSkuDisks", testServerSkuDisksBind)
	t.Run("ServerSkuMemories", testServerSkuMemoriesBind)
	t.Run("ServerSkuNics", testServerSkuNicsBind)
	t.Run("Servers", testServersBind)
	t.Run("VersionedAttributes", testVersionedAttributesBind)
}

func TestOne(t *testing.T) {
	t.Run("AocMacAddresses", testAocMacAddressesOne)
	t.Run("Attributes", testAttributesOne)
	t.Run("AttributesFirmwareSets", testAttributesFirmwareSetsOne)
	t.Run("BiosConfigComponents", testBiosConfigComponentsOne)
	t.Run("BiosConfigSets", testBiosConfigSetsOne)
	t.Run("BiosConfigSettings", testBiosConfigSettingsOne)
	t.Run("BMCMacAddresses", testBMCMacAddressesOne)
	t.Run("BomInfos", testBomInfosOne)
	t.Run("ComponentFirmwareSets", testComponentFirmwareSetsOne)
	t.Run("ComponentFirmwareSetMaps", testComponentFirmwareSetMapsOne)
	t.Run("ComponentFirmwareVersions", testComponentFirmwareVersionsOne)
	t.Run("EventHistories", testEventHistoriesOne)
	t.Run("ServerComponentTypes", testServerComponentTypesOne)
	t.Run("ServerComponents", testServerComponentsOne)
	t.Run("ServerCredentialTypes", testServerCredentialTypesOne)
	t.Run("ServerCredentials", testServerCredentialsOne)
	t.Run("ServerSkus", testServerSkusOne)
	t.Run("ServerSkuAuxDevices", testServerSkuAuxDevicesOne)
	t.Run("ServerSkuDisks", testServerSkuDisksOne)
	t.Run("ServerSkuMemories", testServerSkuMemoriesOne)
	t.Run("ServerSkuNics", testServerSkuNicsOne)
	t.Run("Servers", testServersOne)
	t.Run("VersionedAttributes", testVersionedAttributesOne)
}

func TestAll(t *testing.T) {
	t.Run("AocMacAddresses", testAocMacAddressesAll)
	t.Run("Attributes", testAttributesAll)
	t.Run("AttributesFirmwareSets", testAttributesFirmwareSetsAll)
	t.Run("BiosConfigComponents", testBiosConfigComponentsAll)
	t.Run("BiosConfigSets", testBiosConfigSetsAll)
	t.Run("BiosConfigSettings", testBiosConfigSettingsAll)
	t.Run("BMCMacAddresses", testBMCMacAddressesAll)
	t.Run("BomInfos", testBomInfosAll)
	t.Run("ComponentFirmwareSets", testComponentFirmwareSetsAll)
	t.Run("ComponentFirmwareSetMaps", testComponentFirmwareSetMapsAll)
	t.Run("ComponentFirmwareVersions", testComponentFirmwareVersionsAll)
	t.Run("EventHistories", testEventHistoriesAll)
	t.Run("ServerComponentTypes", testServerComponentTypesAll)
	t.Run("ServerComponents", testServerComponentsAll)
	t.Run("ServerCredentialTypes", testServerCredentialTypesAll)
	t.Run("ServerCredentials", testServerCredentialsAll)
	t.Run("ServerSkus", testServerSkusAll)
	t.Run("ServerSkuAuxDevices", testServerSkuAuxDevicesAll)
	t.Run("ServerSkuDisks", testServerSkuDisksAll)
	t.Run("ServerSkuMemories", testServerSkuMemoriesAll)
	t.Run("ServerSkuNics", testServerSkuNicsAll)
	t.Run("Servers", testServersAll)
	t.Run("VersionedAttributes", testVersionedAttributesAll)
}

func TestCount(t *testing.T) {
	t.Run("AocMacAddresses", testAocMacAddressesCount)
	t.Run("Attributes", testAttributesCount)
	t.Run("AttributesFirmwareSets", testAttributesFirmwareSetsCount)
	t.Run("BiosConfigComponents", testBiosConfigComponentsCount)
	t.Run("BiosConfigSets", testBiosConfigSetsCount)
	t.Run("BiosConfigSettings", testBiosConfigSettingsCount)
	t.Run("BMCMacAddresses", testBMCMacAddressesCount)
	t.Run("BomInfos", testBomInfosCount)
	t.Run("ComponentFirmwareSets", testComponentFirmwareSetsCount)
	t.Run("ComponentFirmwareSetMaps", testComponentFirmwareSetMapsCount)
	t.Run("ComponentFirmwareVersions", testComponentFirmwareVersionsCount)
	t.Run("EventHistories", testEventHistoriesCount)
	t.Run("ServerComponentTypes", testServerComponentTypesCount)
	t.Run("ServerComponents", testServerComponentsCount)
	t.Run("ServerCredentialTypes", testServerCredentialTypesCount)
	t.Run("ServerCredentials", testServerCredentialsCount)
	t.Run("ServerSkus", testServerSkusCount)
	t.Run("ServerSkuAuxDevices", testServerSkuAuxDevicesCount)
	t.Run("ServerSkuDisks", testServerSkuDisksCount)
	t.Run("ServerSkuMemories", testServerSkuMemoriesCount)
	t.Run("ServerSkuNics", testServerSkuNicsCount)
	t.Run("Servers", testServersCount)
	t.Run("VersionedAttributes", testVersionedAttributesCount)
}

func TestHooks(t *testing.T) {
	t.Run("AocMacAddresses", testAocMacAddressesHooks)
	t.Run("Attributes", testAttributesHooks)
	t.Run("AttributesFirmwareSets", testAttributesFirmwareSetsHooks)
	t.Run("BiosConfigComponents", testBiosConfigComponentsHooks)
	t.Run("BiosConfigSets", testBiosConfigSetsHooks)
	t.Run("BiosConfigSettings", testBiosConfigSettingsHooks)
	t.Run("BMCMacAddresses", testBMCMacAddressesHooks)
	t.Run("BomInfos", testBomInfosHooks)
	t.Run("ComponentFirmwareSets", testComponentFirmwareSetsHooks)
	t.Run("ComponentFirmwareSetMaps", testComponentFirmwareSetMapsHooks)
	t.Run("ComponentFirmwareVersions", testComponentFirmwareVersionsHooks)
	t.Run("EventHistories", testEventHistoriesHooks)
	t.Run("ServerComponentTypes", testServerComponentTypesHooks)
	t.Run("ServerComponents", testServerComponentsHooks)
	t.Run("ServerCredentialTypes", testServerCredentialTypesHooks)
	t.Run("ServerCredentials", testServerCredentialsHooks)
	t.Run("ServerSkus", testServerSkusHooks)
	t.Run("ServerSkuAuxDevices", testServerSkuAuxDevicesHooks)
	t.Run("ServerSkuDisks", testServerSkuDisksHooks)
	t.Run("ServerSkuMemories", testServerSkuMemoriesHooks)
	t.Run("ServerSkuNics", testServerSkuNicsHooks)
	t.Run("Servers", testServersHooks)
	t.Run("VersionedAttributes", testVersionedAttributesHooks)
}

func TestInsert(t *testing.T) {
	t.Run("AocMacAddresses", testAocMacAddressesInsert)
	t.Run("AocMacAddresses", testAocMacAddressesInsertWhitelist)
	t.Run("Attributes", testAttributesInsert)
	t.Run("Attributes", testAttributesInsertWhitelist)
	t.Run("AttributesFirmwareSets", testAttributesFirmwareSetsInsert)
	t.Run("AttributesFirmwareSets", testAttributesFirmwareSetsInsertWhitelist)
	t.Run("BiosConfigComponents", testBiosConfigComponentsInsert)
	t.Run("BiosConfigComponents", testBiosConfigComponentsInsertWhitelist)
	t.Run("BiosConfigSets", testBiosConfigSetsInsert)
	t.Run("BiosConfigSets", testBiosConfigSetsInsertWhitelist)
	t.Run("BiosConfigSettings", testBiosConfigSettingsInsert)
	t.Run("BiosConfigSettings", testBiosConfigSettingsInsertWhitelist)
	t.Run("BMCMacAddresses", testBMCMacAddressesInsert)
	t.Run("BMCMacAddresses", testBMCMacAddressesInsertWhitelist)
	t.Run("BomInfos", testBomInfosInsert)
	t.Run("BomInfos", testBomInfosInsertWhitelist)
	t.Run("ComponentFirmwareSets", testComponentFirmwareSetsInsert)
	t.Run("ComponentFirmwareSets", testComponentFirmwareSetsInsertWhitelist)
	t.Run("ComponentFirmwareSetMaps", testComponentFirmwareSetMapsInsert)
	t.Run("ComponentFirmwareSetMaps", testComponentFirmwareSetMapsInsertWhitelist)
	t.Run("ComponentFirmwareVersions", testComponentFirmwareVersionsInsert)
	t.Run("ComponentFirmwareVersions", testComponentFirmwareVersionsInsertWhitelist)
	t.Run("EventHistories", testEventHistoriesInsert)
	t.Run("EventHistories", testEventHistoriesInsertWhitelist)
	t.Run("ServerComponentTypes", testServerComponentTypesInsert)
	t.Run("ServerComponentTypes", testServerComponentTypesInsertWhitelist)
	t.Run("ServerComponents", testServerComponentsInsert)
	t.Run("ServerComponents", testServerComponentsInsertWhitelist)
	t.Run("ServerCredentialTypes", testServerCredentialTypesInsert)
	t.Run("ServerCredentialTypes", testServerCredentialTypesInsertWhitelist)
	t.Run("ServerCredentials", testServerCredentialsInsert)
	t.Run("ServerCredentials", testServerCredentialsInsertWhitelist)
	t.Run("ServerSkus", testServerSkusInsert)
	t.Run("ServerSkus", testServerSkusInsertWhitelist)
	t.Run("ServerSkuAuxDevices", testServerSkuAuxDevicesInsert)
	t.Run("ServerSkuAuxDevices", testServerSkuAuxDevicesInsertWhitelist)
	t.Run("ServerSkuDisks", testServerSkuDisksInsert)
	t.Run("ServerSkuDisks", testServerSkuDisksInsertWhitelist)
	t.Run("ServerSkuMemories", testServerSkuMemoriesInsert)
	t.Run("ServerSkuMemories", testServerSkuMemoriesInsertWhitelist)
	t.Run("ServerSkuNics", testServerSkuNicsInsert)
	t.Run("ServerSkuNics", testServerSkuNicsInsertWhitelist)
	t.Run("Servers", testServersInsert)
	t.Run("Servers", testServersInsertWhitelist)
	t.Run("VersionedAttributes", testVersionedAttributesInsert)
	t.Run("VersionedAttributes", testVersionedAttributesInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("AocMacAddressToBomInfoUsingSerialNumBomInfo", testAocMacAddressToOneBomInfoUsingSerialNumBomInfo)
	t.Run("AttributeToServerUsingServer", testAttributeToOneServerUsingServer)
	t.Run("AttributeToServerComponentUsingServerComponent", testAttributeToOneServerComponentUsingServerComponent)
	t.Run("AttributesFirmwareSetToComponentFirmwareSetUsingFirmwareSet", testAttributesFirmwareSetToOneComponentFirmwareSetUsingFirmwareSet)
	t.Run("BiosConfigComponentToBiosConfigSetUsingFKBiosConfigSet", testBiosConfigComponentToOneBiosConfigSetUsingFKBiosConfigSet)
	t.Run("BiosConfigSettingToBiosConfigComponentUsingFKBiosConfigComponent", testBiosConfigSettingToOneBiosConfigComponentUsingFKBiosConfigComponent)
	t.Run("BMCMacAddressToBomInfoUsingSerialNumBomInfo", testBMCMacAddressToOneBomInfoUsingSerialNumBomInfo)
	t.Run("ComponentFirmwareSetMapToComponentFirmwareSetUsingFirmwareSet", testComponentFirmwareSetMapToOneComponentFirmwareSetUsingFirmwareSet)
	t.Run("ComponentFirmwareSetMapToComponentFirmwareVersionUsingFirmware", testComponentFirmwareSetMapToOneComponentFirmwareVersionUsingFirmware)
	t.Run("EventHistoryToServerUsingTargetServerServer", testEventHistoryToOneServerUsingTargetServerServer)
	t.Run("ServerComponentToServerUsingServer", testServerComponentToOneServerUsingServer)
	t.Run("ServerComponentToServerComponentTypeUsingServerComponentType", testServerComponentToOneServerComponentTypeUsingServerComponentType)
	t.Run("ServerCredentialToServerCredentialTypeUsingServerCredentialType", testServerCredentialToOneServerCredentialTypeUsingServerCredentialType)
	t.Run("ServerCredentialToServerUsingServer", testServerCredentialToOneServerUsingServer)
	t.Run("ServerSkuAuxDeviceToServerSkuUsingSku", testServerSkuAuxDeviceToOneServerSkuUsingSku)
	t.Run("ServerSkuDiskToServerSkuUsingSku", testServerSkuDiskToOneServerSkuUsingSku)
	t.Run("ServerSkuMemoryToServerSkuUsingSku", testServerSkuMemoryToOneServerSkuUsingSku)
	t.Run("ServerSkuNicToServerSkuUsingSku", testServerSkuNicToOneServerSkuUsingSku)
	t.Run("VersionedAttributeToServerUsingServer", testVersionedAttributeToOneServerUsingServer)
	t.Run("VersionedAttributeToServerComponentUsingServerComponent", testVersionedAttributeToOneServerComponentUsingServerComponent)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("BiosConfigComponentToFKBiosConfigComponentBiosConfigSettings", testBiosConfigComponentToManyFKBiosConfigComponentBiosConfigSettings)
	t.Run("BiosConfigSetToFKBiosConfigSetBiosConfigComponents", testBiosConfigSetToManyFKBiosConfigSetBiosConfigComponents)
	t.Run("BomInfoToSerialNumAocMacAddresses", testBomInfoToManySerialNumAocMacAddresses)
	t.Run("BomInfoToSerialNumBMCMacAddresses", testBomInfoToManySerialNumBMCMacAddresses)
	t.Run("ComponentFirmwareSetToFirmwareSetAttributesFirmwareSets", testComponentFirmwareSetToManyFirmwareSetAttributesFirmwareSets)
	t.Run("ComponentFirmwareSetToFirmwareSetComponentFirmwareSetMaps", testComponentFirmwareSetToManyFirmwareSetComponentFirmwareSetMaps)
	t.Run("ComponentFirmwareVersionToFirmwareComponentFirmwareSetMaps", testComponentFirmwareVersionToManyFirmwareComponentFirmwareSetMaps)
	t.Run("ServerComponentTypeToServerComponents", testServerComponentTypeToManyServerComponents)
	t.Run("ServerComponentToAttributes", testServerComponentToManyAttributes)
	t.Run("ServerComponentToVersionedAttributes", testServerComponentToManyVersionedAttributes)
	t.Run("ServerCredentialTypeToServerCredentials", testServerCredentialTypeToManyServerCredentials)
	t.Run("ServerSkuToSkuServerSkuAuxDevices", testServerSkuToManySkuServerSkuAuxDevices)
	t.Run("ServerSkuToSkuServerSkuDisks", testServerSkuToManySkuServerSkuDisks)
	t.Run("ServerSkuToSkuServerSkuMemories", testServerSkuToManySkuServerSkuMemories)
	t.Run("ServerSkuToSkuServerSkuNics", testServerSkuToManySkuServerSkuNics)
	t.Run("ServerToAttributes", testServerToManyAttributes)
	t.Run("ServerToTargetServerEventHistories", testServerToManyTargetServerEventHistories)
	t.Run("ServerToServerComponents", testServerToManyServerComponents)
	t.Run("ServerToServerCredentials", testServerToManyServerCredentials)
	t.Run("ServerToVersionedAttributes", testServerToManyVersionedAttributes)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("AocMacAddressToBomInfoUsingSerialNumAocMacAddresses", testAocMacAddressToOneSetOpBomInfoUsingSerialNumBomInfo)
	t.Run("AttributeToServerUsingAttributes", testAttributeToOneSetOpServerUsingServer)
	t.Run("AttributeToServerComponentUsingAttributes", testAttributeToOneSetOpServerComponentUsingServerComponent)
	t.Run("AttributesFirmwareSetToComponentFirmwareSetUsingFirmwareSetAttributesFirmwareSets", testAttributesFirmwareSetToOneSetOpComponentFirmwareSetUsingFirmwareSet)
	t.Run("BiosConfigComponentToBiosConfigSetUsingFKBiosConfigSetBiosConfigComponents", testBiosConfigComponentToOneSetOpBiosConfigSetUsingFKBiosConfigSet)
	t.Run("BiosConfigSettingToBiosConfigComponentUsingFKBiosConfigComponentBiosConfigSettings", testBiosConfigSettingToOneSetOpBiosConfigComponentUsingFKBiosConfigComponent)
	t.Run("BMCMacAddressToBomInfoUsingSerialNumBMCMacAddresses", testBMCMacAddressToOneSetOpBomInfoUsingSerialNumBomInfo)
	t.Run("ComponentFirmwareSetMapToComponentFirmwareSetUsingFirmwareSetComponentFirmwareSetMaps", testComponentFirmwareSetMapToOneSetOpComponentFirmwareSetUsingFirmwareSet)
	t.Run("ComponentFirmwareSetMapToComponentFirmwareVersionUsingFirmwareComponentFirmwareSetMaps", testComponentFirmwareSetMapToOneSetOpComponentFirmwareVersionUsingFirmware)
	t.Run("EventHistoryToServerUsingTargetServerEventHistories", testEventHistoryToOneSetOpServerUsingTargetServerServer)
	t.Run("ServerComponentToServerUsingServerComponents", testServerComponentToOneSetOpServerUsingServer)
	t.Run("ServerComponentToServerComponentTypeUsingServerComponents", testServerComponentToOneSetOpServerComponentTypeUsingServerComponentType)
	t.Run("ServerCredentialToServerCredentialTypeUsingServerCredentials", testServerCredentialToOneSetOpServerCredentialTypeUsingServerCredentialType)
	t.Run("ServerCredentialToServerUsingServerCredentials", testServerCredentialToOneSetOpServerUsingServer)
	t.Run("ServerSkuAuxDeviceToServerSkuUsingSkuServerSkuAuxDevices", testServerSkuAuxDeviceToOneSetOpServerSkuUsingSku)
	t.Run("ServerSkuDiskToServerSkuUsingSkuServerSkuDisks", testServerSkuDiskToOneSetOpServerSkuUsingSku)
	t.Run("ServerSkuMemoryToServerSkuUsingSkuServerSkuMemories", testServerSkuMemoryToOneSetOpServerSkuUsingSku)
	t.Run("ServerSkuNicToServerSkuUsingSkuServerSkuNics", testServerSkuNicToOneSetOpServerSkuUsingSku)
	t.Run("VersionedAttributeToServerUsingVersionedAttributes", testVersionedAttributeToOneSetOpServerUsingServer)
	t.Run("VersionedAttributeToServerComponentUsingVersionedAttributes", testVersionedAttributeToOneSetOpServerComponentUsingServerComponent)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {
	t.Run("AttributeToServerUsingAttributes", testAttributeToOneRemoveOpServerUsingServer)
	t.Run("AttributeToServerComponentUsingAttributes", testAttributeToOneRemoveOpServerComponentUsingServerComponent)
	t.Run("AttributesFirmwareSetToComponentFirmwareSetUsingFirmwareSetAttributesFirmwareSets", testAttributesFirmwareSetToOneRemoveOpComponentFirmwareSetUsingFirmwareSet)
	t.Run("VersionedAttributeToServerUsingVersionedAttributes", testVersionedAttributeToOneRemoveOpServerUsingServer)
	t.Run("VersionedAttributeToServerComponentUsingVersionedAttributes", testVersionedAttributeToOneRemoveOpServerComponentUsingServerComponent)
}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
	t.Run("BiosConfigComponentToFKBiosConfigComponentBiosConfigSettings", testBiosConfigComponentToManyAddOpFKBiosConfigComponentBiosConfigSettings)
	t.Run("BiosConfigSetToFKBiosConfigSetBiosConfigComponents", testBiosConfigSetToManyAddOpFKBiosConfigSetBiosConfigComponents)
	t.Run("BomInfoToSerialNumAocMacAddresses", testBomInfoToManyAddOpSerialNumAocMacAddresses)
	t.Run("BomInfoToSerialNumBMCMacAddresses", testBomInfoToManyAddOpSerialNumBMCMacAddresses)
	t.Run("ComponentFirmwareSetToFirmwareSetAttributesFirmwareSets", testComponentFirmwareSetToManyAddOpFirmwareSetAttributesFirmwareSets)
	t.Run("ComponentFirmwareSetToFirmwareSetComponentFirmwareSetMaps", testComponentFirmwareSetToManyAddOpFirmwareSetComponentFirmwareSetMaps)
	t.Run("ComponentFirmwareVersionToFirmwareComponentFirmwareSetMaps", testComponentFirmwareVersionToManyAddOpFirmwareComponentFirmwareSetMaps)
	t.Run("ServerComponentTypeToServerComponents", testServerComponentTypeToManyAddOpServerComponents)
	t.Run("ServerComponentToAttributes", testServerComponentToManyAddOpAttributes)
	t.Run("ServerComponentToVersionedAttributes", testServerComponentToManyAddOpVersionedAttributes)
	t.Run("ServerCredentialTypeToServerCredentials", testServerCredentialTypeToManyAddOpServerCredentials)
	t.Run("ServerSkuToSkuServerSkuAuxDevices", testServerSkuToManyAddOpSkuServerSkuAuxDevices)
	t.Run("ServerSkuToSkuServerSkuDisks", testServerSkuToManyAddOpSkuServerSkuDisks)
	t.Run("ServerSkuToSkuServerSkuMemories", testServerSkuToManyAddOpSkuServerSkuMemories)
	t.Run("ServerSkuToSkuServerSkuNics", testServerSkuToManyAddOpSkuServerSkuNics)
	t.Run("ServerToAttributes", testServerToManyAddOpAttributes)
	t.Run("ServerToTargetServerEventHistories", testServerToManyAddOpTargetServerEventHistories)
	t.Run("ServerToServerComponents", testServerToManyAddOpServerComponents)
	t.Run("ServerToServerCredentials", testServerToManyAddOpServerCredentials)
	t.Run("ServerToVersionedAttributes", testServerToManyAddOpVersionedAttributes)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {
	t.Run("ComponentFirmwareSetToFirmwareSetAttributesFirmwareSets", testComponentFirmwareSetToManySetOpFirmwareSetAttributesFirmwareSets)
	t.Run("ServerComponentToAttributes", testServerComponentToManySetOpAttributes)
	t.Run("ServerComponentToVersionedAttributes", testServerComponentToManySetOpVersionedAttributes)
	t.Run("ServerToAttributes", testServerToManySetOpAttributes)
	t.Run("ServerToVersionedAttributes", testServerToManySetOpVersionedAttributes)
}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {
	t.Run("ComponentFirmwareSetToFirmwareSetAttributesFirmwareSets", testComponentFirmwareSetToManyRemoveOpFirmwareSetAttributesFirmwareSets)
	t.Run("ServerComponentToAttributes", testServerComponentToManyRemoveOpAttributes)
	t.Run("ServerComponentToVersionedAttributes", testServerComponentToManyRemoveOpVersionedAttributes)
	t.Run("ServerToAttributes", testServerToManyRemoveOpAttributes)
	t.Run("ServerToVersionedAttributes", testServerToManyRemoveOpVersionedAttributes)
}

func TestReload(t *testing.T) {
	t.Run("AocMacAddresses", testAocMacAddressesReload)
	t.Run("Attributes", testAttributesReload)
	t.Run("AttributesFirmwareSets", testAttributesFirmwareSetsReload)
	t.Run("BiosConfigComponents", testBiosConfigComponentsReload)
	t.Run("BiosConfigSets", testBiosConfigSetsReload)
	t.Run("BiosConfigSettings", testBiosConfigSettingsReload)
	t.Run("BMCMacAddresses", testBMCMacAddressesReload)
	t.Run("BomInfos", testBomInfosReload)
	t.Run("ComponentFirmwareSets", testComponentFirmwareSetsReload)
	t.Run("ComponentFirmwareSetMaps", testComponentFirmwareSetMapsReload)
	t.Run("ComponentFirmwareVersions", testComponentFirmwareVersionsReload)
	t.Run("EventHistories", testEventHistoriesReload)
	t.Run("ServerComponentTypes", testServerComponentTypesReload)
	t.Run("ServerComponents", testServerComponentsReload)
	t.Run("ServerCredentialTypes", testServerCredentialTypesReload)
	t.Run("ServerCredentials", testServerCredentialsReload)
	t.Run("ServerSkus", testServerSkusReload)
	t.Run("ServerSkuAuxDevices", testServerSkuAuxDevicesReload)
	t.Run("ServerSkuDisks", testServerSkuDisksReload)
	t.Run("ServerSkuMemories", testServerSkuMemoriesReload)
	t.Run("ServerSkuNics", testServerSkuNicsReload)
	t.Run("Servers", testServersReload)
	t.Run("VersionedAttributes", testVersionedAttributesReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("AocMacAddresses", testAocMacAddressesReloadAll)
	t.Run("Attributes", testAttributesReloadAll)
	t.Run("AttributesFirmwareSets", testAttributesFirmwareSetsReloadAll)
	t.Run("BiosConfigComponents", testBiosConfigComponentsReloadAll)
	t.Run("BiosConfigSets", testBiosConfigSetsReloadAll)
	t.Run("BiosConfigSettings", testBiosConfigSettingsReloadAll)
	t.Run("BMCMacAddresses", testBMCMacAddressesReloadAll)
	t.Run("BomInfos", testBomInfosReloadAll)
	t.Run("ComponentFirmwareSets", testComponentFirmwareSetsReloadAll)
	t.Run("ComponentFirmwareSetMaps", testComponentFirmwareSetMapsReloadAll)
	t.Run("ComponentFirmwareVersions", testComponentFirmwareVersionsReloadAll)
	t.Run("EventHistories", testEventHistoriesReloadAll)
	t.Run("ServerComponentTypes", testServerComponentTypesReloadAll)
	t.Run("ServerComponents", testServerComponentsReloadAll)
	t.Run("ServerCredentialTypes", testServerCredentialTypesReloadAll)
	t.Run("ServerCredentials", testServerCredentialsReloadAll)
	t.Run("ServerSkus", testServerSkusReloadAll)
	t.Run("ServerSkuAuxDevices", testServerSkuAuxDevicesReloadAll)
	t.Run("ServerSkuDisks", testServerSkuDisksReloadAll)
	t.Run("ServerSkuMemories", testServerSkuMemoriesReloadAll)
	t.Run("ServerSkuNics", testServerSkuNicsReloadAll)
	t.Run("Servers", testServersReloadAll)
	t.Run("VersionedAttributes", testVersionedAttributesReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("AocMacAddresses", testAocMacAddressesSelect)
	t.Run("Attributes", testAttributesSelect)
	t.Run("AttributesFirmwareSets", testAttributesFirmwareSetsSelect)
	t.Run("BiosConfigComponents", testBiosConfigComponentsSelect)
	t.Run("BiosConfigSets", testBiosConfigSetsSelect)
	t.Run("BiosConfigSettings", testBiosConfigSettingsSelect)
	t.Run("BMCMacAddresses", testBMCMacAddressesSelect)
	t.Run("BomInfos", testBomInfosSelect)
	t.Run("ComponentFirmwareSets", testComponentFirmwareSetsSelect)
	t.Run("ComponentFirmwareSetMaps", testComponentFirmwareSetMapsSelect)
	t.Run("ComponentFirmwareVersions", testComponentFirmwareVersionsSelect)
	t.Run("EventHistories", testEventHistoriesSelect)
	t.Run("ServerComponentTypes", testServerComponentTypesSelect)
	t.Run("ServerComponents", testServerComponentsSelect)
	t.Run("ServerCredentialTypes", testServerCredentialTypesSelect)
	t.Run("ServerCredentials", testServerCredentialsSelect)
	t.Run("ServerSkus", testServerSkusSelect)
	t.Run("ServerSkuAuxDevices", testServerSkuAuxDevicesSelect)
	t.Run("ServerSkuDisks", testServerSkuDisksSelect)
	t.Run("ServerSkuMemories", testServerSkuMemoriesSelect)
	t.Run("ServerSkuNics", testServerSkuNicsSelect)
	t.Run("Servers", testServersSelect)
	t.Run("VersionedAttributes", testVersionedAttributesSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("AocMacAddresses", testAocMacAddressesUpdate)
	t.Run("Attributes", testAttributesUpdate)
	t.Run("AttributesFirmwareSets", testAttributesFirmwareSetsUpdate)
	t.Run("BiosConfigComponents", testBiosConfigComponentsUpdate)
	t.Run("BiosConfigSets", testBiosConfigSetsUpdate)
	t.Run("BiosConfigSettings", testBiosConfigSettingsUpdate)
	t.Run("BMCMacAddresses", testBMCMacAddressesUpdate)
	t.Run("BomInfos", testBomInfosUpdate)
	t.Run("ComponentFirmwareSets", testComponentFirmwareSetsUpdate)
	t.Run("ComponentFirmwareSetMaps", testComponentFirmwareSetMapsUpdate)
	t.Run("ComponentFirmwareVersions", testComponentFirmwareVersionsUpdate)
	t.Run("EventHistories", testEventHistoriesUpdate)
	t.Run("ServerComponentTypes", testServerComponentTypesUpdate)
	t.Run("ServerComponents", testServerComponentsUpdate)
	t.Run("ServerCredentialTypes", testServerCredentialTypesUpdate)
	t.Run("ServerCredentials", testServerCredentialsUpdate)
	t.Run("ServerSkus", testServerSkusUpdate)
	t.Run("ServerSkuAuxDevices", testServerSkuAuxDevicesUpdate)
	t.Run("ServerSkuDisks", testServerSkuDisksUpdate)
	t.Run("ServerSkuMemories", testServerSkuMemoriesUpdate)
	t.Run("ServerSkuNics", testServerSkuNicsUpdate)
	t.Run("Servers", testServersUpdate)
	t.Run("VersionedAttributes", testVersionedAttributesUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("AocMacAddresses", testAocMacAddressesSliceUpdateAll)
	t.Run("Attributes", testAttributesSliceUpdateAll)
	t.Run("AttributesFirmwareSets", testAttributesFirmwareSetsSliceUpdateAll)
	t.Run("BiosConfigComponents", testBiosConfigComponentsSliceUpdateAll)
	t.Run("BiosConfigSets", testBiosConfigSetsSliceUpdateAll)
	t.Run("BiosConfigSettings", testBiosConfigSettingsSliceUpdateAll)
	t.Run("BMCMacAddresses", testBMCMacAddressesSliceUpdateAll)
	t.Run("BomInfos", testBomInfosSliceUpdateAll)
	t.Run("ComponentFirmwareSets", testComponentFirmwareSetsSliceUpdateAll)
	t.Run("ComponentFirmwareSetMaps", testComponentFirmwareSetMapsSliceUpdateAll)
	t.Run("ComponentFirmwareVersions", testComponentFirmwareVersionsSliceUpdateAll)
	t.Run("EventHistories", testEventHistoriesSliceUpdateAll)
	t.Run("ServerComponentTypes", testServerComponentTypesSliceUpdateAll)
	t.Run("ServerComponents", testServerComponentsSliceUpdateAll)
	t.Run("ServerCredentialTypes", testServerCredentialTypesSliceUpdateAll)
	t.Run("ServerCredentials", testServerCredentialsSliceUpdateAll)
	t.Run("ServerSkus", testServerSkusSliceUpdateAll)
	t.Run("ServerSkuAuxDevices", testServerSkuAuxDevicesSliceUpdateAll)
	t.Run("ServerSkuDisks", testServerSkuDisksSliceUpdateAll)
	t.Run("ServerSkuMemories", testServerSkuMemoriesSliceUpdateAll)
	t.Run("ServerSkuNics", testServerSkuNicsSliceUpdateAll)
	t.Run("Servers", testServersSliceUpdateAll)
	t.Run("VersionedAttributes", testVersionedAttributesSliceUpdateAll)
}
