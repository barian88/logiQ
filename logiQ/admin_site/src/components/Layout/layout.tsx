import React, {useMemo, useState} from 'react';
import {
    BarChartOutlined,
    HourglassOutlined,
    SettingOutlined,
    TableOutlined, UserOutlined,
} from '@ant-design/icons';
import {Avatar, type MenuProps} from 'antd';
import {Layout, Menu, theme} from 'antd';
import {Header} from 'antd/es/layout/layout';
import {Outlet, useLocation, useNavigate} from 'react-router-dom';

const {Content, Sider} = Layout;

const MyLayout: React.FC = () => {
    const [collapsed, setCollapsed] = useState(false);
    const {
        token: {colorBgContainer},
    } = theme.useToken();
    const navigate = useNavigate();
    const location = useLocation();

    const menuItems = useMemo<MenuProps['items']>(() => ([

        {
            key: '/table',
            icon: React.createElement(TableOutlined),
            label: 'Questions',
        },
        {
            key: '/generation',
            icon: React.createElement(HourglassOutlined),
            label: 'Generation',
        },
        {
            key: '/statistics',
            icon: React.createElement(BarChartOutlined),
            label: 'Statistics',
        },
        {
            key: '/settings',
            icon: React.createElement(SettingOutlined),
            label: 'Settings',
            disabled: true,
        },
    ]), []);

    const selectedKey = useMemo(() => {
        if (location.pathname === '/') {
            return '/table';
        }

        const match = menuItems
            ?.map((item) => (typeof item?.key === 'string' ? item.key : ''))
            .find((key) => key && location.pathname.startsWith(key));

        return match || location.pathname;
    }, [location.pathname, menuItems]);

    return (
        <Layout style={{maxHeight: '100vh', minHeight: '100vh'}}>
            <Sider theme={'light'} collapsible collapsed={collapsed} onCollapse={(value) => setCollapsed(value)}>
                <div style={{ display: 'flex' }}>
                    <Avatar style={{ backgroundColor: '#35abec', margin: '12px auto' }} icon={<UserOutlined />} />
                </div>
                <Menu
                    selectedKeys={[selectedKey]}
                    mode="inline"
                    items={menuItems}
                    onClick={({key}) => {
                        if (typeof key === 'string') {
                            navigate(key);
                        }
                    }}
                />
            </Sider>
            <Layout>
                <Header style={{padding: '0 16px', background: colorBgContainer, height: 'fit-content'}} >
                    <div style={{display: 'flex', alignItems: 'center'}}>
                        <img src={'/src/assets/logo.png'} alt='logo' style={{height: '30px', width: '30px', verticalAlign: 'middle'}} />
                        <h2 style={{marginLeft: '16px',marginBottom: '0'}}>LogiQ Administration</h2>
                    </div>
                </Header>
                <Content style={{margin: '16px 16px 0 16px'}}>
                    <Outlet />
                </Content>
                {/*<Footer style={{textAlign: 'center'}}>*/}
                {/*    LogiQ ©{new Date().getFullYear()} Created by Gu Jianyang*/}
                {/*</Footer>*/}
            </Layout>
        </Layout>
    );
};

export default MyLayout;
