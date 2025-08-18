import {
    PipelineManager,
    PipelineItemStatus,
    BasePipelineItem,
    SimpleTaskItem,
    ChainPipelineItem,
    ChainHandler,
    PipelineContext,
    ExecutionResult
  } from './pipeline-manager';
  
  // Test utilities
  class MockPipelineItem extends BasePipelineItem {
    constructor(
      id: string,
      name: string,
      private executeFn: (context: PipelineContext) => Promise<any>,
      dependencies: string[] = [],
      private canExecuteFn?: (context: PipelineContext) => boolean
    ) {
      super(id, name, 'mock', dependencies);
    }
  
    async execute(context: PipelineContext): Promise<any> {
      return this.executeFn(context);
    }
  
    canExecute(context: PipelineContext): boolean {
      return this.canExecuteFn ? this.canExecuteFn(context) : true;
    }
  }
  
  // Helper function to create delay
  const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));
  
  describe('PipelineManager', () => {
    let manager: PipelineManager;
  
    beforeEach(() => {
      manager = new PipelineManager();
    });
  
    describe('Basic Pipeline Operations', () => {
      test('should add and retrieve items', () => {
        const item = new MockPipelineItem('test1', 'Test Item', async () => 'result');
        manager.addItem(item);
  
        const retrieved = manager.getItem('test1');
        expect(retrieved).toBe(item);
        expect(manager.getAllItems()).toHaveLength(1);
      });
  
      test('should remove items', () => {
        const item = new MockPipelineItem('test1', 'Test Item', async () => 'result');
        manager.addItem(item);
        manager.removeItem('test1');
  
        expect(manager.getItem('test1')).toBeUndefined();
        expect(manager.getAllItems()).toHaveLength(0);
      });
  
      test('should execute simple pipeline without dependencies', async () => {
        const item1 = new MockPipelineItem('item1', 'Item 1', async (ctx) => {
          ctx.item1Result = 'done';
          return 'item1-result';
        });
  
        const item2 = new MockPipelineItem('item2', 'Item 2', async (ctx) => {
          ctx.item2Result = 'done';
          return 'item2-result';
        });
  
        manager.addItem(item1);
        manager.addItem(item2);
  
        const context = { initial: 'data' };
        const results = await manager.execute(context);
  
        expect(results.size).toBe(2);
        expect(results.get('item1')?.status).toBe(PipelineItemStatus.SUCCESS);
        expect(results.get('item2')?.status).toBe(PipelineItemStatus.SUCCESS);
        expect(context.item1Result).toBe('done');
        expect(context.item2Result).toBe('done');
      });
    });
  
    describe('Dependency Management', () => {
      test('should execute items in correct dependency order', async () => {
        const executionOrder: string[] = [];
  
        const item1 = new MockPipelineItem('item1', 'Item 1', async () => {
          executionOrder.push('item1');
          return 'result1';
        });
  
        const item2 = new MockPipelineItem('item2', 'Item 2', async () => {
          executionOrder.push('item2');
          return 'result2';
        }, ['item1']); // depends on item1
  
        const item3 = new MockPipelineItem('item3', 'Item 3', async () => {
          executionOrder.push('item3');
          return 'result3';
        }, ['item1', 'item2']); // depends on both
  
        manager.addItem(item2);
        manager.addItem(item3);
        manager.addItem(item1); // Add in random order
  
        await manager.execute();
  
        expect(executionOrder).toEqual(['item1', 'item2', 'item3']);
      });
  
      test('should detect circular dependencies', () => {
        const item1 = new MockPipelineItem('item1', 'Item 1', async () => 'result', ['item2']);
        const item2 = new MockPipelineItem('item2', 'Item 2', async () => 'result', ['item1']);
  
        manager.addItem(item1);
        expect(() => manager.addItem(item2)).toThrow('Circular dependency detected');
      });
  
      test('should skip items when dependencies fail', async () => {
        const item1 = new MockPipelineItem('item1', 'Item 1', async () => {
          throw new Error('Item 1 failed');
        });
  
        const item2 = new MockPipelineItem('item2', 'Item 2', async () => 'result', ['item1']);
  
        manager.addItem(item1);
        manager.addItem(item2);
  
        const results = await manager.execute();
  
        expect(results.get('item1')?.status).toBe(PipelineItemStatus.FAILED);
        expect(results.get('item2')?.status).toBe(PipelineItemStatus.SKIPPED);
      });
    });
  
    describe('Error Handling and Retries', () => {
      test('should retry failed items', async () => {
        let attempts = 0;
        const item = new MockPipelineItem(
          'retry-item',
          'Retry Item',
          async () => {
            attempts++;
            if (attempts < 3) {
              throw new Error(`Attempt ${attempts} failed`);
            }
            return 'success';
          },
          [],
        );
        item.config.retries = 2;
  
        manager.addItem(item);
        const results = await manager.execute();
  
        expect(attempts).toBe(3);
        expect(results.get('retry-item')?.status).toBe(PipelineItemStatus.SUCCESS);
      });
  
      test('should fail after max retries', async () => {
        let attempts = 0;
        const item = new MockPipelineItem(
          'fail-item',
          'Fail Item',
          async () => {
            attempts++;
            throw new Error(`Attempt ${attempts} failed`);
          },
          [],
        );
        item.config.retries = 2;
  
        manager.addItem(item);
        const results = await manager.execute();
  
        expect(attempts).toBe(3); // initial + 2 retries
        expect(results.get('fail-item')?.status).toBe(PipelineItemStatus.FAILED);
      });
  
      test('should handle timeout', async () => {
        const item = new MockPipelineItem(
          'timeout-item',
          'Timeout Item',
          async () => {
            await delay(200); // Will timeout
            return 'result';
          },
          [],
        );
        item.config.timeout = 100; // 100ms timeout
  
        manager.addItem(item);
        const results = await manager.execute();
  
        expect(results.get('timeout-item')?.status).toBe(PipelineItemStatus.FAILED);
        expect(results.get('timeout-item')?.error?.message).toBe('Execution timeout');
      });
  
      test('should continue pipeline when skipOnFailure is true', async () => {
        const item1 = new MockPipelineItem('item1', 'Item 1', async () => {
          throw new Error('Failed');
        });
        item1.config.skipOnFailure = true;
  
        const item2 = new MockPipelineItem('item2', 'Item 2', async () => 'success');
  
        manager.addItem(item1);
        manager.addItem(item2);
  
        const results = await manager.execute();
  
        expect(results.get('item1')?.status).toBe(PipelineItemStatus.FAILED);
        expect(results.get('item2')?.status).toBe(PipelineItemStatus.SUCCESS);
      });
    });
  
    describe('Conditional Execution', () => {
      test('should skip items based on canExecute condition', async () => {
        const item = new MockPipelineItem(
          'conditional-item',
          'Conditional Item',
          async () => 'result',
          [],
          (ctx) => ctx.shouldExecute === true
        );
  
        manager.addItem(item);
  
        // First execution without condition
        let results = await manager.execute({ shouldExecute: false });
        expect(results.get('conditional-item')?.status).toBe(PipelineItemStatus.SKIPPED);
  
        // Second execution with condition
        results = await manager.execute({ shouldExecute: true });
        expect(results.get('conditional-item')?.status).toBe(PipelineItemStatus.SUCCESS);
      });
    });
  });
  
  describe('SimpleTaskItem', () => {
    test('should execute simple task', async () => {
      const task = vi.fn().mockResolvedValue('task-result');
      const item = new SimpleTaskItem('simple', 'Simple Task', task);
  
      const context = { input: 'test' };
      const result = await item.execute(context);
  
      expect(task).toHaveBeenCalledWith(context);
      expect(result).toBe('task-result');
    });
  });
  
  describe('Chain of Responsibility', () => {
    test('should execute chain handlers in order', async () => {
      const executionOrder: string[] = [];
  
      class TestHandler1 extends ChainHandler {
        protected async process(context: PipelineContext): Promise<PipelineContext> {
          executionOrder.push('handler1');
          return { ...context, step1: 'done' };
        }
      }
  
      class TestHandler2 extends ChainHandler {
        protected async process(context: PipelineContext): Promise<PipelineContext> {
          executionOrder.push('handler2');
          return { ...context, step2: 'done' };
        }
      }
  
      const handler1 = new TestHandler1();
      const handler2 = new TestHandler2();
      handler1.setNext(handler2);
  
      const result = await handler1.handle({ initial: 'data' });
  
      expect(executionOrder).toEqual(['handler1', 'handler2']);
      expect(result).toEqual({
        initial: 'data',
        step1: 'done',
        step2: 'done'
      });
    });
  
    test('should stop chain when shouldContinue returns false', async () => {
      class StopHandler extends ChainHandler {
        protected async process(context: PipelineContext): Promise<PipelineContext> {
          return { ...context, stopped: true };
        }
  
        protected shouldContinue(): boolean {
          return false;
        }
      }
  
      class NextHandler extends ChainHandler {
        protected async process(context: PipelineContext): Promise<PipelineContext> {
          return { ...context, continued: true };
        }
      }
  
      const stopHandler = new StopHandler();
      const nextHandler = new NextHandler();
      stopHandler.setNext(nextHandler);
  
      const result = await stopHandler.handle({ initial: 'data' });
  
      expect(result.stopped).toBe(true);
      expect(result.continued).toBeUndefined();
    });
  });
  
  
  