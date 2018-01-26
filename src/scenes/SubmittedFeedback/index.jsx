// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import Alert from 'shared/Alert';
import IssueCards from 'scenes/SubmittedFeedback/IssueCards';

import { loadIssues } from './ducks';

class SubmittedFeedback extends Component {
  constructor(props) {
    super(props);
    this.state = { issues: null, hasError: false };
  }
  componentDidMount() {
    document.title = 'Transcom PPP: Submitted Feedback';
    this.props.loadIssues();
  }
  render() {
    const { hasError } = this.state;
    const { issues } = this.props;
    return (
      <div className="usa-grid">
        <h1>Submitted Feedback</h1>
        {hasError && (
          <Alert type="error" heading="Server Error">
            There was a problem loading the issues from the server.
          </Alert>
        )}
        {!hasError && <IssueCards issues={issues} />}
      </div>
    );
  }
}

SubmittedFeedback.propTypes = {
  loadIssues: PropTypes.func.isRequired,
  issues: PropTypes.array.isRequired, // add shape
};

function mapStateToProps(state) {
  return { issues: state.issues };
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadIssues }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(SubmittedFeedback);
