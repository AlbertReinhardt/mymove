import { combineReducers } from 'redux';

import { routerReducer } from 'react-router-redux';
import { issues } from 'scenes/SubmittedFeedback/ducks';

export const appReducer = combineReducers({
  router: routerReducer,
  issues,
});

export default appReducer;
