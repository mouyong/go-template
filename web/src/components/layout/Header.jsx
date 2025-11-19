import { Avatar, Dropdown } from '@douyinfe/semi-ui';
import { IconUser } from '@douyinfe/semi-icons';
import { useNavigate } from 'react-router-dom';
import { useUser } from '@/context/User';

export default function Header() {
  const navigate = useNavigate();
  const { user, logout } = useUser();

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  return (
    <>
      <div className='text-lg font-semibold'>Web App</div>
      <Dropdown
        position='bottomRight'
        render={
          <Dropdown.Menu>
            <Dropdown.Item onClick={() => navigate('/profile')}>Profile</Dropdown.Item>
            <Dropdown.Divider />
            <Dropdown.Item onClick={handleLogout}>Logout</Dropdown.Item>
          </Dropdown.Menu>
        }
      >
        <Avatar size='small' style={{ cursor: 'pointer' }}>
          {user?.username?.[0]?.toUpperCase() || <IconUser />}
        </Avatar>
      </Dropdown>
    </>
  );
}
