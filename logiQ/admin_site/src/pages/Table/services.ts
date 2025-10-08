import axiosInstance from '../../utils/request.ts';
import type {Question} from './types.ts';

export interface ApiResponse {
    list: Question[];
    total: number;
    current_page: number;
}

export const fetchTableData = async (
    page: number,
    pageSize: number,
    category: string,
    difficulty: string,
    type: string,
): Promise<ApiResponse> => {
    const res = await axiosInstance.get('/question/questions', {
        params: {
            page,
            page_size: pageSize,
            category,
            difficulty,
            type,
        },
    });
    return res.data;
};

export const deleteQuestions = async (ids: string[]): Promise<void> => {
    await axiosInstance.post('/question/delete', {
        ids,
    });
};
