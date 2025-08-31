# OCPP Server performance testing with k6

This directory contains comprehensive performance tests for OCPP (Open Charge Point Protocol) implementations using
Grafana k6.

## Test Scenarios

### Rapid Connect/Disconnect Test

- **Connection Lifecycle**: 100-200ms per connection
- **Message Frequency**: Random intervals between 50-150ms
- **Reconnection Pattern**: 3-8 reconnections per virtual user
- **Protocol Compliance**: Full OCPP message structure with proper JSON formatting

## Prerequisites

### Required Software

- **k6**: Performance testing framework
  ```bash
  # macOS
  brew install k6
  
  # Ubuntu/Debian
  sudo gpg -k
  sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
  echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
  sudo apt-get update
  sudo apt-get install k6
  
  # Windows
  choco install k6
  ```

## Quick Start

### 1. Clone and Navigate

```bash
cd example/k6
```

### 2. Make Script Executable

```bash
chmod +x run-performance-tests.sh
```

### 3. Run All Tests

```bash
./run-performance-tests.sh
```

### 4. Run Specific Protocol Tests

```bash
# OCPP 1.6 only
./run-performance-tests.sh --ocpp16-only

# OCPP 2.0.1 only
./run-performance-tests.sh --ocpp201-only
```

### 5. Clean Up Infrastructure

```bash
./run-performance-tests.sh --cleanup
```

## Test Configuration

### Default Test Stages

```typescript
stages: [
    {duration: '30s', target: 50},   // Ramp up to 50 VUs
    {duration: '1m', target: 50},    // Stay at 50 VUs
    {duration: '30s', target: 100},  // Ramp up to 100 VUs
    {duration: '2m', target: 100},   // Stay at 100 VUs
    {duration: '30s', target: 0},    // Ramp down to 0 VUs
]
```

### Performance Thresholds

- **Connection Time**: 95% < 200ms
- **Disconnection Time**: 95% < 100ms
- **Message Rate**: > 100 messages/second
- **WebSocket Performance**: Optimized for rapid connect/disconnect cycles

### Test Configurations

The test runner executes multiple configurations:

- **Test 1**: 25 VUs for 2 minutes
- **Test 2**: 50 VUs for 3 minutes
- **Test 3**: 100 VUs for 5 minutes

## Test Execution Details

### Message Generation

Each virtual user generates:

- **Unique IDs**: Per charge point and message
- **Realistic Data**: Random but valid OCPP message content
- **Protocol Compliance**: Proper JSON array format `[type, id, action, payload]`

### Connection Pattern

1. **Connect**: Establish WebSocket connection
2. **Boot**: Send BootNotification immediately
3. **Disconnect**: Close connection after 100-200ms
4. **Repeat**: Perform 3-8 connection cycles per VU

## Results and Analysis

The test results will be pushed to Prometheus and can be visualized in Grafana dashboards.
If running via Docker, this is handled automatically.

### Key Metrics

- **Connection Success Rate**: Percentage of successful WebSocket connections
- **Message Throughput**: Messages per second sent and received
- **Response Times**: Time to receive OCPP responses
- **Error Rates**: Connection failures and message errors
- **Resource Utilization**: CPU, memory, and network usage

### Performance Analysis

1. **Load Testing**: Verify system behavior under increasing load
2. **Stress Testing**: Identify breaking points and failure modes
3. **Stability Testing**: Ensure consistent performance over time

## Customization

### Modifying Test Parameters

Edit the test files to adjust:

- **Connection Intervals**: Change timing between connections
- **Message Types**: Add or remove OCPP message types
- **Load Patterns**: Modify VU counts and test durations
- **Thresholds**: Adjust performance expectations

### Adding New Message Types

```typescript
// Example: Add new OCPP message
function createCustomMessage() {
    const messageId = randomString(8);
    const payload = {
        // Your custom payload
    };

    return JSON.stringify([MESSAGE_TYPE.CALL, messageId, 'CustomAction', payload]);
}
```

### Environment-Specific Configuration

- **Server URLs**: Update Websocket endpoints
- **Authentication**: Add security headers if required
- **Network Settings**: Adjust timeouts and retry logic

### Debug Mode

Enable detailed logging by modifying the test script:

```typescript
// Add debug logging
console.log(`Debug: Sending message ${message}`);
console.log(`Debug: Connection state ${socket.readyState}`);
```

## Support and Resources

### Documentation

- [k6 Documentation](https://k6.io/docs/)
- [OCPP 1.6 Specification](https://www.iso.org/standard/68575.html)
- [OCPP 2.0.1 Specification](https://www.iso.org/standard/78523.html)

### Issues and Contributions

- Report bugs and feature requests
- Submit pull requests for improvements
- Share performance test results and insights