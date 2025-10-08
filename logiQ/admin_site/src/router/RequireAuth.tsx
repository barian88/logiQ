import {Navigate, useLocation} from 'react-router-dom';
import type {ReactElement} from 'react';

import {useAuth} from '../context/AuthContext.tsx';

interface RequireAuthProps {
    children: ReactElement;
}

const RequireAuth = ({children}: RequireAuthProps) => {
    const {isAuthenticated} = useAuth();
    const location = useLocation();

    if (!isAuthenticated) {
        return <Navigate to="/login" state={{from: location}} replace />;
    }

    return children;
};

export default RequireAuth;
