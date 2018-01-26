import { IssuesIndex } from 'shared/api';
//types
const LOAD_ISSUES_REQUEST = 'LOAD_ISSUES_REQUEST';
const LOAD_ISSUES_SUCCESS = 'LOAD_ISSUES_SUCCESS';
const LOAD_ISSUES_FAILURE = 'LOAD_ISSUES_FAILURE';

//actions
const createLoadIssuesRequest = () => ({
  type: LOAD_ISSUES_REQUEST,
});

const createLoadIssuesSuccess = items => ({
  type: LOAD_ISSUES_SUCCESS,
  items: items,
});

const createLoadIssuesError = error => ({
  type: LOAD_ISSUES_FAILURE,
  error,
});

//action creators
export function loadIssues() {
  // Interpreted by the thunk middleware:
  return function(dispatch, getState) {
    dispatch(createLoadIssuesRequest());
    IssuesIndex()
      .then(items => dispatch(createLoadIssuesSuccess(items)))
      .catch(error => dispatch(createLoadIssuesError(error)));
  };
}

//reducer
export function issues(state = [], action) {
  switch (action.type) {
    case LOAD_ISSUES_SUCCESS:
      return action.items;

    default:
      return state;
  }
}
