import {randomIntBetween, randomString} from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import ws from 'k6/ws';
import {sleep, check} from 'k6';

// OCPP 1.6 Core profile message structure constants
const MESSAGE_TYPE = {
    CALL: 2,
    CALL_RESULT: 3,
    CALL_ERROR: 4
};

const ACTIONS = {
    BOOT_NOTIFICATION: 'BootNotification',
    HEARTBEAT: 'Heartbeat',
    STATUS_NOTIFICATION: 'StatusNotification',
    AUTHORIZE: 'Authorize',
    START_TRANSACTION: 'StartTransaction',
    STOP_TRANSACTION: 'StopTransaction',
    METER_VALUES: 'MeterValues',
    FIRMWARE_STATUS_NOTIFICATION: 'FirmwareStatusNotification',
    DIAGNOSTICS_STATUS_NOTIFICATION: 'DiagnosticsStatusNotification',
    DATA_TRANSFER: 'DataTransfer'
};

// Core profile incoming request actions from central system
const INCOMING_ACTIONS = {
    CHANGE_AVAILABILITY: 'ChangeAvailability',
    CHANGE_CONFIGURATION: 'ChangeConfiguration',
    CLEAR_CACHE: 'ClearCache',
    DATA_TRANSFER: 'DataTransfer',
    GET_CONFIGURATION: 'GetConfiguration',
    REMOTE_START_TRANSACTION: 'RemoteStartTransaction',
    REMOTE_STOP_TRANSACTION: 'RemoteStopTransaction',
    RESET: 'Reset',
    UNLOCK_CONNECTOR: 'UnlockConnector',
    GET_DIAGNOSTICS: 'GetDiagnostics',
    UPDATE_FIRMWARE: 'UpdateFirmware',
    GET_LOCAL_LIST_VERSION: 'GetLocalListVersion',
    SEND_LOCAL_LIST: 'SendLocalList',
    RESERVE_NOW: 'ReserveNow',
    CANCEL_RESERVATION: 'CancelReservation',
    SET_CHARGING_PROFILE: 'SetChargingProfile',
    CLEAR_CHARGING_PROFILE: 'ClearChargingProfile',
    GET_COMPOSITE_SCHEDULE: 'GetCompositeSchedule',
    TRIGGER_MESSAGE: 'TriggerMessage'
};

// Generate unique charge point ID
function generateChargePointId() {
    return `CP_${__VU}_${randomString(8)}`;
}

// Generate OCPP 1.6 Core profile BootNotification message
function createBootNotification(chargePointId: string) {
    const messageId = randomString(8);
    const payload = {
        chargePointModel: `Model_${randomString(5)}`,
        chargePointVendor: `Vendor_${randomString(5)}`,
        chargePointSerialNumber: chargePointId,
        firmwareVersion: `v${randomIntBetween(1, 9)}.${randomIntBetween(0, 9)}.${randomIntBetween(0, 9)}`,
        iccid: randomString(20),
        imsi: randomString(15)
    };

    return JSON.stringify([MESSAGE_TYPE.CALL, messageId, ACTIONS.BOOT_NOTIFICATION, payload]);
}

// Generate OCPP 1.6 Core profile Heartbeat message
function createHeartbeat() {
    const messageId = randomString(8);
    return JSON.stringify([MESSAGE_TYPE.CALL, messageId, ACTIONS.HEARTBEAT, {}]);
}

// Generate OCPP 1.6 Core profile StatusNotification message
function createStatusNotification(chargePointId: string, connectorId: number, status: string) {
    const messageId = randomString(8);
    const payload = {
        connectorId: connectorId,
        errorCode: 'NoError',
        status: status,
        timestamp: new Date().toISOString()
    };

    return JSON.stringify([MESSAGE_TYPE.CALL, messageId, ACTIONS.STATUS_NOTIFICATION, payload]);
}

// Generate OCPP 1.6 Core profile Authorize message
function createAuthorize(idTag: string) {
    const messageId = randomString(8);
    const payload = {
        idTag: idTag
    };

    return JSON.stringify([MESSAGE_TYPE.CALL, messageId, ACTIONS.AUTHORIZE, payload]);
}

