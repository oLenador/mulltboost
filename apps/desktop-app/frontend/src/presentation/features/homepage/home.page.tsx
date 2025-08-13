// pages/HomePage.tsx
import React, { Suspense, useMemo } from "react";
import { Cpu, HardDrive, Zap, Activity, Thermometer, ArrowUp, ArrowDown, Play } from "lucide-react";
import { Card, CardHeader, CardTitle, CardContent } from "@/presentation/components/ui/card";
import { Button } from "@/presentation/components/ui/button";
import { Separator } from "@/presentation/components/ui/separator";
import { ScrollArea } from "@/presentation/components/ui/scroll-area";
import { Skeleton } from "@/presentation/components/ui/skeleton";
import BasePage from '@/presentation/components/pages/base-page';
import { SystemMetrics } from "@/core/api/types";
import { useMonitoring } from "@/core/hooks/use-monitoring.hook";

const StatCard = ({ title, icon, value, meta, changeVariant, change, loading }: any) => (
  <Card variant="zincHover" padding="sm">
    <CardHeader className="p-0 mb-3">
      <div className="flex items-center justify-between">
        <div className="p-2 bg-zinc-800 rounded-lg text-zinc-400">{icon}</div>
        <div className="flex items-center text-xs">
          {loading ? (
            <Skeleton className="w-14 h-4" />
          ) : changeVariant === "up" ? (
            <div className="flex items-center text-green-400"><ArrowUp className="w-3 h-3 mr-1" /><span>{change}</span></div>
          ) : changeVariant === "down" ? (
            <div className="flex items-center text-red-400"><ArrowDown className="w-3 h-3 mr-1" /><span>{change}</span></div>
          ) : (
            <span className="text-zinc-400">{change}</span>
          )}
        </div>
      </div>
    </CardHeader>
    <CardContent className="p-0">
      <h3 className="text-sm font-medium text-zinc-400 mb-1">{title}</h3>
      <p className="text-2xl font-semibold text-zinc-100 mb-1">
        {loading ? <Skeleton className="w-20 h-8" /> : value}
      </p>
      <p className="text-xs text-zinc-500">{loading ? <Skeleton className="w-16 h-3" /> : meta}</p>
    </CardContent>
  </Card>
);

