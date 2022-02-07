import React from 'react';
import { NotificationBadge } from '@patternfly/react-core';
import BellIcon from '@patternfly/react-icons/dist/esm/icons/bell-icon';

class BasicHint extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      unreadVariant: 'unread',
      attentionVariant: 'attention'
    };
    this.onFirstClick = () => {
      this.setState({
        unreadVariant: 'read'
      });
    };
    this.onSecondClick = () => {
      this.setState({
        attentionVariant: 'read'
      });
    };
  }

  render() {
    const { unreadVariant, attentionVariant } = this.state;
    return (
      <div className="pf-t-dark">
        <NotificationBadge variant={unreadVariant} onClick={this.onFirstClick} aria-label="First notifications" />
        <NotificationBadge variant={attentionVariant} onClick={this.onSecondClick} aria-label="Second notifications" />
      </div>
    );
  }
}

export default BasicHint