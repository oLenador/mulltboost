export namespace dto {
	
	export class BoosterDto {
	    ID: string;
	    Name: string;
	    Description: string;
	    Category: string;
	    Level: string;
	    Platform: string[];
	    Dependencies: string[];
	    Conflicts: string[];
	    Reversible: boolean;
	    RiskLevel: string;
	    Version: string;
	    Tags: string[];
	
	    static createFrom(source: any = {}) {
	        return new BoosterDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.Description = source["Description"];
	        this.Category = source["Category"];
	        this.Level = source["Level"];
	        this.Platform = source["Platform"];
	        this.Dependencies = source["Dependencies"];
	        this.Conflicts = source["Conflicts"];
	        this.Reversible = source["Reversible"];
	        this.RiskLevel = source["RiskLevel"];
	        this.Version = source["Version"];
	        this.Tags = source["Tags"];
	    }
	}

}

export namespace entities {
	
	export class BoosterResult {
	    Success: boolean;
	    Message: string;
	    BackupData: Record<string, any>;
	    Error: any;
	
	    static createFrom(source: any = {}) {
	        return new BoosterResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Success = source["Success"];
	        this.Message = source["Message"];
	        this.BackupData = source["BackupData"];
	        this.Error = source["Error"];
	    }
	}
	export class BatchResult {
	    TotalCount: number;
	    SuccessCount: number;
	    FailedCount: number;
	    Results: Record<string, BoosterResult>;
	    OverallStatus: string;
	
