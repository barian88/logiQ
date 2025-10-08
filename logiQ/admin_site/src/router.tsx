import {createBrowserRouter, Navigate} from 'react-router-dom';
import MyLayout from './components/Layout/layout.tsx';
import LoginPage from './pages/Login/login.tsx';
import TablePage from './pages/Table/table.tsx';
import GenerationPage from './pages/Generation/generation.tsx';
import RequireAuth from './router/RequireAuth.tsx';
import StatisticsPage from "./pages/Statistics/statistics.tsx";

export const router = createBrowserRouter([
    {
        path: '/login',
        element: <LoginPage />,
    },
    {
        path: '/',
        element: (
            <RequireAuth>
                <MyLayout />
            </RequireAuth>
        ),
        children: [
            {
                index: true,
                element: <Navigate to="table" replace />,
            },
            {
                path: 'table',
                element: <TablePage />,
            },
            {
                path: 'generation',
                element: <GenerationPage />,
            },
            {
                path: 'statistics',
                element: <StatisticsPage />,
            },
        ],
    },
]);
