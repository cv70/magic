import { Layout, Menu, Dropdown, Avatar, Space, Breadcrumb, Switch } from 'antd'
import { LogoutOutlined, SettingOutlined, UserOutlined, BgColorsOutlined, GlobalOutlined } from '@ant-design/icons'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { useAuthStore } from '../../stores/authStore'
import { useThemeStore } from '../../stores/themeStore'

const { Sider, Header, Content } = Layout

const menuItems: any[] = [
  { key: '/dashboard', label: 'Dashboard', title: '概览' },
  {
    key: 'drafts',
    label: '内容创作',
    title: '内容创作',
    children: [
      { key: '/drafts', label: '草稿中心' },
      { key: '/content', label: '已发布' },
      { key: '/ai-studio', label: 'AI工作室' },
    ],
  },
  {
    key: 'publishing',
    label: '发布运营',
    title: '发布运营',
    children: [
      { key: '/publish', label: '发布管理' },
      { key: '/analytics', label: '数据分析' },
    ],
  },
  { key: '/square', label: '内容广场', icon: <GlobalOutlined />, title: '内容广场' },
  { key: '/settings', label: '设置', title: '设置' },
]

export function AppLayout() {
  const location = useLocation()
  const navigate = useNavigate()
  const { user, logout } = useAuthStore()
  const { isDarkMode, toggleDarkMode } = useThemeStore()

  const selectedKey = menuItems.find((item) => item.key === location.pathname)?.key ||
    menuItems.flatMap((item) => item.children || []).find((item) => item.key === location.pathname)?.key ||
    location.pathname

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const userMenu: any[] = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人资料',
    },
    {
      key: 'theme',
      label: (
        <Space>
          <BgColorsOutlined />
          <span>深色模式</span>
          <Switch size="small" checked={isDarkMode} onChange={toggleDarkMode} />
        </Space>
      ),
      disabled: true,
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '设置',
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ]

  const breadcrumbMap: { [key: string]: string[] } = {
    '/dashboard': ['首页'],
    '/drafts': ['内容创作', '草稿中心'],
    '/content': ['内容创作', '已发布'],
    '/ai-studio': ['内容创作', 'AI工作室'],
    '/publish': ['发布运营', '发布管理'],
    '/analytics': ['发布运营', '数据分析'],
    '/square': ['内容广场'],
    '/settings': ['系统', '设置'],
  }

  const breadcrumb = breadcrumbMap[location.pathname] || ['首页']

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider width={240} theme="dark" breakpoint="lg" collapsedWidth={0}>
        <div style={{ padding: '16px', color: 'white', fontSize: '18px', fontWeight: 'bold' }}>
          🚀 Magic Engine
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[selectedKey as string]}
          items={menuItems}
          onClick={(e) => {
            navigate(e.key)
          }}
        />
      </Sider>

      <Layout>
        <Header
          style={{
            background: '#fff',
            padding: '0 24px',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
          }}
        >
          <Breadcrumb items={breadcrumb.map((item) => ({ title: item }))} />
          <Space>
            <span>{user?.username}</span>
            <Dropdown menu={{ items: userMenu }}>
              <Avatar style={{ backgroundColor: '#87d068' }} icon={<UserOutlined />} />
            </Dropdown>
          </Space>
        </Header>

        <Content style={{ margin: '24px 16px', padding: 24, background: '#fff', borderRadius: '4px' }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
