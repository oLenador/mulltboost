import React, { createContext, lazy, Suspense, useContext, useEffect, useState } from 'react'
import PageListingLoading from '../../components/pages/PageListingLoading';
import { AuthContext } from '../middleware';
import LoadingState from '../../components/pages/loading/loadingState.component';
import { DashboardHeader } from '../../components/header/dashboard-header';
import { UserProvider, UserProviderHook } from '../../providers/user.provider';
import SettingsPage from './settings/settings-page';
import BoxesPage from './boxes/boxes-page';

const PlayPage = lazy(() => import('./play/play-page'))


interface PagesProviderI {
    currentPage: "settings" | "boxes" | "play";
    handleChagePage: (newValue: any) => void;
}

const PagesProvider_INITIAL: PagesProviderI = {
    currentPage: "play",
    handleChagePage: (newValue: any) => { },
}
export const PagesProvider = createContext<PagesProviderI>(PagesProvider_INITIAL)




export function DashboardPages() {
    const [currentPage, setPage] = useState<PagesProviderI["currentPage"]>("play");
    const { isAuthenticated } = useContext(AuthContext)
    const [isLoading, setLoading] = useState<boolean>(true)
    const handleChagePage = (newValue: any) => {
        setPage(newValue);
    };

    const UserProviderValues = UserProviderHook()
    // console.log(UserProviderValues)
    const pages = {
        play: currentPage === "play" ? (
            <Suspense fallback={<PageListingLoading />}>
                <PlayPage />
            </Suspense>)
            :
            null,
        settings: currentPage === "settings" ? (
            <Suspense fallback={<PageListingLoading />}>
                <SettingsPage />
            </Suspense>)
            :
            null,
        boxes: currentPage === "boxes" ? (
            <Suspense fallback={<PageListingLoading />}>
                <BoxesPage />
            </Suspense>)
            :
            null,

        // 'crm': currentPage === "crm" ? (
        //     <Suspense fallback={<PageListingLoading />}>
        //         {
        //             <CRMPage />
        //         }
        //     </Suspense>)
        //     :
        //     null,
    }

    return (
        <>
            <Suspense fallback={<LoadingState />}>
                {
                    <UserProvider.Provider value={UserProviderValues}>
                        <PagesProvider.Provider value={{ handleChagePage, currentPage }}>
                            <section className='flex flex-row w-screen h-[100vh] bg-neutral-dark-0 text-white overflow-hidden'>
                                <DashboardHeader />
                                <section className='flex bg-black flex-row items-center w-full h-full '>
                                    {   // Render the pages
                                        pages[currentPage]
                                    }
           
                                </section>
                            </section>
                        </PagesProvider.Provider>
                    </UserProvider.Provider>
                }
            </Suspense>
        </>
    )
}

