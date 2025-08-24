export enum EventType {
    Processing = "booster.processing",
    Success = "booster.success",
    Error = "booster.error",
    Failed = "booster.failed",
    Queued = "booster.queued",
    BatchQueued = "booster.batch_queued",
    Cancelled = "booster.cancelled",
  }
  
  export interface EventData {
    readonly eventType: EventType;
    readonly timestamp: Date;
    readonly operationType: string;
    readonly operationId: string;
    readonly boosterId: string;
    readonly sequence: number;
    readonly idempotencyId: string;
    readonly status: string;
    readonly endAt?: Date;
    readonly error?: string;
    readonly queueSize: number;
  }
  
  export function createEventData(rawData: any): EventData {
    return {
      eventType: rawData.EventType as EventType,
      timestamp: new Date(rawData.Timestamp),
      operationType: rawData.OperationType,
      operationId: rawData.OperationID,
      boosterId: rawData.BoosterID,
      sequence: rawData.Sequency || 0,
      idempotencyId: rawData.IdempotencyID,
      status: rawData.Status,
      endAt: rawData.EndAt ? new Date(rawData.EndAt) : undefined,
      error: rawData.Error,
      queueSize: rawData.QueueSize || 0,
    };
  }