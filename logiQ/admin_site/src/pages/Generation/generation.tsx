import {useCallback, useEffect, useMemo, useState} from 'react';
import {Button, Card, Col, Form, InputNumber, Row, Select, Space, Spin, Typography, message} from 'antd';

import type {Question} from '../Table/types.ts';
import GenerationPreview, {type GenerationSummary} from './components/GenerationPreview.tsx';
import {createJsonDownloadHandle, type JsonDownloadHandle} from '../../utils/download.ts';
import {CATEGORY_OPTIONS, DIFFICULTY_OPTIONS, TYPE_OPTIONS} from '../../constants/option.ts';
import {generateQuestions, type GenerateParam} from './services.ts';
import './generation.css';

const {Title, Paragraph} = Typography;

const UNSPECIFIED_KEY = 'unspecified';
const UNSPECIFIED_LABEL = 'Unspecified';


const buildLabelMap = (options: { label: string; value: string }[]) => {
    const map = new Map<string, string>();
    options.forEach(({value, label}) => {
        map.set(value, label);
    });
    map.set(UNSPECIFIED_KEY, UNSPECIFIED_LABEL);
    return map;
};

const GenerationPage = () => {
    const [form] = Form.useForm<GenerateParam>();
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [messageApi, contextHolder] = message.useMessage();
    const [generatedQuestions, setGeneratedQuestions] = useState<Question[]>([]);
    const [generatedAt, setGeneratedAt] = useState<Date | null>(null);
    const [downloadHandle, setDownloadHandle] = useState<JsonDownloadHandle | null>(null);

    const categoryLabelMap = useMemo(() => buildLabelMap(CATEGORY_OPTIONS), []);
    const difficultyLabelMap = useMemo(() => buildLabelMap(DIFFICULTY_OPTIONS), []);
    const typeLabelMap = useMemo(() => buildLabelMap(TYPE_OPTIONS), []);

    const updateDownloadHandle = useCallback((questions: Question[]) => {
        setDownloadHandle((previous) => {
            if (previous) {
                previous.revoke();
            }

            if (questions.length === 0) {
                return null;
            }

            return createJsonDownloadHandle(questions);
        });
    }, []);

    const handleGenerate = async (values: GenerateParam) => {
        setIsSubmitting(true);
        const messageKey = 'generate-loading';
        messageApi.open({key: messageKey, type: 'loading', content: 'Generating questions...', duration: 0});
        try {
            const questions = await generateQuestions({
                number: values.number,
                category: values.category,
                difficulty: values.difficulty,
                type: values.type,
            });
            setGeneratedQuestions(questions);
            setGeneratedAt(new Date());
            updateDownloadHandle(questions);
            messageApi.open({key: messageKey, type: 'success', content: 'Generation request submitted'});
        } catch (err) {
            const errorMessage = err instanceof Error ? err.message : 'Failed to generate questions';
            messageApi.open({key: messageKey, type: 'error', content: errorMessage});
        } finally {
            setIsSubmitting(false);
        }
    };

    const summary: GenerationSummary = useMemo(() => {
        const countBy = (
            key: keyof Pick<Question, 'category' | 'difficulty' | 'type'>,
            labelMap: Map<string, string>,
        ) => {
            const counts = new Map<string, number>();

            generatedQuestions.forEach((item) => {
                const rawValue = (item[key] ?? '') as string;
                const value = rawValue || UNSPECIFIED_KEY;
                counts.set(value, (counts.get(value) ?? 0) + 1);
            });

            return Array.from(counts.entries())
                .map(([value, count]) => ({
                    value,
                    label: labelMap.get(value) ?? value,
                    count,
                }))
                .sort((a, b) => b.count - a.count);
        };

        return {
            total: generatedQuestions.length,
            category: countBy('category', categoryLabelMap),
            difficulty: countBy('difficulty', difficultyLabelMap),
            type: countBy('type', typeLabelMap),
        };
    }, [categoryLabelMap, difficultyLabelMap, generatedQuestions, typeLabelMap]);

    const previewItems = useMemo(() => generatedQuestions.slice(0, 5), [generatedQuestions]);

    useEffect(() => () => {
        if (downloadHandle) {
            downloadHandle.revoke();
        }
    }, [downloadHandle]);

    return (
        <div className={'container'}>
            {contextHolder}
            <Row gutter={[24, 24]} align="stretch" className="generation-layout-row">
                <Col xs={24} lg={14} className="generation-left-column">
                    <div className="generation-left-content">
                        <Card className="generation-form-card">
                            <Space direction="vertical" size={24} style={{width: '100%'}}>
                            <div>
                                <Title level={3} style={{marginBottom: 8}}>Generate Questions</Title>
                                <Paragraph type="secondary" style={{marginBottom: 0}}>
                                    Fill out the fields below to generate new logic questions. Only the quantity is
                                    mandatory;
                                    any filters you leave empty will be treated as "any" on the server side.
                                </Paragraph>
                            </div>
                            <Spin spinning={isSubmitting} tip="Generating...">
                                <Form<GenerateParam>
                                    layout="vertical"
                                    form={form}
                                    initialValues={{number: 10}}
                                    onFinish={handleGenerate}
                                >
                                    <Form.Item
                                        label="Number of Questions"
                                        name="number"
                                        rules={[
                                            {required: true, message: 'Please provide a number of questions'},
                                            {
                                                validator: (_, value) => {
                                                    if (typeof value === 'number' && Number.isFinite(value) && value >= 1 && value <= 200) {
                                                        return Promise.resolve();
                                                    }

                                                    return Promise.reject(new Error('Please enter a valid number between 1 and 200'));
                                                },
                                            },
                                        ]}
                                    >
                                        <InputNumber min={1} max={200} step={1} style={{width: '100%'}}
                                                     placeholder="e.g. 20"/>
                                    </Form.Item>

                                    <Form.Item label="Category" name="category">
                                        <Select allowClear placeholder="Select a category" options={CATEGORY_OPTIONS}/>
                                    </Form.Item>

                                    <Form.Item label="Difficulty" name="difficulty">
                                        <Select allowClear placeholder="Select difficulty"
                                                options={DIFFICULTY_OPTIONS}/>
                                    </Form.Item>

                                    <Form.Item label="Type" name="type">
                                        <Select allowClear placeholder="Select question type" options={TYPE_OPTIONS}/>
                                    </Form.Item>

                                    <Form.Item>
                                        <Button
                                            type="primary"
                                            size="large"
                                            loading={isSubmitting}
                                            htmlType="submit"
                                        >
                                            Generate
                                        </Button>
                                    </Form.Item>
                                </Form>
                            </Spin>
                            </Space>
                        </Card>
                        <div className="generation-tips-wrapper">
                            <Card>
                                <Paragraph type="secondary">
                                    Generation Tips: Combining filters helps you target a specific set of questions. For
                                    example, pick Equivalence with Hard difficulty to challenge advanced students.
                                </Paragraph>
                            </Card>
                        </div>
                    </div>
                </Col>
                <Col xs={24} lg={10} className="generation-right-column">
                    <div className="generation-preview-wrapper">
                        <GenerationPreview
                            summary={summary}
                            previewItems={previewItems}
                            generatedAt={generatedAt}
                            downloadUrl={downloadHandle?.url}
                        />
                    </div>
                </Col>
            </Row>
        </div>
    );
};

export default GenerationPage;
