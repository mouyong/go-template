import { Nav, Button } from '@douyinfe/semi-ui';
import { IconHome, IconSetting, IconUser, IconChevronLeft } from '@douyinfe/semi-icons';
import { useNavigate, useLocation } from 'react-router-dom';

export default function SideBar({ collapsed, onToggleCollapse }) {
  const navigate = useNavigate();
  const location = useLocation();
  const sidebarWidth = collapsed ? 60 : 240;

  const items = [
    { itemKey: '/dashboard', text: 'Dashboard', icon: <IconHome /> },
    { itemKey: '/users', text: 'Users', icon: <IconUser /> },
    { itemKey: '/settings', text: 'Settings', icon: <IconSetting /> },
  ];

  return (
    <div
      className='sidebar-container'
      style={{
        width: `${sidebarWidth}px`,
      }}
    >
      <Nav
        className='sidebar-nav'
        isCollapsed={collapsed}
        selectedKeys={[location.pathname]}
        items={items}
        renderWrapper={({ itemElement, props }) => {
          const to = props.itemKey;
          if (!to) return itemElement;

          return (
            <a
              href={to}
              style={{ textDecoration: 'none', color: 'inherit', display: 'block' }}
              onClick={(e) => {
                // 正常点击使用 SPA 导航
                if (!e.ctrlKey && !e.metaKey && e.button === 0) {
                  e.preventDefault();
                  navigate(to);
                }
                // Ctrl/Cmd + Click 或右键使用原生行为，在新标签页打开
              }}
            >
              {itemElement}
            </a>
          );
        }}
      />

      {/* 底部折叠按钮 */}
      <div className='sidebar-collapse-button'>
        <Button
          theme='borderless'
          type='tertiary'
          size='small'
          icon={
            <IconChevronLeft
              style={{
                transform: collapsed ? 'rotate(180deg)' : 'rotate(0deg)',
                transition: 'transform 0.3s ease',
              }}
            />
          }
          onClick={onToggleCollapse}
          block={!collapsed}
          style={
            collapsed
              ? { width: 36, height: 32, padding: 0 }
              : { width: '100%', height: 32 }
          }
        >
          {!collapsed ? '收起侧边栏' : null}
        </Button>
      </div>
    </div>
  );
}
