// pages/HomePage.tsx
import React, { Suspense } from "react"
import { Cpu, HardDrive, Zap, Activity, Thermometer, ArrowUp, ArrowDown, Play } from "lucide-react"
import { Card, CardHeader, CardTitle, CardContent } from "@/presentation/components/ui/card"
import { Button } from "@/presentation/components/ui/button"
import { Separator } from "@/presentation/components/ui/separator"
import { ScrollArea } from "@/presentation/components/ui/scroll-area"
import { Skeleton } from "@/presentation/components/ui/skeleton"
import BasePage from '@/presentation/components/pages/base-page';

const stats = [
  { id: "cpu", icon: <Cpu />, title: "CPU", value: "23%", meta: "Intel i7-12700K", change: "+12%", changeVariant: "up" },
  { id: "ram", icon: <Activity />, title: "RAM", value: "67%", meta: "16GB DDR4", change: "+8%", changeVariant: "up" },
  { id: "ssd", icon: <HardDrive />, title: "SSD", value: "45%", meta: "1TB NVMe", change: "-3%", changeVariant: "down" },
  { id: "temp", icon: <Thermometer />, title: "Temp", value: "64°C", meta: "RTX 4080", change: "Normal", changeVariant: "normal" },
]

const StatCard = ({ item }: any) => (
  <Card variant="zincHover" padding="sm">
    <CardHeader className="p-0 mb-3">
      <div className="flex items-center justify-between">
        <div className="p-2 bg-zinc-800 rounded-lg text-zinc-400">{item.icon}</div>
        <div className="flex items-center text-xs">
          {item.changeVariant === "up" && <div className="flex items-center text-green-400"><ArrowUp className="w-3 h-3 mr-1" /><span>{item.change}</span></div>}
          {item.changeVariant === "down" && <div className="flex items-center text-red-400"><ArrowDown className="w-3 h-3 mr-1" /><span>{item.change}</span></div>}
          {item.changeVariant === "normal" && <span className="text-zinc-400">{item.change}</span>}
        </div>
      </div>
    </CardHeader>
    <CardContent className="p-0">
      <h3 className="text-sm font-medium text-zinc-400 mb-1">{item.title}</h3>
      <p className="text-2xl font-semibold text-zinc-100 mb-1">{item.value}</p>
      <p className="text-xs text-zinc-500">{item.meta}</p>
    </CardContent>
  </Card>
)

export default function HomePage() {
  const recent = [
    { action: "FPS Boost aplicado", time: "2 min atrás", status: "success" },
    { action: "Limpeza de cache realizada", time: "15 min atrás", status: "success" },
    { action: "Análise de rede concluída", time: "1h atrás", status: "info" },
    { action: "Otimização de jogos executada", time: "2h atrás", status: "success" },
  ]

  return (
    <BasePage>
      <>
        <header>
          <h1 className="text-2xl font-semibold">Dashboard</h1>
          <p className="text-zinc-400 text-sm">Visão geral do desempenho do sistema</p>
        </header>

        {/* System Stats */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {stats.map((s) => <StatCard key={s.id} item={s} />)}
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
                  <p className="text-zinc-500 text-sm">Gráfico de Performance</p>
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
                <Button variant="zincLight" className="w-full justify-start px-3 py-3">
                  <Play className="mr-3" /> Programar Limpeza
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
  )
}
