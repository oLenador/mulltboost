import * as React from "react"
import * as ProgressPrimitive from "@radix-ui/react-progress"

import { cn } from "@/presentation/lib/utils"

interface ProgressProps extends React.ComponentPropsWithoutRef<typeof ProgressPrimitive.Root> {
  bgColor?: string
  indicatorColor?: string
}

const Progress = React.forwardRef<
  React.ElementRef<typeof ProgressPrimitive.Root>,
  ProgressProps
>(({ className, value, bgColor, indicatorColor, ...props }, ref) => (
  <ProgressPrimitive.Root
    ref={ref}
    className={cn(
      "relative h-3 w-full overflow-hidden rounded-full border border-solid border-neutral-light-0/10 outline-none ring-offset-0 ",
      bgColor || "bg-neutral-dark-200 dark:bg-neutral-50/20",
      className
    )}
    {...props}
  >
    {value !== 0 && (
      <ProgressPrimitive.Indicator
        className={cn(
          "h-full w-full flex-1 transition-all",
          indicatorColor || "bg-neutral-50 dark:bg-neutral-50"
        )}
        style={{ transform: `translateX(-${100 - (value || 0)}%)` }}
      />
    )}
  </ProgressPrimitive.Root>
))
Progress.displayName = ProgressPrimitive.Root.displayName

export { Progress }
