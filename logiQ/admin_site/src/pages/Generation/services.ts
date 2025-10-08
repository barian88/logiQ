import axiosInstance from '../../utils/request.ts';
import type {Question} from '../Table/types.ts';

export interface GenerateParam {
    number: number;
    category?: string;
    difficulty?: string;
    type?: string;
}

export const generateQuestions = async ({number, category, difficulty, type}: GenerateParam): Promise<Question[]> => {
    const res = await axiosInstance.post('/question/generate', {
        number,
        category,
        difficulty,
        type,
    });

    const data = res?.data;

    if (Array.isArray(data)) {
        return data as Question[];
    }

    if (Array.isArray(data?.questions)) {
        return data.questions as Question[];
    }

    if (Array.isArray(data?.list)) {
        return data.list as Question[];
    }

    return [];
};
