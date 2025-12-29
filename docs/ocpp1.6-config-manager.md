# OCPP 1.6 Configuration Manager

The OCPP 1.6 Configuration Manager provides a robust solution for managing configuration keys defined in the OCPP 1.6
specification. It handles validation, custom validators, update callbacks, and ensures mandatory keys are present based
on the supported OCPP profiles.

## Overview

The configuration manager (`ManagerV16`) is designed to:

- Validate configuration keys and values according to OCPP 1.6 specifications
- Ensure all mandatory keys for selected profiles are present
- Support custom key validators for application-specific validation rules
- Trigger callbacks when configuration keys are updated
- Provide thread-safe operations for concurrent access
- Support multiple OCPP profiles (Core, LocalAuth, SmartCharging, Firmware, ISO15118, Security)

## Package Import

```go
import configManager "github.com/xBlaz3kx/ocpp-go/ocpp1.6/config_manager"
```

## Getting Started

### Creating a Configuration Manager

To create a new configuration manager, you need to:

1. Define the supported OCPP profiles
2. Create a default configuration (or use the provided defaults)
3. Initialize the manager with the configuration

```go
package main

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	configManager "github.com/lorenzodonini/ocpp-go/ocpp1.6/config_manager"
)

func main() {
	// Define supported profiles
	supportedProfiles := []string{core.ProfileName, smartcharging.ProfileName}

	// Create default configuration for the profiles
	defaultConfig, err := configManager.DefaultConfigurationFromProfiles(supportedProfiles...)
	if err != nil {
		// Handle error
		return
	}

	// Create the manager
	manager, err := configManager.NewV16ConfigurationManager(*defaultConfig, supportedProfiles...)
	if err != nil {
		// Handle error - usually means mandatory keys are missing
		return
	}

	// Manager is ready to use
}
```

### Using Default Configurations

The package provides helper functions to create default configurations for each profile:

- `DefaultConfigurationFromProfiles(profiles ...string)` - Creates a default configuration for multiple profiles
- `DefaultCoreConfiguration()` - Returns default Core profile configuration keys
- `DefaultLocalAuthConfiguration()` - Returns default LocalAuth profile configuration keys
- `DefaultSmartChargingConfiguration()` - Returns default SmartCharging profile configuration keys
- `DefaultFirmwareConfiguration()` - Returns default Firmware profile configuration keys
- `NewEmptyConfiguration()` - Creates an empty configuration (useful for custom setups)

## Core Operations

### Getting Configuration Values

Retrieve the value of a specific configuration key:

```go
value, err := manager.GetConfigurationValue(configManager.HeartbeatInterval)
if err != nil {
// Handle error (e.g., key not found)
return
}
if value != nil {
fmt.Printf("Heartbeat interval: %s\n", *value)
}
```

### Updating Configuration Keys

Update a configuration key value:

```go
newValue := "120"
err := manager.UpdateKey(configManager.HeartbeatInterval, &newValue)
if err != nil {
// Handle error - could be validation failure, readonly key, etc.
return
}
```

To remove a value (set to nil), pass `nil`:

```go
err := manager.UpdateKey(configManager.SomeOptionalKey, nil)
```

### Getting Full Configuration

Retrieve the complete configuration as OCPP ConfigurationKey array:

```go
config, err := manager.GetConfiguration()
if err != nil {
// Handle error
return
}

for _, key := range config {
fmt.Printf("Key: %s, Value: %v, Readonly: %v\n",
key.Key, key.Value, key.Readonly)
}
```

### Setting Complete Configuration

Replace the entire configuration (useful for loading from storage):

```go
newConfig := configManager.Config{
Version: 1,
Keys: []core.ConfigurationKey{
// ... your configuration keys
},
}

err := manager.SetConfiguration(newConfig)
if err != nil {
// Handle error - validation will ensure mandatory keys are present
return
}
```

## Advanced Features

### Custom Key Validators

Register a custom validator function to enforce application-specific validation rules:

