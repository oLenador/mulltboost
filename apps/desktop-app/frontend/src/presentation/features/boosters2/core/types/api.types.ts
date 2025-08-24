export interface QueueStateData {
    readonly items: QueueItem[];
    readonly inProgress: number;
  }
  
  export interface QueueItem {
    readonly boosterId: string;
    readonly operationId?: string;
    readonly progress?: number;
    readonly error?: string;
  }
  
  export function createQueueStateData(rawData: any): QueueStateData {
    if (Array.isArray(rawData)) {
      return {
        items: rawData.map(createQueueItem),
        inProgress: 0,
      };
    }
    
    return {
      items: (rawData.Items || rawData.items || []).map(createQueueItem),
      inProgress: rawData.InProgress || rawData.inProgress || 0,
    };
  }
  
  function createQueueItem(rawItem: any): QueueItem {
    return {
      boosterId: rawItem.BoosterID || rawItem.boosterId,
      operationId: rawItem.OperationID || rawItem.operationId,
      progress: rawItem.Progress || rawItem.progress,
      error: rawItem.Error || rawItem.error,
    };
  }
  
  export interface ApiResponse<T = any> {
    readonly success: boolean;
    readonly data?: T;
    readonly error?: string;
    readonly message?: string;
    readonly operationId?: string;
  }
  
  export function createApiResponse<T>(
    success: boolean,
    data?: T,
    error?: string,
    message?: string,
    operationId?: string
  ): ApiResponse<T> {
    return {
      success,
      data,
      error,
      message,
      operationId,
    };
  }