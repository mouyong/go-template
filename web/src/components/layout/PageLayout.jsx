import { Outlet } from 'react-router-dom';
import { Layout } from '@douyinfe/semi-ui';
import SideBar from './SideBar';
import Header from './Header';
import useSidebarCollapsed from '@/hooks/common/useSidebarCollapsed';

const { Sider, Content } = Layout;

export default function PageLayout() {
  const [collapsed, toggleCollapsed] = useSidebarCollapsed();
  const sidebarWidth = collapsed ? 60 : 240;
  const contentMarginLeft = collapsed ? 68 : 264; // 侧边栏宽度 + 间距 (8px 或 24px)

  return (
    <Layout style={{ height: '100vh', display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
      {/* Header - fixed 定位在顶部 */}
      <Layout.Header
        style={{
          position: 'fixed',
          width: '100%',
          top: 0,
          zIndex: 100,
          height: '64px',
          padding: 0,
        }}
      >
        <Header />
      </Layout.Header>

      {/* 中间层 - 包含 Sider 和 Content */}
      <Layout style={{ overflow: 'auto', display: 'flex', flexDirection: 'column', marginTop: '64px' }}>
        {/* Sider - fixed 定位在左侧，Header 下方 */}
        <Sider
          style={{
            position: 'fixed',
            left: 0,
            top: '64px',
            zIndex: 99,
            border: 'none',
            paddingRight: '0',
            height: 'calc(100vh - 64px)',
            width: `${sidebarWidth}px`,
          }}
        >
          <SideBar collapsed={collapsed} onToggleCollapse={toggleCollapsed} />
        </Sider>

        {/* 内容容器 - marginLeft 避让侧边栏并留出间距 */}
        <Layout
          style={{
            marginLeft: `${contentMarginLeft}px`,
            transition: 'margin-left 0.3s ease',
            flex: '1 1 auto',
            display: 'flex',
            flexDirection: 'column',
          }}
        >
          <Content
            style={{
              flex: '1 0 auto',
              overflowY: 'hidden',
              padding: '24px',
              position: 'relative',
            }}
          >
            <Outlet />
          </Content>
        </Layout>
      </Layout>
    </Layout>
  );
}
