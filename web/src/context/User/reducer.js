export const initialState = {
  user: null,
  isAuthenticated: false,
};

export const actions = {
  SET_USER: 'SET_USER',
  LOGOUT: 'LOGOUT',
};

export default function userReducer(state, action) {
  switch (action.type) {
    case actions.SET_USER:
      return {
        ...state,
        user: action.payload,
        isAuthenticated: true,
      };
    case actions.LOGOUT:
      return {
        ...state,
        user: null,
        isAuthenticated: false,
      };
    default:
      return state;
  }
}
