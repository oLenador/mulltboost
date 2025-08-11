declare global {
    interface Window {
      go: {
        handlers: {
          OptimizationHandler: OptimizationHandler;
          MonitoringHandler: MonitoringHandler;
          SystemHandler: SystemHandler;
        };
      };
    }
  }
  
  interface OptimizationHandler {
    GetAvailableOptimizations(): Promise<Optimization[]>;
    GetOptimizationState(id: string): Promise<OptimizationState>;
    ApplyOptimization(id: string): Promise<OptimizationResult>;
    RevertOptimization(id: string): Promise<OptimizationResult>;
    ApplyOptimizationBatch(ids: string[]): Promise<BatchResult>;
    GetOptimizationsByCategory(category: string): Promise<Optimization[]>;
  }
  
  interface MonitoringHandler {
    GetSystemMetrics(): Promise<SystemMetrics>;
    StartRealTimeMonitoring(intervalSeconds: number): Promise<void>;
    StopRealTimeMonitoring(): Promise<void>;
    IsMonitoring(): Promise<boolean>;
  }
  
  interface SystemHandler {
    GetSystemInfo(): Promise<SystemInfo>;
    GetHardwareInfo(): Promise<SystemInfo>;
    RefreshSystemInfo(): Promise<void>;
  }
  
  export class ApiClient {
    private static instance: ApiClient;
    
    private constructor() {}
    
    static getInstance(): ApiClient {
      if (!ApiClient.instance) {
        ApiClient.instance = new ApiClient();
      }
      return ApiClient.instance;
    }
  
    // Optimization methods
    async getAvailableOptimizations(): Promise<Optimization[]> {
      try {
        return await window.go.handlers.OptimizationHandler.GetAvailableOptimizations();
      } catch (error) {
        console.error('Failed to get available optimizations:', error);
        throw new Error('Failed to get available optimizations');
      }
    }
  
    async getOptimizationState(id: string): Promise<OptimizationState> {
      try {
        return await window.go.handlers.OptimizationHandler.GetOptimizationState(id);
      } catch (error) {
        console.error(`Failed to get optimization state for ${id}:`, error);
        throw new Error(`Failed to get optimization state for ${id}`);
      }
    }
  
    async applyOptimization(id: string): Promise<OptimizationResult> {
      try {
        return await window.go.handlers.OptimizationHandler.ApplyOptimization(id);
      } catch (error) {
        console.error(`Failed to apply optimization ${id}:`, error);
        throw new Error(`Failed to apply optimization ${id}`);
      }
    }
  
    async revertOptimization(id: string): Promise<OptimizationResult> {
      try {
        return await window.go.handlers.OptimizationHandler.RevertOptimization(id);
      } catch (error) {
        console.error(`Failed to revert optimization ${id}:`, error);
        throw new Error(`Failed to revert optimization ${id}`);
      }
    }
  
    async applyOptimizationBatch(ids: string[]): Promise<BatchResult> {
      try {
        return await window.go.handlers.OptimizationHandler.ApplyOptimizationBatch(ids);
      } catch (error) {
        console.error('Failed to apply optimization batch:', error);
        throw new Error('Failed to apply optimization batch');
      }
    }
  
    async getOptimizationsByCategory(category: OptimizationCategory): Promise<Optimization[]> {
      try {
        return await window.go.handlers.OptimizationHandler.GetOptimizationsByCategory(category);
      } catch (error) {
        console.error(`Failed to get optimizations for category ${category}:`, error);
        throw new Error(`Failed to get optimizations for category ${category}`);
      }
    }
  
    // Monitoring methods
    async getSystemMetrics(): Promise<SystemMetrics> {
      try {
        return await window.go.handlers.MonitoringHandler.GetSystemMetrics();
      } catch (error) {
        console.error('Failed to get system metrics:', error);
        throw new Error('Failed to get system metrics');
      }
    }
  
    async startRealTimeMonitoring(intervalSeconds: number = 1): Promise<void> {
      try {
        await window.go.handlers.MonitoringHandler.StartRealTimeMonitoring(intervalSeconds);
      } catch (error) {
        console.error('Failed to start real-time monitoring:', error);
        throw new Error('Failed to start real-time monitoring');
      }
    }
  
    async stopRealTimeMonitoring(): Promise<void> {
      try {
        await window.go.handlers.MonitoringHandler.StopRealTimeMonitoring();
      } catch (error) {
        console.error('Failed to stop real-time monitoring:', error);
        throw new Error('Failed to stop real-time monitoring');
      }
    }
  
    async isMonitoring(): Promise<boolean> {
      try {
        return await window.go.handlers.MonitoringHandler.IsMonitoring();
      } catch (error) {
        console.error('Failed to check monitoring status:', error);
        return false;
      }
    }
  
    // System info methods
    async getSystemInfo(): Promise<SystemInfo> {
      try {
        return await window.go.handlers.SystemHandler.GetSystemInfo();
      } catch (error) {
        console.error('Failed to get system info:', error);
        throw new Error('Failed to get system info');
      }
    }
  
    async getHardwareInfo(): Promise<SystemInfo> {
      try {
        return await window.go.handlers.SystemHandler.GetHardwareInfo();
      } catch (error) {
        console.error('Failed to get hardware info:', error);
        throw new Error('Failed to get hardware info');
      }
    }
  
    async refreshSystemInfo(): Promise<void> {
      try {
        await window.go.handlers.SystemHandler.RefreshSystemInfo();
      } catch (error) {
        console.error('Failed to refresh system info:', error);
        throw new Error('Failed to refresh system info');
      }
    }
  }
  
  export const apiClient = ApiClient.getInstance();