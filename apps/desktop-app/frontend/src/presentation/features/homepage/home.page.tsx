import React, { Suspense, useMemo } from "react";
import { Cpu, HardDrive, Zap, Activity, Thermometer, ArrowUp, ArrowDown, Play } from "lucide-react";
import { Card, CardHeader, CardTitle, CardContent } from "@/presentation/components/ui/card";
import { Button } from "@/presentation/components/ui/button";
import { Separator } from "@/presentation/components/ui/separator";
import { ScrollArea } from "@/presentation/components/ui/scroll-area";
import { Skeleton } from "@/presentation/components/ui/skeleton";
import BasePage from '@/presentation/components/pages/base-page';
import { useMonitoring } from "@/core/hooks/use-monitoring.hook";
import { useTranslation } from "react-i18next";

/**
 * Responsiveness fixes applied:
 * - Wrap main content in a scrollable container so very-vertical viewports won't let cards overflow the visible area.
 * - Ensure cards are `h-full` / `min-h-0` so grid/flex children can shrink and inner scroll areas work.
 * - Use responsive heights (h-36 / h-44 / h-56) for chart area so it consumes less vertical space on small/tall screens.
 * - Limit heights on scrollable areas (recent activity, quick actions content) using viewport-relative `max-h-[..vh]` so they adapt to small heights.
 * - Prevent text/icon overflow with `min-w-0` and `truncate` where appropriate.
 *
 * You asked for the complete code without obfuscation — here it is.
 */

const StatCard = ({ title, icon, value, meta, changeVariant, change, loading }: any) => (
  <Card variant="zincHover" padding="sm" className="h-full min-h-0 flex flex-col">
    <CardHeader className="p-0 mb-3">
      <div className="flex items-center justify-between min-w-0">
        <div className="p-2 bg-zinc-800 rounded-lg text-zinc-400 flex-shrink-0">{icon}</div>
        <div className="flex items-center text-xs min-w-0 ml-3">
          {loading ? (
            <Skeleton className="w-14 h-4" />
          ) : changeVariant === "up" ? (
            <div className="flex items-center text-green-400 whitespace-nowrap"><ArrowUp className="w-3 h-3 mr-1" /><span>{change}</span></div>
          ) : changeVariant === "down" ? (
            <div className="flex items-center text-red-400 whitespace-nowrap"><ArrowDown className="w-3 h-3 mr-1" /><span>{change}</span></div>
          ) : (
            <span className="text-zinc-400 whitespace-nowrap">{change}</span>
          )}
        </div>
      </div>
    </CardHeader>

    <CardContent className="p-0 flex-1 min-h-0">
      <h3 className="text-sm font-medium text-zinc-400 mb-1 truncate">{title}</h3>
      <p className="text-2xl font-semibold text-zinc-100 mb-1 truncate">
        {loading ? <Skeleton className="w-20 h-8" /> : value}
      </p>
      <p className="text-xs text-zinc-500 truncate">{loading ? <Skeleton className="w-16 h-3" /> : meta}</p>
    </CardContent>
  </Card>
);

function computeChange(current?: number, previous?: number) {
  if (current == null || previous == null) return { text: "—", variant: "normal" as const };
  if (previous === 0) {
    return { text: `${current > 0 ? '+' : ''}${Math.round(current)}%`, variant: current >= 0 ? "up" as const : "down" as const };
  }
  const diff = current - previous;
  const pct = (diff / Math.abs(previous)) * 100;
  const sign = pct > 0 ? "+" : "";
  const variant = pct > 0 ? "up" : pct < 0 ? "down" : "normal";
  return { text: `${sign}${Math.round(pct)}%`, variant };
}

