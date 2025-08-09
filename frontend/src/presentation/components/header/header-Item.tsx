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
                <button
                    onClick={() => type === "path" ? handleChagePage(link) : handleOpenLink()}
                    className="rounded-md w-full flex flex-row items-center justify-start py-3  px-4 text-lg
                    data-[active=true]:z-50 data-[active=false]:text-neutral-light-0/60 data-[active=true]:bg-neutral-light-0/5 "
                    data-active={isActive}
                    disabled={isActive}
                >
                    <div 
                        className={
                            "text-nowrap " +
                            "data-[has-icon=true]:flex items-center flex-row data-[has-icon=true]:gap-4 data-[has-icon=true]:justify-center" +
                            "data-[active=false]:rounded-lg data-[active=false]:hover:text-neutral-light-0/60" +
                            "min-w-5"
                        }
                        data-has-icon={!!icon}
                        data-active={isActive}
                    >
                        {icon && icon}

                        <span className="-mb-1 h-fit">{title}</span>
                    </div>
                </button>
    );
}
