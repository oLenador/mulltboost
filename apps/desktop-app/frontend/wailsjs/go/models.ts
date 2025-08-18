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
	export class GPUInfo {
	    id: string;
	    name: string;
	    vendor: string;
	    device_id: string;
	    driver_version: string;
	    vram_size: number;
	    vram_used: number;
	    core_clock: number;
	    memory_clock: number;
	    temperature: number;
	    usage: number;
	    power_usage: number;
	    is_primary: boolean;
	    is_discrete: boolean;
	    supports_directx: string;
	    supports_opengl: string;
	    supports_vulkan: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new GPUInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.vendor = source["vendor"];
	        this.device_id = source["device_id"];
	        this.driver_version = source["driver_version"];
	        this.vram_size = source["vram_size"];
	        this.vram_used = source["vram_used"];
	        this.core_clock = source["core_clock"];
	        this.memory_clock = source["memory_clock"];
	        this.temperature = source["temperature"];
	        this.usage = source["usage"];
	        this.power_usage = source["power_usage"];
	        this.is_primary = source["is_primary"];
	        this.is_discrete = source["is_discrete"];
	        this.supports_directx = source["supports_directx"];
	        this.supports_opengl = source["supports_opengl"];
	        this.supports_vulkan = source["supports_vulkan"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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

}