export default function HomePage() {
  const { t } = useTranslation("homepage");

  const { currentMetrics, metricsHistory, isLoading, clearHistory } = useMonitoring(true);

  const stats = useMemo(() => {
    const cur = currentMetrics;
    const hist = metricsHistory ?? [];

    const cpuCur = cur?.CPU?.Usage ?? (hist.length ? hist[hist.length - 1].cpu?.usage : undefined);
    const cpuPrev = hist.length >= 2 ? hist[hist.length - 2].cpu?.usage : undefined;
    const cpuChange = computeChange(cpuCur, cpuPrev);

    const ramCur = cur?.Memory?.UsagePercent ?? (hist.length ? hist[hist.length - 1].memory?.usagePercent : undefined);
    const ramPrev = hist.length >= 2 ? hist[hist.length - 2].memory?.usagePercent : undefined;
    const ramChange = computeChange(ramCur, ramPrev);

    const diskCur = cur?.Disk?.Drives?.[0]?.UsagePercent ?? (hist.length ? hist[hist.length - 1].disk?.drives?.[0]?.usagePercent : undefined);
    const diskPrev = hist.length >= 2 ? hist[hist.length - 2].disk?.drives?.[0]?.usagePercent : undefined;
    const diskChange = computeChange(diskCur, diskPrev);

    const tempCur = cur?.CPU?.Temperature ?? cur?.Temperature?.CPU ?? (hist.length ? hist[hist.length - 1].cpu?.temperature ?? hist[hist.length - 1].temperature?.cpu : undefined);
    const tempPrev = hist.length >= 2 ? hist[hist.length - 2].cpu?.temperature ?? hist[hist.length - 2].temperature?.cpu : undefined;
    const tempChange = computeChange(tempCur, tempPrev);

    return [
      {
        id: "cpu",
        icon: <Cpu />,
        title: "CPU",
        value: cpuCur != null ? `${Math.round(cpuCur)}%` : "—",
        meta: cur?.CPU ? `${cur?.CPU.CoreCount ?? "?"} ${t("cores")} • ${Math.round(cur?.CPU?.Frequency ?? 0)} MHz` : "—",
        change: cpuChange.text,
        changeVariant: cpuChange.variant,
      },
      {
        id: "ram",
        icon: <Activity />,
        title: "RAM",
        value: ramCur != null ? `${Math.round(ramCur)}%` : "—",
        meta: cur?.Memory ? `${(cur.Memory.Total ?? 0) > 0 ? `${Math.round(((cur.Memory.Used ?? 0) / (cur.Memory.Total ?? 1)) * 100)}% ${t("used")}` : "—"}` : "—",
        change: ramChange.text,
        changeVariant: ramChange.variant,
      },
      {
        id: "ssd",
        icon: <HardDrive />,
        title: "SSD",
        value: diskCur != null ? `${Math.round(diskCur)}%` : "—",
        meta: cur?.Disk?.Drives?.[0] ? `${Math.round((cur.Disk.Drives[0].Used ?? 0) / 1024 / 1024 / 1024)}GB ${t("usedPlural")}` : "—",
        change: diskChange.text,
        changeVariant: diskChange.variant,
      },
      {
        id: "temp",
        icon: <Thermometer />,
        title: t("temp"),
        value: tempCur != null ? `${Math.round(tempCur)}°C` : "—",
        meta: cur?.GPU ? `${cur.GPU?.Name ?? ""}` : (cur?.Temperature ? t("systemTemperatures") : "—"),
        change: tempChange.text,
        changeVariant: tempChange.variant,
      },
    ];
  }, [currentMetrics, metricsHistory, t]);

  const recent = [
    { action: t("recentActions.fpsBoost"), time: "2 min", status: "success" },
    { action: t("recentActions.cacheClean"), time: "15 min", status: "success" },
    { action: t("recentActions.networkAnalysis"), time: "1h", status: "info" },
    { action: t("recentActions.gameOptimization"), time: "2h", status: "success" },
  ];

  return (
    <BasePage>
        <div className="max-w-4xl space-y-8">
                  <div className="flex flex-col gap-6 p-4 min-h-[calc(100vh-3.5rem)] max-h-[100vh] overflow-auto">
        <header>
          <h1 className="text-2xl font-semibold">{t("title")}</h1>
          <p className="text-zinc-400 text-sm">{t("subTitle")}</p>
        </header>

        <img src="" alt="" />

        {/* Stats grid: ensure children stretch and can shrink (h-full, min-h-0) */}
        <div className="w-full grid grid-cols-1 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-4 gap-4 items-stretch">
          {stats.map((s) => (
            <div key={s.id} className="min-w-0">
              <StatCard
                item={s}
                title={s.title}
                icon={s.icon}
                value={s.value}
                meta={s.meta}
                change={s.change}
                changeVariant={s.changeVariant}
                loading={isLoading}
              />
            </div>
          ))}
        </div>

        {/* Performance + Quick Actions */}
        <div className="flex w-full gap-6 items-start">
          <Suspense fallback={<Skeleton />}>
            <Card variant="zinc" padding="default" className="w-full lg:col-span-2 h-full min-h-0 flex flex-col">
              <CardHeader className="p-0 mb-4 flex items-center justify-between">
                <CardTitle className="text-lg">{t("performaceCard.title")}</CardTitle>
                <div className="flex space-x-1">
                  <Button size="sm" variant="zinc">{t("performaceCard.optionSelector.1h")}</Button>
                  <Button size="sm" variant="ghost" className="text-zinc-500">{t("performaceCard.optionSelector.6h")}</Button>
                  <Button size="sm" variant="ghost" className="text-zinc-500">{t("performaceCard.optionSelector.1D")}</Button>
                </div>
              </CardHeader>

              <CardContent className="p-0 flex-1 min-h-0">
                {/* Use responsive heights and allow shrinking on short viewports */}
                <div className="h-36 sm:h-44 md:h-56 lg:h-48 bg-zinc-800 rounded-xl flex items-center justify-center border border-zinc-700 w-full min-h-0">
                  {isLoading ? (
                    <p className="text-zinc-500 text-sm">{t("loadingMetrics")}</p>
                  ) : (
                    <p className="text-zinc-500 text-sm text-center px-4">{t("performaceCard.noDataFallback", { count: metricsHistory?.length ?? 0 })}</p>
                  )}
                </div>
              </CardContent>
            </Card>
          </Suspense>
        </div>

        {/* Recent Activity */}
        <Suspense fallback={<Skeleton />}>
          <Card variant="zinc" padding="default" className="h-full min-h-0 flex flex-col">
            <CardHeader className="p-0 mb-4">
              <CardTitle className="text-lg">{t("recentActivityCard.title")}</CardTitle>
            </CardHeader>
            <Separator />
            <CardContent className="p-0 mt-4 flex-1 min-h-0">
              {/* Make this scrollable with a max-height relative to viewport to avoid pushing content off-screen */}
              <ScrollArea className="max-h-[28vh] sm:max-h-[32vh] md:max-h-[36vh]">
                <div className="space-y-3 pr-4">
                  {recent.map((item, idx) => (
                    <div key={idx} className="flex items-center justify-between p-3 bg-zinc-800/50 rounded-lg border border-zinc-700/50">
                      <div className="flex items-center space-x-3 min-w-0">
                        <div className={`w-2 h-2 rounded-full ${item.status === "success" ? "bg-green-400" : "bg-blue-400"}`} />
                        <span className="text-sm font-medium text-zinc-300 truncate">{item.action}</span>
                      </div>
                      <span className="text-xs text-zinc-500">{item.time}</span>
                    </div>
                  ))}
                </div>
              </ScrollArea>
            </CardContent>
          </Card>
        </Suspense>
      </div>
      </div>
    </BasePage>
  );
}
