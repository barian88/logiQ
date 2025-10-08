import {useState} from 'react';
import {useLocation, useNavigate, type Location} from 'react-router-dom';
import {LockOutlined, MailOutlined} from '@ant-design/icons';
import {Button, Card, Col, Form, Input, Row, Typography, message} from 'antd';

import './login.css';
import {login as loginRequest} from './services.ts';
import {useAuth} from '../../context/AuthContext.tsx';

const {Title, Paragraph} = Typography;

interface LoginFormValues {
    email: string;
    password: string;
}

const LoginPage = () => {
    const [form] = Form.useForm<LoginFormValues>();
    const navigate = useNavigate();
    const location = useLocation();
    const {login} = useAuth();
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [messageApi, contextHolder] = message.useMessage();

    const redirectPath = (location.state as {from?: Location})?.from?.pathname ?? '/';

    const handleSubmit = async (values: LoginFormValues) => {
        setIsSubmitting(true);
        try {
            const response = await loginRequest(values.email, values.password);
            const token = typeof response === 'string' ? response : response.token;
            if (!token) {
                throw new Error('Invalid token received');
            }
            login(token);
            navigate(redirectPath, {replace: true});
        } catch (error) {
            const errorMessage = error instanceof Error ? error.message : 'Login failed';
            messageApi.error(errorMessage);
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className="login-page">
            {contextHolder}
            <Row justify="center" align="middle" className="login-container">
                <Col xs={22} sm={16} md={12} lg={8}>
                    <Card className="login-card">
                        <div className="login-header">
                            <Title level={3} style={{marginBottom: 8}}>Welcome Back</Title>
                            <Paragraph type="secondary" style={{margin: 0}}>
                                Sign in with your email and password to continue.
                            </Paragraph>
                        </div>
                        <Form<LoginFormValues>
                            layout="vertical"
                            form={form}
                            requiredMark={false}
                            onFinish={handleSubmit}
                        >
                            <Form.Item label="Email" name="email" rules={[{required: true, message: 'Please enter your email'}]}>
                                <Input size="large" prefix={<MailOutlined />} placeholder="you@example.com" allowClear />
                            </Form.Item>

                            <Form.Item label="Password" name="password" rules={[{required: true, message: 'Please enter your password'}]}>
                                <Input.Password size="large" prefix={<LockOutlined />} placeholder="Password" />
                            </Form.Item>

                            <Form.Item>
                                <Button type="primary" htmlType="submit" size="large" block loading={isSubmitting}>
                                    Log in
                                </Button>
                            </Form.Item>
                        </Form>
                    </Card>
                </Col>
            </Row>
        </div>
    );
};

export default LoginPage;
