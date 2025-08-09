
import React, { createContext, useState, useEffect, ReactNode } from 'react';
import { CheckLoginUseCase } from '../../domain/usecases/auth/check-login.usecase';


interface AuthContextValue {
    isAuthenticated: boolean;
    setIsAuthenticated: (value: boolean) => void;
    loading: boolean;
}

const AuthContext_Default: AuthContextValue = {
    isAuthenticated: false,
    setIsAuthenticated: (value: boolean) => { },
    loading: false,
};

export const AuthContext = createContext<AuthContextValue>(AuthContext_Default);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [isAuthenticated, setIsAuthenticated] = useState(true);

    async function checkLogin() {
        
        let response = await CheckLoginUseCase();

        if (response.success) {
            setIsAuthenticated(true);
        } else {
            setIsAuthenticated(false);
            // chrome.storage.local.remove("authentication");
        }
    }

    async function initMiddleware() {
        const userAuthResponse: any | undefined = await chrome.storage.local.get("authentication");
        const authentication = userAuthResponse?.authentication;

        if (!!authentication?.refreshToken) setIsAuthenticated(true)
    }

    useEffect(() => {
        initMiddleware()
        checkLogin();
        monitorUser();
    }, []);

    async function monitorUser() {
        await new Promise( r => setTimeout(r, 30000))
        setInterval(checkLogin, 300000)
    }   

    // useStorageListener(checkLogin, 'authentication');

    return (
        <AuthContext.Provider value={{ ...AuthContext_Default, isAuthenticated, setIsAuthenticated }}>
            {
            children
            }
        </AuthContext.Provider>
    );
};

