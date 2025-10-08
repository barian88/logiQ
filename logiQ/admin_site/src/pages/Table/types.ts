export interface Question {
    _id: string;
    question_text: string;
    options: string[];
    correct_answer_index: number[];
    category: string;
    difficulty: string;
    type: string;
    is_active: boolean;

    // Statistics fields
    total_answers: number;
    correct_answers: number;
    accuracy_rate: number;
}

export interface TableFilter {
    category?: string;
    difficulty?: string;
    type?: string;
}

export type RawTableFilter = Record<string, (string | number | boolean)[] | null>;