```go
manager.RegisterCustomKeyValidator(func (key configManager.Key, value *string) bool {
// Custom validation logic
// Return true if valid, false otherwise

if key == configManager.HeartbeatInterval {
if value == nil {
return false // Don't allow nil values for heartbeat interval
}
// Parse and validate the interval value
interval, err := strconv.Atoi(*value)
if err != nil {
return false
}
// Ensure interval is between 60 and 3600 seconds
return interval >= 60 && interval <= 3600
}

// Allow other keys to be validated by default rules
return true
})

// Now updates will be validated using your custom validator
err := manager.UpdateKey(configManager.HeartbeatInterval, lo.ToPtr("30"))
if err != nil {
// Will fail because 30 < 60
}
```

**Note:** The custom validator is called for ALL keys. If you want to allow default validation for some keys, return
`true` for those keys in your validator function.

### Update Handlers

Register handlers to be called automatically when specific keys are updated:

```go
err := manager.OnUpdateKey(configManager.HeartbeatInterval, func(value *string) error {
if value != nil {
interval, err := strconv.Atoi(*value)
if err != nil {
return err
}
// Apply the new heartbeat interval to your system
fmt.Printf("Updating heartbeat interval to %d seconds\n", interval)
// ... your custom logic here
}
return nil
})
if err != nil {
// Handle error - key must exist in configuration
return
}

// When you update the key, the handler will be called automatically
err = manager.UpdateKey(configManager.HeartbeatInterval, lo.ToPtr("120"))
// Handler is called after successful update
```

**Important:** Update handlers are called synchronously after the key update succeeds. Keep handler logic fast to avoid
blocking other operations.

### Mandatory Keys Management

The manager automatically tracks mandatory keys based on the profiles you specify. You can:

**Get mandatory keys:**

```go
mandatoryKeys := manager.GetMandatoryKeys()
for _, key := range mandatoryKeys {
fmt.Println(key.String())
}
```

**Add additional mandatory keys:**

```go
additionalKeys := []configManager.Key{
configManager.SomeCustomKey,
}
err := manager.SetMandatoryKeys(additionalKeys)
if err != nil {
// Handle error
return
}
```

**Note:** When you set mandatory keys, they are added to the existing list. The manager will ensure all mandatory keys
are present when validating configurations.

## Supported Configuration Keys

The package provides constants for all OCPP 1.6 configuration keys:

### Core Profile Keys

- `AuthorizeRemoteTxRequests`
- `HeartbeatInterval`
- `ConnectionTimeOut`
- `MeterValueSampleInterval`
- `NumberOfConnectors`
- `SupportedFeatureProfiles`
- And many more...

### LocalAuth Profile Keys

- `LocalAuthListEnabled`
- `LocalAuthListMaxLength`
- `SendLocalListMaxLength`

### SmartCharging Profile Keys

- `ChargeProfileMaxStackLevel`
- `ChargingScheduleAllowedChargingRateUnit`
- `ChargingScheduleMaxPeriods`
- `MaxChargingProfilesInstalled`

### Firmware Profile Keys

- `SupportedFileTransferProtocols`

### ISO15118 Profile Keys

- `ISO15118PnCEnabled`
- `ContractValidationOffline`
- `CentralContractValidationAllowed`
- And more...

### Security Extension Keys

- `SecurityProfile`
- `CpoName`
- `AuthorizationData`
- And more...

See `keys.go` for the complete list of available keys.

## Error Handling

The manager returns specific errors for different failure scenarios:

- `ErrKeyCannotBeEmpty` - Attempted to use an empty key
- `ErrKeyNotFound` - Configuration key doesn't exist
- `ErrReadOnly` - Attempted to update a readonly key
- Validation errors - Custom validator rejected the value

Always check errors when performing operations:

```go
value, err := manager.GetConfigurationValue(configManager.HeartbeatInterval)
if err != nil {
switch err {
case configManager.ErrKeyNotFound:
// Key doesn't exist in configuration
case configManager.ErrKeyCannotBeEmpty:
// Invalid key provided
default:
// Other error
}
return
}
```

