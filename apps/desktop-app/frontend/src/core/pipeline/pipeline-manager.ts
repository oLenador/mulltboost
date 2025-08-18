// Types e Interfaces
export enum PipelineItemStatus {
    PENDING = 'pending',
    RUNNING = 'running',
    SUCCESS = 'success',
    FAILED = 'failed',
    SKIPPED = 'skipped'
  }
  
  export interface PipelineContext {
    [key: string]: any;
  }
  
  export interface PipelineItemConfig {
    timeout?: number;
    retries?: number;
    resources?: Record<string, any>;
    skipOnFailure?: boolean;
  }
  
  export interface ExecutionResult {
    status: PipelineItemStatus;
    output?: any;
    error?: Error;
    duration: number;
    logs?: string[];
  }
  
  export interface PipelineItem {
    id: string;
    name: string;
    type: string;
    status: PipelineItemStatus;
    dependencies: string[];
    config: PipelineItemConfig;
    execute(context: PipelineContext): Promise<any>;
    canExecute?(context: PipelineContext): boolean;
  }
  
  // Base Pipeline Item Class
  export abstract class BasePipelineItem implements PipelineItem {
    public status: PipelineItemStatus = PipelineItemStatus.PENDING;
    
    constructor(
      public readonly id: string,
      public readonly name: string,
      public readonly type: string,
      public readonly dependencies: string[] = [],
      public readonly config: PipelineItemConfig = {}
    ) {}
  
    abstract execute(context: PipelineContext): Promise<any>;
  
    canExecute?(context: PipelineContext): boolean {
      return true;
    }
  }
  
  // Pipeline Execution Engine
  export class PipelineManager {
    private items: Map<string, PipelineItem> = new Map();
    private executionOrder: string[] = [];
    private results: Map<string, ExecutionResult> = new Map();
  
    addItem(item: PipelineItem): void {
      this.items.set(item.id, item);
      this.calculateExecutionOrder();
    }
  
    removeItem(id: string): void {
      this.items.delete(id);
      this.calculateExecutionOrder();
    }
  
    private calculateExecutionOrder(): void {
      const visited = new Set<string>();
      const visiting = new Set<string>();
      const order: string[] = [];
  
      const visit = (itemId: string): void => {
        if (visiting.has(itemId)) {
          throw new Error(`Circular dependency detected involving ${itemId}`);
        }
        if (visited.has(itemId)) return;
  
        const item = this.items.get(itemId);
        if (!item) throw new Error(`Item ${itemId} not found`);
  
        visiting.add(itemId);
  
        for (const dep of item.dependencies) {
          if (!this.items.has(dep)) continue;
          visit(dep);
        }
  
        visiting.delete(itemId);
        visited.add(itemId);
        order.push(itemId);
      };
  
      for (const itemId of this.items.keys()) {
        if (!visited.has(itemId)) {
          visit(itemId);
        }
      }
  
      this.executionOrder = order;
    }
  
    async execute(context: PipelineContext = {}): Promise<Map<string, ExecutionResult>> {
      this.results.clear();
      this.calculateExecutionOrder();

      for (const itemId of this.executionOrder) {
        const item = this.items.get(itemId)!;
        
        // Check dependencies
        if (!this.areDependenciesSatisfied(item)) {
          const result: ExecutionResult = {
            status: PipelineItemStatus.SKIPPED,
            duration: 0,
            logs: ['Skipped due to failed dependencies']
          };
          this.results.set(itemId, result);
          item.status = PipelineItemStatus.SKIPPED;
          continue;
        }
  
        // Check if item can execute
        if (item.canExecute && !item.canExecute(context)) {
          const result: ExecutionResult = {
            status: PipelineItemStatus.SKIPPED,
            duration: 0,
            logs: ['Skipped due to execution conditions']
          };
          this.results.set(itemId, result);
          item.status = PipelineItemStatus.SKIPPED;
          continue;
        }
  
        // Execute item
        const result = await this.executeItem(item, context);
        this.results.set(itemId, result);
  
        // Stop pipeline if item failed and no skip on failure
        if (result.status === PipelineItemStatus.FAILED && !item.config.skipOnFailure) {
          break;
        }
      }
  
      return this.results;
    }
  
    private areDependenciesSatisfied(item: PipelineItem): boolean {
      return item.dependencies.every(depId => {
        if (!this.items.has(depId)) {
          return true; 
        }
        const result = this.results.get(depId);
        return result && result.status === PipelineItemStatus.SUCCESS;
      });
    }
  
    private async executeItem(item: PipelineItem, context: PipelineContext): Promise<ExecutionResult> {
      const startTime = Date.now();
      const logs: string[] = [];
      let retries = item.config.retries || 0;
  
      item.status = PipelineItemStatus.RUNNING;
      logs.push(`Started execution of ${item.name}`);
  
      while (retries >= 0) {
        try {
          const timeoutPromise = item.config.timeout 
            ? new Promise((_, reject) => 
                setTimeout(() => reject(new Error('Execution timeout')), item.config.timeout)
              )
            : new Promise(() => {}); // Never resolves
  
          const executionPromise = item.execute(context);
          const output = await Promise.race([executionPromise, timeoutPromise]);
  
          const duration = Date.now() - startTime;
          item.status = PipelineItemStatus.SUCCESS;
          
          logs.push(`Completed execution of ${item.name} in ${duration}ms`);
          
          return {
            status: PipelineItemStatus.SUCCESS,
            output,
            duration,
            logs
          };
        } catch (error) {
          retries--;
          const err = error as Error;
          
          if (retries >= 0) {
            logs.push(`Retry ${(item.config.retries || 0) - retries} for ${item.name}: ${err.message}`);
            continue;
          }
  
          const duration = Date.now() - startTime;
          item.status = PipelineItemStatus.FAILED;
          logs.push(`Failed execution of ${item.name}: ${err.message}`);
  
          return {
            status: PipelineItemStatus.FAILED,
            error: err,
            duration,
            logs
          };
        }
      }
  
      // Should never reach here, but TypeScript needs it
      throw new Error('Unexpected execution path');
    }
  
    getExecutionOrder(): string[] {
      return [...this.executionOrder];
    }
  
    getResults(): Map<string, ExecutionResult> {
      return new Map(this.results);
    }
  
    getItem(id: string): PipelineItem | undefined {
      return this.items.get(id);
    }
  
    getAllItems(): PipelineItem[] {
      return Array.from(this.items.values());
    }
  }
  
  // Chain of Responsibility Implementation
  export abstract class ChainHandler {
    private nextHandler?: ChainHandler;
  
    setNext(handler: ChainHandler): ChainHandler {
      this.nextHandler = handler;
      return handler;
    }
  
    async handle(context: PipelineContext): Promise<PipelineContext> {
      const result = await this.process(context);
      
      if (this.nextHandler && this.shouldContinue(result)) {
        return this.nextHandler.handle(result);
      }
      
      return result;
    }
  
    protected abstract process(context: PipelineContext): Promise<PipelineContext>;
    
    protected shouldContinue(context: PipelineContext): boolean {
      return true;
    }
  }
  
  // Chain-based Pipeline Item
  export class ChainPipelineItem extends BasePipelineItem {
    private chainHead?: ChainHandler;
  
    constructor(
      id: string,
      name: string,
      type: string,
      dependencies: string[] = [],
      config: PipelineItemConfig = {}
    ) {
      super(id, name, type, dependencies, config);
    }
  
    setChain(handler: ChainHandler): void {
      this.chainHead = handler;
    }
  
    async execute(context: PipelineContext): Promise<any> {
      if (!this.chainHead) {
        throw new Error(`No chain configured for item ${this.id}`);
      }
  
      return this.chainHead.handle({ ...context });
    }
  }
  
  // Example Implementations
  export class SimpleTaskItem extends BasePipelineItem {
    constructor(
      id: string,
      name: string,
      private task: (context: PipelineContext) => Promise<any>,
      dependencies: string[] = [],
      config: PipelineItemConfig = {}
    ) {
      super(id, name, 'simple-task', dependencies, config);
    }
  
    async execute(context: PipelineContext): Promise<any> {
      return this.task(context);
    }
  }
