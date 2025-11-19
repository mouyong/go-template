import { post, get } from './api';

const TOKEN_KEY = 'token';
const USER_KEY = 'user';

export const login = async (credentials) => {
  const response = await post('/api/auth/login', credentials);
  const { token, user } = response.data;
  localStorage.setItem(TOKEN_KEY, token);
  localStorage.setItem(USER_KEY, JSON.stringify(user));
  return response.data;
};

export const register = async (userData) => {
  const response = await post('/api/auth/register', userData);
  return response.data;
};

export const logout = () => {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(USER_KEY);
};

export const getToken = () => {
  return localStorage.getItem(TOKEN_KEY);
};

export const getUser = () => {
  const user = localStorage.getItem(USER_KEY);
  return user ? JSON.parse(user) : null;
};

export const isAuthenticated = () => {
  return !!getToken();
};

export const getCurrentUser = async () => {
  const response = await get('/api/user/me');
  return response.data;
};
