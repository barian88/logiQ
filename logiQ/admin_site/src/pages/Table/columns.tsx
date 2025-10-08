import {Tag, type TableProps} from 'antd';

import type {Question} from './types.ts';
import {CATEGORY_OPTIONS, DIFFICULTY_OPTIONS, TYPE_OPTIONS} from "../../constants/option.ts";

export const questionTableColumns: TableProps<Question>['columns'] = [
    {
        title: 'Question',
        dataIndex: 'question_text',
        key: 'question_text',
        width: 350,
    },
    {
        title: 'Options',
        dataIndex: 'options',
        key: 'options',
        width: 250,
        render: (options: string[], record: Question) => {
            const labels = ['A', 'B', 'C', 'D'];
            return options.map((option, index) => (
                <div key={index} >
                    {record.type === 'trueFalse' ? option : labels[index] + ': ' + option}
                </div>
            ));
        },
    },
    {
        title: 'Answer',
        dataIndex: 'correct_answer_index',
        key: 'correct_answer_index',
        width: 110,
        render: (correct_answer_index: number[], record: Question) => {
            const choices = record.type === 'trueFalse'
                ? ['True', 'False']
                : ['A', 'B', 'C', 'D'];
            return correct_answer_index.map((answerIndex) => (
                <Tag color="green" bordered={false} key={answerIndex}>{choices[answerIndex]}</Tag>
            ));
        },
    },

    {
        title: 'Category',
        dataIndex: 'category',
        key: 'category',
        width: 110,
        render: (category: string) => (
            <Tag color="geekblue" bordered={false}>  {CATEGORY_OPTIONS.find(opt => opt.value === category)?.label || 'Unknown'}
            </Tag>
        ),
        showSorterTooltip: {target: 'full-header'},
        filters: [
            {text: CATEGORY_OPTIONS[0].label, value: CATEGORY_OPTIONS[0].value},
            {text: CATEGORY_OPTIONS[1].label, value: CATEGORY_OPTIONS[1].value},
            {text: CATEGORY_OPTIONS[2].label, value: CATEGORY_OPTIONS[2].value},
        ],
        filterMultiple: false,
        onFilter: (value, record) => record.category.indexOf(value as string) === 0,
    },
    {
        title: 'Difficulty',
        dataIndex: 'difficulty',
        key: 'difficulty',
        width: 110,
        render: (difficulty: string) => {
            let color = 'green';
            if (difficulty === 'medium') {
                color = 'orange';
            } else if (difficulty === 'hard') {
                color = 'red';
            }
            return <Tag color={color} bordered={false}>  {DIFFICULTY_OPTIONS.find(opt => opt.value === difficulty)?.label || 'Unknown'}
            </Tag>;
        },
        showSorterTooltip: {target: 'full-header'},
        filters: [
            {text: DIFFICULTY_OPTIONS[0].label, value: DIFFICULTY_OPTIONS[0].value},
            {text: DIFFICULTY_OPTIONS[1].label, value: DIFFICULTY_OPTIONS[1].value},
            {text: DIFFICULTY_OPTIONS[2].label, value: DIFFICULTY_OPTIONS[2].value},
        ],
        filterMultiple: false,
        onFilter: (value, record) => record.difficulty.indexOf(value as string) === 0,

    },
    {
        title: 'Type',
        dataIndex: 'type',
        key: 'type',
        width: 130,
        render: (type: string) => (
            <Tag color="purple" bordered={false}> {TYPE_OPTIONS.find(opt => opt.value === type)?.label || 'Unknown'}
            </Tag>
        ),
        showSorterTooltip: {target: 'full-header'},
        filters: [
            {text: TYPE_OPTIONS[0].label, value: TYPE_OPTIONS[0].value},
            {text: TYPE_OPTIONS[1].label, value: TYPE_OPTIONS[1].value},
            {text: TYPE_OPTIONS[2].label, value: TYPE_OPTIONS[2].value},
        ],
        filterMultiple: false,
        onFilter: (value, record) => record.type.indexOf(value as string) === 0,
    },
    {
        title: 'Accuracy',
        dataIndex: 'accuracy_rate',
        key: 'accuracy_rate',
        width: 100,
        render: (accuracy_rate: number, record ) => {
            if (record.total_answers === 0) {
                return <Tag color="default" bordered={false}>--</Tag>;
            }
            const percentage = (accuracy_rate * 100).toFixed(1) + '%';
            if (accuracy_rate >= 0.70) {
                return <Tag color="green" bordered={false}>{percentage}</Tag>;
            } else if (accuracy_rate >= 0.30) {
                return <Tag color="orange" bordered={false}>{percentage}</Tag>;
            }
            return <Tag color="red" bordered={false}>{percentage}</Tag>;
        },
    },
];
