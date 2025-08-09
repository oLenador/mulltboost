import DashHeaderItem from "./header-Item";
import SARSIcon from "../SARSIcon";
import React, { useContext } from "react";
import { PagesProvider, PageType } from "../../pages/dashboard/dashboard";
import {
  Box,
  LayoutDashboard,
  Settings,
  Zap,
  Gauge,
  Wifi,
  Crosshair,
  Dices,
  MessagesSquare,
  BarChart3,
  UserCog
} from "lucide-react";

export function DashboardHeader() {
  const { currentPage, handleChangePage } = useContext(PagesProvider);

  return (
    <nav className="h-full flex flex-col w-fit justify-between items-center pb-6 border-neutral-light-0/5 border-r-[2px]">
      <div className="pl-1 pr-2 w-full flex flex-col justify-between items-center gap-8">

        {/* Logo */}
        <div
          onClick={() => handleChangePage(PageType.HOMEPAGE)}
          className="cursor-pointer"
        >
          <div className="h-[3rem] py-6 pb-12">
            <SARSIcon width={56} />
          </div>
        </div>

        {/* Itens */}
        <div className="flex flex-col items-start gap-4 overflow-hidden">

          {/* Geral */}
          <div className="flex flex-col gap-2">
            <span className="text-xs font-light text-neutral-light-0/40 ml-3 pl-1 pt-2 uppercase">Geral</span>
            <div>
              <DashHeaderItem icon={<LayoutDashboard size={16} />} link={PageType.HOMEPAGE} title="Hub" />
              <DashHeaderItem icon={<Settings size={16} />} link={PageType.SETTINGS} title="Configurações" />
            </div>
          </div>

          {/* Otimizações */}
          <div className="flex flex-col gap-2">
            <span className="text-xs font-light text-neutral-light-0/40 ml-3 pl-1 pt-2 uppercase">Otimizações</span>
            <div>
              <DashHeaderItem icon={<Gauge size={16} />} link={PageType.FPS_BOOST} title="FPS Boost" />
              <DashHeaderItem icon={<Wifi size={16} />} link={PageType.CONNECTION} title="Conexão" />
              <DashHeaderItem icon={<Crosshair size={16} />} link={PageType.PRECISION} title="Precisão" />
              <DashHeaderItem icon={<Dices size={16} />} link={PageType.GAMES} title="Games" />
              <DashHeaderItem icon={<Zap size={16} />} link={PageType.FLUSHER} title="Flusher" />
            </div>
          </div>

          {/* Multi AI */}
          <div className="flex flex-col gap-2">
            <span className="text-xs font-light text-neutral-light-0/40 ml-3 pl-1 pt-2 uppercase">Multi AI</span>
            <DashHeaderItem icon={<MessagesSquare size={16} />} link={PageType.CHAT} title="Chat" />
            <DashHeaderItem icon={<BarChart3 size={16} />} link={PageType.ANALYTICS} title="Análises" />
            <DashHeaderItem icon={<Zap size={16} />} link={PageType.SMART_BOOST} title="Smart Boost" />
          </div>
        </div>
      </div>

      {/* Rodapé */}
      <DashHeaderItem
        type={"link"}
        link={PageType.SETTINGS}
        title="Gerenciar"
        icon={<UserCog size={16} />}
      />
    </nav>
  );
}
