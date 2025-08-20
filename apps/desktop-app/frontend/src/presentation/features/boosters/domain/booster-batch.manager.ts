import { BaseBatchManager } from "@/core/domain/batch/batch.manager";
import { BoosterItem } from "../types/booster.types";
import { GetExecutionQueueState, InitBoosterApply } from "bindings/github.com/oLenador/mulltbost/internal/app/handlers/boosterhandler";
import { QueueItem } from "bindings/github.com/oLenador/mulltbost/internal/core/domain/entities";
import { BoosterBatchProgressEventData, BoosterEvent, BoosterEventData } from "@/core/domain/batch/events.types";
import { Events } from '@wailsio/runtime';
import { ItemStatus } from "@/core/domain/batch/batch.types";

export class BoosterBatchManager extends BaseBatchManager<BoosterItem> {
  private eventUnsubscribeMap: Map<string, () => void> = new Map();

  constructor() {
    super(5); 
    this.setupEventListening();
  }

  private setupEventListening(): void {
    const boosterEvents = [
      BoosterEvent.BatchQueued,
      BoosterEvent.Cancelled,
      BoosterEvent.Error,
      BoosterEvent.Failed,
      BoosterEvent.Processing,
      BoosterEvent.Queued,
      BoosterEvent.Success,
    ];

    boosterEvents.forEach(eventType => {
      const unsub = Events.On(eventType, (data: any) => {
        this.handleEvent(eventType, data);
      });
      this.eventUnsubscribeMap.set(eventType, unsub);
    });
  }

  private handleEvent(eventType: string, data: BoosterEventData | BoosterBatchProgressEventData): void {
    if (this.isSingleBoosterEvent(data)) {
      this.handleSingleBoosterEvent(eventType, data);
    } else {
      this.syncWithBackend()
    }  
  }

  private handleSingleBoosterEvent(eventType: string, data: BoosterEventData): void {
    const item = this.getItem(data.BoosterID);
    if (!item) return;

    switch (eventType) {
      case BoosterEvent.Queued:
        this.updateItemStatus(data.BoosterID, 'queued');
        break;
      case BoosterEvent.Processing:
        this.updateItemStatus(data.BoosterID, 'processing');
        break;
      case BoosterEvent.Success:
        this.updateItemStatus(data.BoosterID, 'completed', 100);
        break;
      case BoosterEvent.Error:
      case BoosterEvent.Failed:
        this.updateItemStatus(data.BoosterID, 'error', item.progress, data.Error);
        break;
      case BoosterEvent.Cancelled:
        this.updateItemStatus(data.BoosterID, 'cancelled');
        break;
    }
  }

  protected async processItem(item: BoosterItem): Promise<void> {
    try {
      const result = await InitBoosterApply(item.id);
      if (!result.Success) {
        throw new Error(result.Error || result.Message || 'Failed to process booster');
      }
      console.log(`Initiated booster processing: ${item.name}, Operation ID: ${result.OperationID}`);
    } catch (error) {
      console.error(`Failed to process booster ${item.name}:`, error);
      throw error;
    }
  }

  async syncWithBackend(): Promise<void> {
    try {
      const stateOrArray: any = await GetExecutionQueueState();
      if (!stateOrArray) return;
  
      // Normaliza: aceita tanto um array direto quanto um objeto { Items, InProgress, ... }
      const itemsArray: any[] = Array.isArray(stateOrArray)
        ? stateOrArray
        : (stateOrArray.Items ?? []);
  
      // pega InProgress de forma defensiva (possui nomes possíveis)
      const inProgressCount: number = (typeof stateOrArray === 'object')
        ? (stateOrArray.InProgress ?? stateOrArray.inProgress ?? 0)
        : 0;
  
      // Mapa para busca rápida por BoosterID
      const indexByBooster = new Map<string, number>();
      itemsArray.forEach((it, idx) => {
        if (it && it.BoosterID) indexByBooster.set(it.BoosterID, idx);
      });
  
      // Determina quais BoosterIDs estão em processing (heurística: primeiros inProgressCount itens)
      const processingSet = new Set<string>();
      for (let i = 0; i < Math.min(inProgressCount, itemsArray.length); i++) {
        const id = itemsArray[i]?.BoosterID;
        if (id) processingSet.add(id);
      }
  
      // 1) Atualiza/insere itens que estão no backend
      for (const qItem of itemsArray) {
        if (!qItem || !qItem.BoosterID) continue;
        const localItem = this.getItem(qItem.BoosterID);
        if (!localItem) continue;
  
        // Infer status: se está entre os primeiros inProgressCount => processing, senão queued
        let inferredStatus: ItemStatus = processingSet.has(qItem.BoosterID) ? 'processing' : 'queued';
  
        // Se quiser, priorize informação local de operação (se sua classe expõe operations)
        // Exemplo (opcional): if (this.operations?.has(qItem.OperationID)) { inferir a partir da operação local }
  
        // Tentativa de obter progress (se backend não fornecer, mantém progress local)
        const progress = (typeof qItem.Progress === 'number') ? qItem.Progress : localItem.progress;
  
        this.updateItemStatus(qItem.BoosterID, inferredStatus, progress, qItem.Error ?? undefined);
      }
  
      // 2) Verifica itens que o frontend conhece, mas que não aparecem no backend: marcar como completed
      //    (assume-se que se sumiu da fila -> já terminou; se foi cancelado/erro remoto, backend deveria reportar)
      const allLocal = this.getItems(); // supõe-se que retorna ProcessableItemWithStatus<T>[]
      for (const local of allLocal) {
        const boosterId = local.item.id;
        if (!indexByBooster.has(boosterId)) {
          // só atualiza se estava em 'processing' ou 'queued' — senão mantém (ex: 'idle' permanece)
          if (local.status === 'processing' || local.status === 'queued') {
            this.updateItemStatus(boosterId, 'completed', 100);
          }
        }
      }
    } catch (error) {
      console.error('Failed to sync with backend:', error);
    }
  }

  private isSingleBoosterEvent(data: any): data is BoosterEventData {
    return 'BoosterID' in data;
  }
}
