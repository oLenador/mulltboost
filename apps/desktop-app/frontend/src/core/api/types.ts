export interface Optimization {
    id: string;
    name: string;
    description: string;
    category: OptimizationCategory;
    level: OptimizationLevel;
    platform: Platform[];
    dependencies: string[];
    conflicts: string[];
    reversible: boolean;
    riskLevel: RiskLevel;
    version: string;
  }
  
  export interface OptimizationState {
    id: string;
    applied: boolean;
    appliedAt?: string;
    revertedAt?: string;
    version: string;
    backupData: Record<string, any>;
    status: ExecutionStatus;
    errorMsg?: string;
  }
  
  export interface OptimizationResult {
    success: boolean;
    message: string;
    backupData?: Record<string, any>;
    error?: string;
  }
  
  export interface BatchResult {
    totalCount: number;
    successCount: number;
    failedCount: number;
    results: Record<string, OptimizationResult>;
    overallStatus: string;
  }
  
  export interface SystemMetrics {
    cpu: CPUMetrics;
    memory: MemoryMetrics;
    gpu: GPUMetrics;
    network: NetworkMetrics;
    temperature: TemperatureMetrics;
    disk: DiskMetrics;
    timestamp: string;
  }
  
  export interface CPUMetrics {
    usage: number;
    coreCount: number;
    threadCount: number;
    frequency: number;
    temperature: number;
    cores: CoreMetric[];
  }
  
  export interface CoreMetric {
    index: number;
    usage: number;
    frequency: number;
    temperature: number;
  }
  
  export interface MemoryMetrics {
    total: number;
    used: number;
    available: number;
    usagePercent: number;
    cached: number;
    buffers: number;
  }
  
  export interface GPUMetrics {
    name: string;
    usage: number;
    memoryUsed: number;
    memoryTotal: number;
    temperature: number;
    powerDraw: number;
    clockSpeed: number;
    memoryClockSpeed: number;
  }
  
  export interface NetworkMetrics {
    interfaces: NetworkInterface[];
    totalSent: number;
    totalRecv: number;
  }
  
  export interface NetworkInterface {
    name: string;
    bytesSent: number;
    bytesRecv: number;
    speed: number;
    isUp: boolean;
  }
  
  export interface TemperatureMetrics {
    cpu: number;
    gpu: number;
    motherboard: number;
    drives: DriveTemp[];
  }
  
  export interface DriveTemp {
    name: string;
    temperature: number;
  }
  
  export interface DiskMetrics {
    drives: DriveMetric[];
  }
  
  export interface DriveMetric {
    name: string;
    total: number;
    used: number;
    free: number;
    usagePercent: number;
    readSpeed: number;
    writeSpeed: number;
  }
  
  export interface SystemInfo {
    os: OSInfo;
    cpu: CPUInfo;
    memory: MemoryInfo;
    gpu: GPUInfo[];
    storage: StorageInfo[];
    network: NetworkInfo[];
    motherboard: MotherboardInfo;
  }
  
  export interface OSInfo {
    name: string;
    version: string;
    architecture: string;
    buildNumber: string;
  }
  
  export interface CPUInfo {
    name: string;
    manufacturer: string;
    cores: number;
    threads: number;
    baseFreq: number;
    maxFreq: number;
    architecture: string;
  }
  
  export interface MemoryInfo {
    totalRAM: number;
    totalSlots: number;
    usedSlots: number;
    memoryType: string;
    speed: number;
  }
  
  export interface GPUInfo {
    name: string;
    manufacturer: string;
    memory: number;
    driver: string;
    directX: string;
    openGL: string;
  }
  
  export interface StorageInfo {
    name: string;
    type: string;
    size: number;
    model: string;
    interface: string;
  }
  
  export interface NetworkInfo {
    name: string;
    type: string;
    speed: number;
    macAddress: string;
    ipAddress: string;
  }
  
  export interface MotherboardInfo {
    manufacturer: string;
    model: string;
    bios: string;
    chipset: string;
  }
  
  // Enums
  export type OptimizationCategory = 'precision' | 'performance' | 'network' | 'system';
  export type OptimizationLevel = 'free' | 'premium';
  export type RiskLevel = 'low' | 'medium' | 'high';
  export type Platform = 'windows' | 'linux';
  export type ExecutionStatus = 'not_applied' | 'applied' | 'failed' | 'reverting' | 'reverted';
  
  // API Response Types
  export interface ApiResponse<T> {
    data: T;
    error?: string;
    success: boolean;
  }