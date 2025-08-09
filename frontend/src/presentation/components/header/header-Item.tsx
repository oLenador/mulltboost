import React, { memo, useContext, useEffect, useState } from "react";
import { PagesProvider, PageType } from "../../pages/dashboard/dashboard";

interface DashHeaderItemI {
    link: PageType;
    title: string;
    icon: React.ReactNode;
    type?: "path" | "link"
}

function DashHeaderItem({ link, title, icon, type = "path" }: DashHeaderItemI): React.ReactElement {
    const { handleChangePage, currentPage } = useContext(PagesProvider);
    const [ isActive, setIsActive ] = useState(currentPage === link)

    function handleClick() {
        setIsActive(true)
        handleChangePage(link)
    }
    useEffect(() => {
        setIsActive(currentPage === link)

    }, [currentPage])

    return (
                <button
                    onClick={handleClick}
                    className="rounded-md w-full flex flex-row items-center justify-start py-3 px-4 text-sm
                    data-[active=true]:z-50 data-[active=false]:text-neutral-light-0/60 data-[active=true]:bg-neutral-light-0/10 "
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

export default memo(DashHeaderItem)