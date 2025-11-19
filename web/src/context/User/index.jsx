import { createContext, useContext, useReducer, useEffect } from 'react';
import { getUser, logout as authLogout } from '@/helpers/auth';
import userReducer, { initialState, actions } from './reducer';

const UserContext = createContext(null);

export function UserProvider({ children }) {
  const [state, dispatch] = useReducer(userReducer, initialState);

  useEffect(() => {
    const user = getUser();
    if (user) {
      dispatch({ type: actions.SET_USER, payload: user });
    }
  }, []);

  const setUser = (user) => {
    dispatch({ type: actions.SET_USER, payload: user });
  };

  const logout = () => {
    authLogout();
    dispatch({ type: actions.LOGOUT });
  };

  const value = {
    user: state.user,
    isAuthenticated: state.isAuthenticated,
    setUser,
    logout,
  };

  return <UserContext.Provider value={value}>{children}</UserContext.Provider>;
}

export function useUser() {
  const context = useContext(UserContext);
  if (!context) {
    throw new Error('useUser must be used within a UserProvider');
  }
  return context;
}