function computeChange(current?: number, previous?: number) {
  if (current == null || previous == null) return { text: "—", variant: "normal" as const };
  // avoid divide by zero
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
  // autoStart = true para já começar a receber métricas
  const { currentMetrics, metricsHistory, isLoading, clearHistory } = useMonitoring(true);

  // derive stats from currentMetrics + history
  const stats = useMemo(() => {
    const cur = currentMetrics
    const hist = metricsHistory ?? [];

    // helper to get last & previous for a numeric path
    const last = (selector: (m: SystemMetrics) => number | undefined) => {
      const n = hist.length;
      const curVal = n >= 1 ? selector(hist[n - 1]) : undefined;
      const prevVal = n >= 2 ? selector(hist[n - 2]) : undefined;
      return { curVal, prevVal };
    };

    // CPU usage
    const cpuCur = cur?.CPU?.Usage ?? (hist.length ? hist[hist.length - 1].cpu?.usage : undefined);
    const cpuPrev = hist.length >= 2 ? hist[hist.length - 2].cpu?.usage : undefined;
    const cpuChange = computeChange(cpuCur, cpuPrev);

    // RAM usage percent (assumes Memory.UsagePercent is percent already)
    const ramCur = cur?.Memory?.UsagePercent ?? (hist.length ? hist[hist.length - 1].memory?.usagePercent : undefined);
    const ramPrev = hist.length >= 2 ? hist[hist.length - 2].memory?.usagePercent : undefined;
    const ramChange = computeChange(ramCur, ramPrev);

    // Disk (first drive) usagePercent
    const diskCur = cur?.Disk?.Drives?.[0]?.UsagePercent ?? (hist.length ? hist[hist.length - 1].disk?.drives?.[0]?.usagePercent : undefined);
    const diskPrev = hist.length >= 2 ? hist[hist.length - 2].disk?.drives?.[0]?.usagePercent : undefined;
    const diskChange = computeChange(diskCur, diskPrev);

    // Temperature: prefer CPU temperature from CPU struct, fallback to Temperature.CPU
    const tempCur = cur?.CPU?.Temperature ?? cur?.Temperature?.CPU ?? (hist.length ? hist[hist.length - 1].cpu?.temperature ?? hist[hist.length - 1].temperature?.cpu : undefined);
    const tempPrev = hist.length >= 2 ? hist[hist.length - 2].cpu?.temperature ?? hist[hist.length - 2].temperature?.cpu : undefined;
    const tempChange = computeChange(tempCur, tempPrev);

    return [
      {
        id: "cpu",
        icon: <Cpu />,
        title: "CPU",
        value: cpuCur != null ? `${Math.round(cpuCur)}%` : "—",
        meta: cur?.CPU ? `${cur?.CPU.CoreCount ?? "?"} cores • ${Math.round(cur?.CPU?.Frequency ?? 0)} MHz` : "—",
        change: cpuChange.text,
        changeVariant: cpuChange.variant,
      },
      {
        id: "ram",
        icon: <Activity />,
        title: "RAM",
        value: ramCur != null ? `${Math.round(ramCur)}%` : "—",
        meta: cur?.Memory ? `${(cur.Memory.Total ?? 0) > 0 ? `${Math.round(((cur.Memory.Used ?? 0) / (cur.Memory.Total ?? 1)) * 100)}% usado` : "—"}` : "—",
        change: ramChange.text,
        changeVariant: ramChange.variant,
      },
      {
        id: "ssd",
        icon: <HardDrive />,
        title: "SSD",
        value: diskCur != null ? `${Math.round(diskCur)}%` : "—",
        meta: cur?.Disk?.Drives?.[0] ? `${Math.round((cur.Disk.Drives[0].Used ?? 0) / 1024 / 1024 / 1024)}GB usados` : "—",
        change: diskChange.text,
        changeVariant: diskChange.variant,
      },
      {
        id: "temp",
        icon: <Thermometer />,
        title: "Temp",
        value: tempCur != null ? `${Math.round(tempCur)}°C` : "—",
        meta: cur?.GPU ? `${cur.GPU?.Name ?? ""}` : (cur?.Temperature ? "Temperaturas do sistema" : "—"),
        change: tempChange.text,
        changeVariant: tempChange.variant,
      },
    ];
  }, [currentMetrics, metricsHistory]);

  const recent = [
    { action: "FPS Boost aplicado", time: "2 min atrás", status: "success" },
    { action: "Limpeza de cache realizada", time: "15 min atrás", status: "success" },
    { action: "Análise de rede concluída", time: "1h atrás", status: "info" },
    { action: "Otimização de jogos executada", time: "2h atrás", status: "success" },
  ];

  return (
    <BasePage>
      <>
        <header>
          <h1 className="text-2xl font-semibold">Dashboard</h1>
          <p className="text-zinc-400 text-sm">Visão geral do desempenho do sistema</p>
        </header>

        {/* System Stats */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {stats.map((s) => (
            <StatCard
              key={s.id}
              item={s}
              title={s.title}
              icon={s.icon}
              value={s.value}
              meta={s.meta}
              change={s.change}
              changeVariant={s.changeVariant}
              loading={isLoading}
            />
          ))}
        </div>

        {/* Performance + Quick Actions */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Suspense fallback={<Skeleton />}>
            <Card variant="zinc" padding="default" className="lg:col-span-2">
              <CardHeader className="p-0 mb-4 flex items-center justify-between">
                <CardTitle className="text-lg">Performance</CardTitle>
                <div className="flex space-x-1">
                  <Button size="sm" variant="zinc">1H</Button>
                  <Button size="sm" variant="ghost" className="text-zinc-500">6H</Button>
                  <Button size="sm" variant="ghost" className="text-zinc-500">1D</Button>
                </div>
              </CardHeader>
              <CardContent className="p-0">
                <div className="h-48 bg-zinc-800 rounded-xl flex items-center justify-center border border-zinc-700">
                  {isLoading ? (
                    <p className="text-zinc-500 text-sm">Carregando métricas...</p>
                  ) : (
                    <p className="text-zinc-500 text-sm">Gráfico de Performance (histórico: {metricsHistory?.length ?? 0} pontos)</p>
                  )}
                </div>
              </CardContent>
            </Card>
          </Suspense>
          <Suspense fallback={<Skeleton />}>
            <Card variant="zinc" padding="default">
              <CardHeader className="p-0 mb-4">
                <CardTitle className="text-lg">Ações Rápidas</CardTitle>
              </CardHeader>
              <CardContent className="p-0 space-y-3">
                <Button variant="zinc" className="w-full justify-start px-3 py-3">
                  <Zap className="mr-3" /> Otimizar Agora
                </Button>
                <Button variant="zincLight" className="w-full justify-start px-3 py-3">
                  <Activity className="mr-3" /> Análise Completa
                </Button>
                <Button variant="zincLight" className="w-full justify-start px-3 py-3" onClick={() => clearHistory()}>
                  <Play className="mr-3" /> Limpar Histórico
                </Button>
              </CardContent>
            </Card>
          </Suspense>
        </div>

        {/* Recent Activity */}
        <Suspense fallback={<Skeleton />}>
          <Card variant="zinc" padding="default">
            <CardHeader className="p-0 mb-4">
              <CardTitle className="text-lg">Atividade Recente</CardTitle>
            </CardHeader>
            <Separator />
            <CardContent className="p-0 mt-4">
              <ScrollArea className="h-56">
                <div className="space-y-3 pr-4">
                  {recent.map((item, idx) => (
                    <div key={idx} className="flex items-center justify-between p-3 bg-zinc-800/50 rounded-lg border border-zinc-700/50">
                      <div className="flex items-center space-x-3">
                        <div className={`w-2 h-2 rounded-full ${item.status === "success" ? "bg-green-400" : "bg-blue-400"}`} />
                        <span className="text-sm font-medium text-zinc-300">{item.action}</span>
                      </div>
                      <span className="text-xs text-zinc-500">{item.time}</span>
                    </div>
                  ))}
                </div>
              </ScrollArea>
            </CardContent>
          </Card>
        </Suspense>
      </>
    </BasePage>
  );
}
