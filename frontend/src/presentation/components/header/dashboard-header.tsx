
import DashHeaderItem from "./header-Item";
import SARSIcon from "../SARSIcon";
import React, { useContext } from "react";
import { MdAccountCircle, MdPerson } from "react-icons/md";
import { PagesProvider } from "../../pages/dashboard/dashboard";
import { Box, LucidePlay, Settings } from "lucide-react";

export function DashboardHeader() {
  const { handleChagePage, currentPage } = useContext(PagesProvider);


  return (
    <nav className="h-full flex flex-col w-fit justify-between  items-center pb-6 bg-[#181C1F] border-neutral-light-0/5 border-r-[2px] ">
      <div className="px-3 w-full flex flex-col justify-between items-center gap-8">

        <div
          onClick={() => handleChagePage("play")}
          className="cursor-pointer"
        >
          <div className="h-[3rem] py-6 pb-12">
            <SARSIcon width={56} />
          </div>
        </div>

        <hr className="w-full border-solid border-b-2 border-[#2F3335]"/>

        <div className="flex flex-col items-center gap-7 overflow-hidden">
          <DashHeaderItem icon={<LucidePlay/>} link="play" title="Jogar" />
          <DashHeaderItem icon={<Settings/>} link="settings" title="Configurações" />
          <DashHeaderItem icon={<Box/>} link="boxes" title="Caixa"  />
          {
            // <DashHeaderItem link="crm" title="CRM" />
          }
        </div>

      </div>
        <DashHeaderItem type={"link"} link="config" title="Gerenciar" icon={<MdAccountCircle size={32} />} />
    </nav>
  );
}