## Complete Example

Here's a complete example demonstrating all major features:

```go
package main

import (
	"fmt"
	"log"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/samber/lo"
	configManager "github.com/lorenzodonini/ocpp-go/ocpp1.6/config_manager"
)

func main() {
	// 1. Create default configuration
	supportedProfiles := []string{core.ProfileName, smartcharging.ProfileName}
	defaultConfig, err := configManager.DefaultConfigurationFromProfiles(supportedProfiles...)
	if err != nil {
		log.Fatalf("Failed to create default config: %v", err)
	}

	// 2. Initialize manager
	manager, err := configManager.NewV16ConfigurationManager(*defaultConfig, supportedProfiles...)
	if err != nil {
		log.Fatalf("Failed to create manager: %v", err)
	}

	// 3. Register custom validator
	manager.RegisterCustomKeyValidator(func(key configManager.Key, value *string) bool {
		// Prevent updating AuthorizeRemoteTxRequests to false
		if key == configManager.AuthorizeRemoteTxRequests && value != nil && *value == "false" {
			return false
		}
		return true
	})

	// 4. Register update handler
	err = manager.OnUpdateKey(configManager.HeartbeatInterval, func(value *string) error {
		if value != nil {
			log.Printf("Heartbeat interval updated to: %s", *value)
		}
		return nil
	})
	if err != nil {
		log.Printf("Failed to register update handler: %v", err)
	}

	// 5. Get current value
	value, err := manager.GetConfigurationValue(configManager.HeartbeatInterval)
	if err != nil {
		log.Fatalf("Failed to get value: %v", err)
	}
	fmt.Printf("Current heartbeat interval: %s\n", *value)

	// 6. Update value (handler will be called automatically)
	newValue := "180"
	err = manager.UpdateKey(configManager.HeartbeatInterval, &newValue)
	if err != nil {
		log.Fatalf("Failed to update: %v", err)
	}

	// 7. Verify update
	value, _ = manager.GetConfigurationValue(configManager.HeartbeatInterval)
	fmt.Printf("Updated heartbeat interval: %s\n", *value)

	// 8. Try to update with invalid value (will fail due to validator)
	err = manager.UpdateKey(configManager.AuthorizeRemoteTxRequests, lo.ToPtr("false"))
	if err != nil {
		fmt.Printf("Update rejected (as expected): %v\n", err)
	}

	// 9. Get full configuration
	config, err := manager.GetConfiguration()
	if err != nil {
		log.Fatalf("Failed to get configuration: %v", err)
	}
	fmt.Printf("Total configuration keys: %d\n", len(config))

	// 10. Display mandatory keys
	mandatoryKeys := manager.GetMandatoryKeys()
	fmt.Printf("Mandatory keys count: %d\n", len(mandatoryKeys))
}
```

## Best Practices

1. **Always validate errors** - The manager returns errors for invalid operations. Always check and handle them
   appropriately.

2. **Use default configurations** - Start with `DefaultConfigurationFromProfiles()` to ensure you have all required keys
   with sensible defaults.

3. **Register validators early** - Set up custom validators before allowing updates to prevent invalid configurations.

4. **Keep handlers lightweight** - Update handlers are called synchronously. Perform heavy operations asynchronously if
   needed.

5. **Handle nil values** - Configuration values can be `nil` (for optional keys). Always check before dereferencing.

6. **Use constants for keys** - Always use the provided constants (e.g., `configManager.HeartbeatInterval`) instead of
   string literals to avoid typos.

7. **Profile-specific keys** - Remember that mandatory keys depend on the profiles you support. Ensure your
   configuration includes keys for all supported profiles.

## See Also

- [OCPP 1.6 Documentation](./ocpp-1.6.md) - General OCPP 1.6 usage guide
- [Package Source Code](../ocpp1.6/config_manager/) - Full implementation details

