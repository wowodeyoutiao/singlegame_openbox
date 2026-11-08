package global

import (
	"nanoserver/pkg/errutil"
	"reflect"
	"testing"
)

func TestGetRandomCreatCfgV2(t *testing.T) {
	// Setup: Populate the global variable with test data.
	// The slice for equipSlot 1001 is sorted by MaxLevel, as the function expects.
	mockConfig1 := &RandomCreatV2Config{EquipSlot: 1001, MaxLevel: 10}
	mockConfig2 := &RandomCreatV2Config{EquipSlot: 1001, MaxLevel: 20}
	mockConfig3 := &RandomCreatV2Config{EquipSlot: 1001, MaxLevel: 30}

	RandomCreatCfgV2 = map[int32][]*RandomCreatV2Config{
		1001: {mockConfig1, mockConfig2, mockConfig3},
		1002: {}, // Test case for an existing slot with an empty config slice
	}

	// Teardown: Clean up the global variable after the test.
	defer func() {
		RandomCreatCfgV2 = nil
	}()

	testCases := []struct {
		name          string
		equipSlot     int32
		level         int32
		expectedCfg   *RandomCreatV2Config
		expectedError error
	}{
		{
			name:          "Equip slot not found",
			equipSlot:     9999,
			level:         15,
			expectedCfg:   nil,
			expectedError: errutil.ErrConfigNotFound,
		},
		{
			name:          "Level too high, no config found",
			equipSlot:     1001,
			level:         31,
			expectedCfg:   nil,
			expectedError: errutil.ErrConfigNotFound,
		},
		{
			name:          "Level lower than all configs, should get first",
			equipSlot:     1001,
			level:         5,
			expectedCfg:   mockConfig1,
			expectedError: nil,
		},
		{
			name:          "Level matches first config exactly",
			equipSlot:     1001,
			level:         10,
			expectedCfg:   mockConfig1,
			expectedError: nil,
		},
		{
			name:          "Level between first and second config",
			equipSlot:     1001,
			level:         11,
			expectedCfg:   mockConfig2,
			expectedError: nil,
		},
		{
			name:          "Level matches second config exactly",
			equipSlot:     1001,
			level:         20,
			expectedCfg:   mockConfig2,
			expectedError: nil,
		},
		{
			name:          "Level matches last config exactly",
			equipSlot:     1001,
			level:         30,
			expectedCfg:   mockConfig3,
			expectedError: nil,
		},
		{
			name:          "Equip slot found but config slice is empty",
			equipSlot:     1002,
			level:         15,
			expectedCfg:   nil,
			expectedError: errutil.ErrConfigNotFound,
		},
		{
			name:          "Level is zero",
			equipSlot:     1001,
			level:         0,
			expectedCfg:   mockConfig1,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := GetRandomCreatCfgV2(tc.equipSlot, tc.level)

			if err != tc.expectedError {
				t.Errorf("expected error %v, but got %v", tc.expectedError, err)
			}

			if !reflect.DeepEqual(cfg, tc.expectedCfg) {
				t.Errorf("expected config %+v, but got %+v", tc.expectedCfg, cfg)
			}
		})
	}
}
