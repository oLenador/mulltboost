// src/presentation/features/boosters/components/VirtualizedBoosterList.tsx

import React, { useMemo } from 'react';
import { WindowScroller, AutoSizer, List } from 'react-virtualized';
import 'react-virtualized/styles.css';
import { BoosterItem } from '../types/booster.types';
import { BoosterExecution } from '../domain/booster-queue.types';
import { BoosterCard } from './base-page/booster-card.component';

interface EffectiveBoosterItem extends BoosterItem {
  effectiveState: 'original' | 'staged-apply' | 'staged-revert';
  hasChanges: boolean;
}

interface VirtualizedBoosterListProps {
  boosters: EffectiveBoosterItem[];
  executions: Map<string, BoosterExecution>;
  onToggleBooster: (id: string) => void;
  height?: number;
  itemHeight?: number;
}

export const VirtualizedBoosterList: React.FC<VirtualizedBoosterListProps> = ({
  boosters,
  executions,
  onToggleBooster,
  height = 600,
  itemHeight = 200,
}) => {
  if (!boosters || boosters.length === 0) return null;

  // memorize data to avoid useless re-renders when passed to child components (optional)
  const memoBoosters = useMemo(() => boosters, [boosters]);
  const memoExecutions = useMemo(() => executions, [executions]);

  return (
    <div className="w-full">
      <WindowScroller>
        {({ height: windowScrollerHeight, isScrolling, onChildScroll, registerChild, scrollTop }) => (
          // registerChild é uma callback do WindowScroller — passamos pro container
          <div ref={registerChild as any}>
            <AutoSizer disableHeight>
              {({ width }) => {
                const itemsPerRow = width >= 1024 ? 2 : 1; // breakpoint lg
                const rowCount = Math.ceil(memoBoosters.length / itemsPerRow);
                const rowHeight = itemHeight;

                // rowRenderer fechado com acesso ao itemsPerRow / dados
                const rowRenderer = ({ index, key, style }: { index: number; key: string; style: React.CSSProperties }) => {
                  const startIndex = index * itemsPerRow;
                  const rowBoosters = memoBoosters.slice(startIndex, startIndex + itemsPerRow);

                  return (
                    <div key={key} style={style} className="flex gap-4 items-start">
                      {rowBoosters.map((booster) => (
                        <div key={booster.id} className="flex-1">
                          <BoosterCard
                            booster={booster}
                            execution={memoExecutions.get(booster.id)}
                            onToggle={onToggleBooster}
                          />
                        </div>
                      ))}

                      {/* Preenche colunas vazias para manter altura/align */}
                      {Array.from({ length: itemsPerRow - rowBoosters.length }).map((_, emptyIndex) => (
                        <div key={`empty-${emptyIndex}`} className="flex-1" />
                      ))}
                    </div>
                  );
                };

                return (
                  <List
                    autoHeight
                    width={width}
                    height={Math.min(windowScrollerHeight || height, rowCount * rowHeight)}
                    rowCount={rowCount}
                    rowHeight={rowHeight}
                    rowRenderer={rowRenderer}
                    scrollTop={scrollTop}
                    isScrolling={isScrolling}
                    onScroll={onChildScroll}
                    overscanRowCount={3}
                    className="scrollbar-thin scrollbar-track-zinc-800 scrollbar-thumb-zinc-600"
                  />
                );
              }}
            </AutoSizer>
          </div>
        )}
      </WindowScroller>
    </div>
  );
};
