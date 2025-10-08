import axiosInstance from '../../utils/request.ts';

export interface DimensionDistributions {
    difficulty: DimensionPortion[];
    category: DimensionPortion[];
    type: DimensionPortion[];
}

export interface DimensionPortion {
    value: string;    // 例如："truthTable" 或 "easy"
    count: number;    // 该维度下的问题数量
    portion: number;  // 占总问题的比例，例如 0.25 表示 25%
}

export const getDimensionDistributions = async (): Promise<DimensionDistributions> => {
    const res = await axiosInstance.get('/question-stats/dimension-distribution');
    return res.data;
}

export interface AccuracyByDimension {
    value: string; // 维度值，例如 "truthTable" 或 "easy"
    accuracy: number;
}

export interface AccuracyDistributions {
    difficulty: AccuracyByDimension[];
    category: AccuracyByDimension[];
    type: AccuracyByDimension[];
}

export const getAccuracyDistributions = async (): Promise<AccuracyDistributions> => {
    const res = await axiosInstance.get('/question-stats/dimension-accuracy');
    return res.data;
}