// Generate OCPP 1.6 Core profile StartTransaction message
function createStartTransaction(connectorId: number, idTag: string) {
    const messageId = randomString(8);
    const payload = {
        connectorId: connectorId,
        idTag: idTag,
        meterStart: randomIntBetween(0, 1000),
        timestamp: new Date().toISOString()
    };

    return JSON.stringify([MESSAGE_TYPE.CALL, messageId, ACTIONS.START_TRANSACTION, payload]);
}

// Generate OCPP 1.6 Core profile StopTransaction message
function createStopTransaction(transactionId: number, meterStop: number) {
    const messageId = randomString(8);
    const payload = {
        transactionId: transactionId,
        meterStop: meterStop,
        timestamp: new Date().toISOString(),
        reason: 'Remote'
    };

    return JSON.stringify([MESSAGE_TYPE.CALL, messageId, ACTIONS.STOP_TRANSACTION, payload]);
}

// Generate OCPP 1.6 Core profile MeterValues message
function createMeterValues(connectorId: number, transactionId?: number) {
    const messageId = randomString(8);
    const payload = {
        connectorId: connectorId,
        transactionId: transactionId,
        meterValue: [{
            timestamp: new Date().toISOString(),
            sampledValue: [{
                value: randomIntBetween(0, 100000).toString(),
                context: 'Sample.Periodic',
                measurand: 'Energy.Active.Import.Register',
                unit: 'Wh'
            }]
        }]
    };

    return JSON.stringify([MESSAGE_TYPE.CALL, messageId, ACTIONS.METER_VALUES, payload]);
}

// Generate OCPP 1.6 Core profile CallResult response
function createCallResult(messageId: string, payload: any) {
    return JSON.stringify([MESSAGE_TYPE.CALL_RESULT, messageId, payload]);
}

// Generate OCPP 1.6 Core profile CallError response
function createCallError(messageId: string, errorCode: string, errorDescription: string) {
    return JSON.stringify([MESSAGE_TYPE.CALL_ERROR, messageId, errorCode, errorDescription, {}]);
}

// Generate OCPP 1.6 Core profile FirmwareStatusNotification message
function createFirmwareStatusNotification(status: string) {
    const messageId = randomString(8);
    const payload = {
        status: status
    };

    return JSON.stringify([MESSAGE_TYPE.CALL, messageId, ACTIONS.FIRMWARE_STATUS_NOTIFICATION, payload]);
}

// Generate OCPP 1.6 Core profile DiagnosticsStatusNotification message
function createDiagnosticsStatusNotification(status: string) {
    const messageId = randomString(8);
    const payload = {
        status: status
    };

    return JSON.stringify([MESSAGE_TYPE.CALL, messageId, ACTIONS.DIAGNOSTICS_STATUS_NOTIFICATION, payload]);
}

// Generate OCPP 1.6 Core profile DataTransfer message (outgoing)
function createDataTransfer(vendorId: string, messageId?: string, data?: string) {
    const msgId = randomString(8);
    const payload = {
        vendorId: vendorId,
        messageId: messageId,
        data: data
    };

    return JSON.stringify([MESSAGE_TYPE.CALL, msgId, ACTIONS.DATA_TRANSFER, payload]);
}

// Handle incoming Core profile ChangeAvailability request
function handleChangeAvailability(request: any) {
    return {
        status: 'Accepted'
    };
}

// Handle incoming Core profile ChangeConfiguration request
function handleChangeConfiguration(request: any) {
    return {
        status: 'Accepted'
    };
}

// Handle incoming Core profile ClearCache request
function handleClearCache(request: any) {
    return {
        status: 'Accepted'
    };
}

// Handle incoming Core profile DataTransfer request
function handleDataTransfer(request: any) {
    return {
        status: 'Accepted',
        data: `Response to ${request.vendorId}:${request.messageId || 'default'}`
    };
}

// Handle incoming Core profile GetConfiguration request
function handleGetConfiguration(request: any) {
    return {
        configurationKey: [
            {
                key: 'HeartbeatInterval',
                readonly: false,
                value: '60'
            },
            {
                key: 'ConnectionTimeOut',
                readonly: false,
                value: '60'
            }
        ],
        unknownKey: []
    };
}

// Handle incoming Core profile RemoteStartTransaction request
function handleRemoteStartTransaction(request: any) {
    return {
        status: 'Accepted'
    };
}

// Handle incoming Core profile RemoteStopTransaction request
function handleRemoteStopTransaction(request: any) {
    return {
        status: 'Accepted'
    };
}

