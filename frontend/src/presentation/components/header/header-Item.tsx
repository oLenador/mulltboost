import React, { useContext } from "react";
import { PagesProvider } from "../../pages/dashboard/dashboard";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "../ui/tooltip";

interface DashHeaderItemI {
    link: string;
    title: string;
    icon: React.ReactNode;
    type?: "path" | "link"
}

export default function DashHeaderItem({ link, title, icon, type = "path" }: DashHeaderItemI): React.ReactElement {
    const { handleChagePage, currentPage } = useContext(PagesProvider);
    const isActive = currentPage === link;

    function handleOpenLink() {

    }

    return (
        <TooltipProvider delayDuration={500}>
        <Tooltip >
            <TooltipTrigger asChild>
                <button
                    onClick={() => type === "path" ? handleChagePage(link) : handleOpenLink()}
                    className="rounded-lg w-fit bg-[#1A1B1F] flex flex-row items-center justify-center py-4 px-4 text-lg !aspect-square data-[active=true]:border-white 
                    border-2 data-[active=true]:z-50 data-[active=false]:text-neutral-light-0/60 data-[active=false]:border-white/10"
                    data-active={isActive}
                    disabled={isActive}
                >
                    <div 
                        className={
                            "text-nowrap " +
                            "data-[has-icon=true]:flex data-[has-icon=true]:flex-row data-[has-icon=true]:gap-4 data-[has-icon=true]:items-center data-[has-icon=true]:justify-center data-[has-icon=true]:-my-[14px] " +
                            "data-[active=false]:rounded-lg data-[active=false]:hover:bg-neutral-light-0/15 data-[active=false]:hover:text-neutral-light-0/60" +
                            "min-w-5 max-w-5 w-5 min-h-5 max-h-5 h-5"
                        }
                        data-has-icon={!!icon}
                        data-active={isActive}
                    >
                        {icon && icon}
                    </div>
                </button>
            </TooltipTrigger>
            <TooltipContent>
                <p className="text-base">{title}</p>
            </TooltipContent>
        </Tooltip>
    </TooltipProvider>
    );
}
