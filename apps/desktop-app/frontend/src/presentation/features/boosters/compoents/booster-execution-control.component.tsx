
import React from 'react';
import { RotateCcw } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { Button } from '@/presentation/components/ui/button';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/presentation/components/ui/tooltip';
import { FloatElement } from '@/presentation/components/floating-manager';
import { ExecutionStats } from '../domain/booster-queue.types';

interface BoosterExecutionControlsProps {
    hasChanges: boolean;
    onReset: () => void;
}

export const BoosterExecutionControls: React.FC<BoosterExecutionControlsProps> = ({
    hasChanges,
    onReset,
}) => {
    const { t } = useTranslation('boosters');

    if (!hasChanges) {
        return null;
    }

    return (
        <FloatElement
            id="booster-execution-controls"
            type="custom"
            position="bottom-right"
            className="right-[112px]"
        >
            <div className='flex flex-row gap-2 bg-zinc-800/50 p-3 w-fit rounded-xl border border-zinc-700/60'>

                <TooltipProvider>
                    <Tooltip>
                        <TooltipTrigger asChild>
                            <Button
                                onClick={onReset}
                                variant="zinc"
                                size="sm"
                                className="border-zinc-700/60 hover:bg-zinc-700"
                            >
                                <RotateCcw className="w-4 h-4" />
                            </Button>
                        </TooltipTrigger>
                        <TooltipContent>
                            {t('actions.resetChanges')}
                        </TooltipContent>
                    </Tooltip>
                </TooltipProvider>

            </div>
        </FloatElement>
    );
};