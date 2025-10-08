import axiosInstance from '../../utils/request.ts';

export const login = async (username: string, password: string): Promise<string | {token: string}> => {
    const res = await axiosInstance.post('/auth/login-admin', {
        email: username,
        password,
    });
    return res.data;
};
