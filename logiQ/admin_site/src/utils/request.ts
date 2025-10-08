import axios from 'axios';
import {AUTH_TOKEN_KEY} from '../constants/auth.ts';
const axiosInstance = axios.create({
    baseURL: 'http://localhost:3000', // 默认请求路径
    timeout: 5000, // 超时时间
});

// 添加请求拦截器
axiosInstance.interceptors.request.use(
    (config) => {
        // 从 localStorage 中获取 token
        const token = localStorage.getItem(AUTH_TOKEN_KEY);
        // 如果 token 存在，则在请求头中添加 Authorization
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }

        return config;
    },
    (error) => {
        return Promise.reject(error);
    },
);

axiosInstance.interceptors.response.use(
    (response) => response,
    (error) => {
        // 如果响应状态码是 401，表示未授权，执行登出操作
        if (error?.response?.status === 401) {
            localStorage.removeItem(AUTH_TOKEN_KEY);
            window.dispatchEvent(new Event('auth:logout'));
        }
        return Promise.reject(error);
    },
);

export default axiosInstance;