// Handle incoming Core profile Reset request
function handleReset(request: any) {
    return {
        status: 'Accepted'
    };
}

// Handle incoming Core profile UnlockConnector request
function handleUnlockConnector(request: any) {
    return {
        status: 'Accepted'
    };
}

// Handle incoming Core profile GetDiagnostics request
function handleGetDiagnostics(request: any) {
    return {
        fileName: `diagnostics_${Date.now()}.zip`
    };
}

// Handle incoming Core profile UpdateFirmware request
function handleUpdateFirmware(request: any) {
    return {
        status: 'Accepted'
    };
}

// Handle incoming Core profile GetLocalListVersion request
function handleGetLocalListVersion(request: any) {
    return {
        listVersion: randomIntBetween(1, 100)
    };
}

// Handle incoming Core profile SendLocalList request
function handleSendLocalList(request: any) {
    return {
        status: 'Accepted'
    };
}

// Handle incoming Core profile ReserveNow request
function handleReserveNow(request: any) {
    return {
        status: 'Accepted'
    };
}

// Handle incoming Core profile CancelReservation request
function handleCancelReservation(request: any) {
    return {
        status: 'Accepted'
    };
}

// Handle incoming Core profile SetChargingProfile request
function handleSetChargingProfile(request: any) {
    return {
        status: 'Accepted'
    };
}

// Handle incoming Core profile ClearChargingProfile request
function handleClearChargingProfile(request: any) {
    return {
        status: 'Accepted'
    };
}

// Handle incoming Core profile GetCompositeSchedule request
function handleGetCompositeSchedule(request: any) {
    return {
        status: 'Accepted',
        connectorId: request.connectorId,
        scheduleStart: new Date(Date.now() + 60000).toISOString(),
        chargingSchedule: {
            duration: 3600,
            startSchedule: new Date().toISOString(),
            chargingRateUnit: 'A',
            chargingSchedulePeriod: [
                {
                    startPeriod: 0,
                    limit: 32
                }
            ]
        }
    };
}

// Handle incoming Core profile TriggerMessage request
function handleTriggerMessage(request: any) {
    return {
        status: 'Accepted'
    };
}

// Route incoming Core profile requests to appropriate handlers
function handleIncomingRequest(action: string, payload: any, messageId: string) {
    let response;

    try {
        switch (action) {
            case INCOMING_ACTIONS.CHANGE_AVAILABILITY:
                response = handleChangeAvailability(payload);
                break;
            case INCOMING_ACTIONS.CHANGE_CONFIGURATION:
                response = handleChangeConfiguration(payload);
                break;
            case INCOMING_ACTIONS.CLEAR_CACHE:
                response = handleClearCache(payload);
                break;
            case INCOMING_ACTIONS.DATA_TRANSFER:
                response = handleDataTransfer(payload);
                break;
            case INCOMING_ACTIONS.GET_CONFIGURATION:
                response = handleGetConfiguration(payload);
                break;
            case INCOMING_ACTIONS.REMOTE_START_TRANSACTION:
                response = handleRemoteStartTransaction(payload);
                break;
            case INCOMING_ACTIONS.REMOTE_STOP_TRANSACTION:
                response = handleRemoteStopTransaction(payload);
                break;
            case INCOMING_ACTIONS.RESET:
                response = handleReset(payload);
                break;
            case INCOMING_ACTIONS.UNLOCK_CONNECTOR:
                response = handleUnlockConnector(payload);
                break;
            case INCOMING_ACTIONS.GET_DIAGNOSTICS:
                response = handleGetDiagnostics(payload);
                break;
            case INCOMING_ACTIONS.UPDATE_FIRMWARE:
                response = handleUpdateFirmware(payload);
                break;
            case INCOMING_ACTIONS.GET_LOCAL_LIST_VERSION:
                response = handleGetLocalListVersion(payload);
                break;
            case INCOMING_ACTIONS.SEND_LOCAL_LIST:
                response = handleSendLocalList(payload);
                break;
            case INCOMING_ACTIONS.RESERVE_NOW:
                response = handleReserveNow(payload);
                break;
            case INCOMING_ACTIONS.CANCEL_RESERVATION:
                response = handleCancelReservation(payload);
                break;
            case INCOMING_ACTIONS.SET_CHARGING_PROFILE:
                response = handleSetChargingProfile(payload);
                break;
            case INCOMING_ACTIONS.CLEAR_CHARGING_PROFILE:
                response = handleClearChargingProfile(payload);
                break;
            case INCOMING_ACTIONS.GET_COMPOSITE_SCHEDULE:
                response = handleGetCompositeSchedule(payload);
                break;
            case INCOMING_ACTIONS.TRIGGER_MESSAGE:
                response = handleTriggerMessage(payload);
                break;
            default:
                console.log(`VU ${__VU}: Unknown incoming action: ${action}`);
                return createCallError(messageId, 'NotImplemented', `Action ${action} not implemented`);
        }

        return createCallResult(messageId, response);
    } catch (error) {
        console.log(`VU ${__VU}: Error handling ${action}: ${error}`);
        return createCallError(messageId, 'InternalError', `Failed to process ${action}`);
    }
}

