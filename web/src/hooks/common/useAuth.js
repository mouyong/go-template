import { useUser } from '@/context/User';
import { isAuthenticated } from '@/helpers/auth';

export default function useAuth() {
  const { user, setUser, logout } = useUser();

  return {
    user,
    setUser,
    logout,
    isAuthenticated: isAuthenticated(),
    isAdmin: user?.role === 'admin',
  };
}