	    static createFrom(source: any = {}) {
	        return new BatchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TotalCount = source["TotalCount"];
	        this.SuccessCount = source["SuccessCount"];
	        this.FailedCount = source["FailedCount"];
	        this.Results = this.convertValues(source["Results"], BoosterResult, true);
	        this.OverallStatus = source["OverallStatus"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class BoosterRollbackState {
	    ID: string;
	    Applied: boolean;
	    // Go type: time
	    AppliedAt?: any;
	    // Go type: time
	    RevertedAt?: any;
	    Version: string;
	    BackupData: Record<string, any>;
	    Status: string;
	    ErrorMsg: string;
	
	    static createFrom(source: any = {}) {
	        return new BoosterRollbackState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Applied = source["Applied"];
	        this.AppliedAt = this.convertValues(source["AppliedAt"], null);
	        this.RevertedAt = this.convertValues(source["RevertedAt"], null);
	        this.Version = source["Version"];
	        this.BackupData = source["BackupData"];
	        this.Status = source["Status"];
	        this.ErrorMsg = source["ErrorMsg"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CPUInfo {
	    Name: string;
	    Manufacturer: string;
	    Cores: number;
	    Threads: number;
	    BaseFreq: number;
	    MaxFreq: number;
	    Architecture: string;
	
	    static createFrom(source: any = {}) {
	        return new CPUInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Manufacturer = source["Manufacturer"];
	        this.Cores = source["Cores"];
	        this.Threads = source["Threads"];
	        this.BaseFreq = source["BaseFreq"];
	        this.MaxFreq = source["MaxFreq"];
	        this.Architecture = source["Architecture"];
	    }
	}
	export class CoreMetric {
	    Index: number;
	    Usage: number;
	    Frequency: number;
	    Temperature: number;
	
	    static createFrom(source: any = {}) {
	        return new CoreMetric(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Index = source["Index"];
	        this.Usage = source["Usage"];
	        this.Frequency = source["Frequency"];
	        this.Temperature = source["Temperature"];
	    }
	}
	export class CPUMetrics {
	    Usage: number;
	    CoreCount: number;
	    ThreadCount: number;
	    Frequency: number;
	    Temperature: number;
	    Cores: CoreMetric[];
	
	    static createFrom(source: any = {}) {
	        return new CPUMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Usage = source["Usage"];
	        this.CoreCount = source["CoreCount"];
	        this.ThreadCount = source["ThreadCount"];
	        this.Frequency = source["Frequency"];
	        this.Temperature = source["Temperature"];
	        this.Cores = this.convertValues(source["Cores"], CoreMetric);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class DriveMetric {
	    Name: string;
	    Total: number;
	    Used: number;
	    Free: number;
	    UsagePercent: number;
	    ReadSpeed: number;
	    WriteSpeed: number;
	
	    static createFrom(source: any = {}) {
	        return new DriveMetric(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Total = source["Total"];
	        this.Used = source["Used"];
	        this.Free = source["Free"];
	        this.UsagePercent = source["UsagePercent"];
	        this.ReadSpeed = source["ReadSpeed"];
	        this.WriteSpeed = source["WriteSpeed"];
	    }
	}
	export class DiskMetrics {
	    Drives: DriveMetric[];
	
	    static createFrom(source: any = {}) {
	        return new DiskMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Drives = this.convertValues(source["Drives"], DriveMetric);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class DriveTemp {
	    Name: string;
	    Temperature: number;
	
	    static createFrom(source: any = {}) {
	        return new DriveTemp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Temperature = source["Temperature"];
	    }
	}
	export class GPUInfo {
	    Name: string;
	    Manufacturer: string;
	    Memory: number;
	    Driver: string;
	    DirectX: string;
	    OpenGL: string;
	
	    static createFrom(source: any = {}) {
	        return new GPUInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Manufacturer = source["Manufacturer"];
	        this.Memory = source["Memory"];
	        this.Driver = source["Driver"];
	        this.DirectX = source["DirectX"];
	        this.OpenGL = source["OpenGL"];
	    }
	}
	export class GPUMetrics {
	    Name: string;
	    Usage: number;
	    MemoryUsed: number;
	    MemoryTotal: number;
	    Temperature: number;
	    PowerDraw: number;
	    ClockSpeed: number;
	    MemoryClockSpeed: number;
	
	    static createFrom(source: any = {}) {
	        return new GPUMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Usage = source["Usage"];
	        this.MemoryUsed = source["MemoryUsed"];
	        this.MemoryTotal = source["MemoryTotal"];
	        this.Temperature = source["Temperature"];
	        this.PowerDraw = source["PowerDraw"];
	        this.ClockSpeed = source["ClockSpeed"];
	        this.MemoryClockSpeed = source["MemoryClockSpeed"];
	    }
	}
	export class MemoryInfo {
	    TotalRAM: number;
	    TotalSlots: number;
	    UsedSlots: number;
	    MemoryType: string;
	    Speed: number;
	
	    static createFrom(source: any = {}) {
	        return new MemoryInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TotalRAM = source["TotalRAM"];
	        this.TotalSlots = source["TotalSlots"];
	        this.UsedSlots = source["UsedSlots"];
	        this.MemoryType = source["MemoryType"];
	        this.Speed = source["Speed"];
	    }
	}
	export class MemoryMetrics {
	    Total: number;
	    Used: number;
	    Available: number;
	    UsagePercent: number;
	    Cached: number;
	    Buffers: number;
	
	    static createFrom(source: any = {}) {
	        return new MemoryMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Total = source["Total"];
	        this.Used = source["Used"];
	        this.Available = source["Available"];
	        this.UsagePercent = source["UsagePercent"];
	        this.Cached = source["Cached"];
	        this.Buffers = source["Buffers"];
	    }
	}
	export class MotherboardInfo {
	    Manufacturer: string;
	    Model: string;
	    BIOS: string;
	    Chipset: string;
	
	    static createFrom(source: any = {}) {
	        return new MotherboardInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Manufacturer = source["Manufacturer"];
	        this.Model = source["Model"];
	        this.BIOS = source["BIOS"];
	        this.Chipset = source["Chipset"];
	    }
	}
	export class NetworkInfo {
	    Name: string;
	    Type: string;
	    Speed: number;
	    MACAddress: string;
	    IPAddress: string;
	
	    static createFrom(source: any = {}) {
	        return new NetworkInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Speed = source["Speed"];
	        this.MACAddress = source["MACAddress"];
	        this.IPAddress = source["IPAddress"];
	    }
	}
	export class NetworkInterface {
	    Name: string;
	    BytesSent: number;
	    BytesRecv: number;
	    Speed: number;
	    IsUp: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NetworkInterface(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.BytesSent = source["BytesSent"];
	        this.BytesRecv = source["BytesRecv"];
	        this.Speed = source["Speed"];
	        this.IsUp = source["IsUp"];
	    }
	}
	export class NetworkMetrics {
	    Interfaces: NetworkInterface[];
	    TotalSent: number;
	    TotalRecv: number;
	
	    static createFrom(source: any = {}) {
	        return new NetworkMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Interfaces = this.convertValues(source["Interfaces"], NetworkInterface);
	        this.TotalSent = source["TotalSent"];
	        this.TotalRecv = source["TotalRecv"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class OSInfo {
	    Name: string;
	    Version: string;
	    Architecture: string;
	    BuildNumber: string;
	
	    static createFrom(source: any = {}) {
	        return new OSInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Version = source["Version"];
	        this.Architecture = source["Architecture"];
	        this.BuildNumber = source["BuildNumber"];
	    }
	}
	export class StorageInfo {
	    Name: string;
	    Type: string;
	    Size: number;
	    Model: string;
	    Interface: string;
	
	    static createFrom(source: any = {}) {
	        return new StorageInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Size = source["Size"];
	        this.Model = source["Model"];
	        this.Interface = source["Interface"];
	    }
	}
	export class SystemInfo {
	    OS: OSInfo;
	    CPU: CPUInfo;
	    Memory: MemoryInfo;
	    GPU: GPUInfo[];
	    Storage: StorageInfo[];
	    Network: NetworkInfo[];
	    Motherboard: MotherboardInfo;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.OS = this.convertValues(source["OS"], OSInfo);
	        this.CPU = this.convertValues(source["CPU"], CPUInfo);
	        this.Memory = this.convertValues(source["Memory"], MemoryInfo);
	        this.GPU = this.convertValues(source["GPU"], GPUInfo);
	        this.Storage = this.convertValues(source["Storage"], StorageInfo);
	        this.Network = this.convertValues(source["Network"], NetworkInfo);
	        this.Motherboard = this.convertValues(source["Motherboard"], MotherboardInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TemperatureMetrics {
	    CPU: number;
	    GPU: number;
	    Motherboard: number;
	    Drives: DriveTemp[];
	
	    static createFrom(source: any = {}) {
	        return new TemperatureMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.CPU = source["CPU"];
	        this.GPU = source["GPU"];
	        this.Motherboard = source["Motherboard"];
	        this.Drives = this.convertValues(source["Drives"], DriveTemp);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SystemMetrics {
	    CPU: CPUMetrics;
	    Memory: MemoryMetrics;
	    GPU: GPUMetrics;
	    Network: NetworkMetrics;
	    Temperature: TemperatureMetrics;
	    Disk: DiskMetrics;
	    // Go type: time
	    Timestamp: any;
	
	    static createFrom(source: any = {}) {
	        return new SystemMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.CPU = this.convertValues(source["CPU"], CPUMetrics);
	        this.Memory = this.convertValues(source["Memory"], MemoryMetrics);
	        this.GPU = this.convertValues(source["GPU"], GPUMetrics);
	        this.Network = this.convertValues(source["Network"], NetworkMetrics);
	        this.Temperature = this.convertValues(source["Temperature"], TemperatureMetrics);
	        this.Disk = this.convertValues(source["Disk"], DiskMetrics);
	        this.Timestamp = this.convertValues(source["Timestamp"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

