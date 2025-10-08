import {createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode} from 'react';

import {AUTH_TOKEN_KEY} from '../constants/auth.ts';

interface AuthContextValue {
    isAuthenticated: boolean;
    login: (token: string) => void;
    logout: () => void;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export const AuthProvider = ({children}: {children: ReactNode}) => {
    const [isAuthenticated, setIsAuthenticated] = useState<boolean>(() => Boolean(localStorage.getItem(AUTH_TOKEN_KEY)));

    useEffect(() => {
        const handleForcedLogout = () => {
            setIsAuthenticated(false);
        };

        window.addEventListener('auth:logout', handleForcedLogout);

        return () => {
            window.removeEventListener('auth:logout', handleForcedLogout);
        };
    }, []);

    useEffect(() => {
        const handleStorageChange = (event: StorageEvent) => {
            if (event.key === AUTH_TOKEN_KEY) {
                setIsAuthenticated(Boolean(event.newValue));
            }
        };

        window.addEventListener('storage', handleStorageChange);

        return () => {
            window.removeEventListener('storage', handleStorageChange);
        };
    }, []);

    const login = useCallback((token: string) => {
        localStorage.setItem(AUTH_TOKEN_KEY, token);
        setIsAuthenticated(true);
    }, []);

    const logout = useCallback(() => {
        localStorage.removeItem(AUTH_TOKEN_KEY);
        setIsAuthenticated(false);
    }, []);

    const value = useMemo(
        () => ({isAuthenticated, login, logout}),
        [isAuthenticated, login, logout],
    );

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// eslint-disable-next-line react-refresh/only-export-components
export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};