export default function () {
    const chargePointId = generateChargePointId();

    // Configurable WebSocket URL from environment variables
    const wsHost = __ENV.WS_HOST || 'central-system';
    const wsPort = __ENV.WS_PORT || '8887';
    const wsProtocol = __ENV.WS_PROTOCOL || 'ws';
    const wsPath = __ENV.WS_PATH || '';

    const url = `${wsProtocol}://${wsHost}:${wsPort}${wsPath}/${chargePointId}`;

    console.log(`VU ${__VU}: Connecting to ${url}`);

    const params = {
        // OCPP 1.6 requires specific WebSocket subprotocol
        headers: {
            'Sec-WebSocket-Protocol': 'ocpp1.6'
        }
    };

    // Track connection metrics
    const startTime = Date.now();
    let messageCount = 0;
    let reconnectCount = 0;
    let maxReconnects = 10;

    while (reconnectCount < maxReconnects) {
        const res = ws.connect(url, params, function (socket) {
            const connectionStart = Date.now();

            socket.on('open', function open() {
                const connectTime = Date.now() - connectionStart;
                console.log(`VU ${__VU}: Connected in ${connectTime}ms (attempt ${reconnectCount + 1}/${maxReconnects})`);

                // Send BootNotification immediately after connection
                const bootMsg = createBootNotification(chargePointId);
                socket.send(bootMsg);
                messageCount++;

                // Disconnect after random time (100-200ms)
                const disconnectDelay = randomIntBetween(100, 200);

                socket.setInterval(function timeout() {
                    socket.ping();
                    console.log('Pinging');
                }, disconnectDelay - 10);

                socket.setTimeout(function () {
                    console.log(`${disconnectDelay} ms passed, closing the socket`);
                    socket.close();
                }, disconnectDelay);
            });

            socket.on('message', function (message) {
                try {
                    const data = JSON.parse(message);

                    if (data[0] === MESSAGE_TYPE.CALL_RESULT) {
                        console.log(`VU ${__VU}: Received response for message ${data[1]}`);
                    } else if (data[0] === MESSAGE_TYPE.CALL_ERROR) {
                        console.log(`VU ${__VU}: Received error for message ${data[1]}: ${data[2]}`);
                    } else if (data[0] === MESSAGE_TYPE.CALL) {
                        // Handle incoming request from central system
                        const messageId = data[1];
                        const action = data[2];
                        const payload = data[3] || {};

                        console.log(`VU ${__VU}: Received incoming request: ${action} (ID: ${messageId})`);

                        // Generate and send response
                        const response = handleIncomingRequest(action, payload, messageId);
                        socket.send(response);

                        // Track incoming request handling
                        messageCount++;
                    }
                } catch (e) {
                    console.log(`VU ${__VU}: Received non-JSON message: ${message}`);
                }
            });

            socket.on('close', function () {
                const totalTime = Date.now() - connectionStart;
                console.log(`VU ${__VU}: Disconnected after ${totalTime}ms`);
                reconnectCount++;
            });

            socket.on('error', function (error) {
                console.log(`VU ${__VU}: Websocket error: ${error}`);
            });
        });

        check(res, {'status is 101': (r) => r && r.status === 101});

        // Small delay between reconnections to avoid overwhelming the server
        sleep(randomIntBetween(50, 100) / 1000);
    }

    // Final metrics
    console.log(`VU ${__VU}: Completed ${reconnectCount} connections with ${messageCount} messages in ${Date.now() - startTime}ms`);
}
