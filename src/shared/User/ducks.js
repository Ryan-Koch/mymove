import * as Cookies from 'js-cookie';
const LOGOUT = 'USER|LOGOUT';

const LOAD_USER_AND_TOKEN = 'LOAD_USER_AND_TOKEN';

function getUserInfo() {
  const cookie = Cookies.get('user_session');
  return {
    jwt: cookie || null,
    isLoggedIn: cookie && true,
  };
}

export function loadUserAndToken() {
  console.log('action creator');
  const jwt = getUserInfo();
  return { type: LOAD_USER_AND_TOKEN, payload: jwt };
}

export function logOut() {
  //todo: call server endpoint, clean up cookies?
  return {
    type: LOGOUT,
  };
}

const initialState = getUserInfo();

const userReducer = (state = initialState, action) => {
  switch (action.type) {
    case LOAD_USER_AND_TOKEN:
      return action.payload;
    case LOGOUT:
      return { isLoggedIn: false, jwt: null };
    default: {
      return state;
    }
  }
};

export default userReducer;
