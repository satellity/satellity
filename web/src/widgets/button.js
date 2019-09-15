import React, {Component} from 'react';
import {Link} from 'react-router-dom';

class Button extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <Link to={this.props.action}>{this.props.text}</Link>
    )
  }
}

export default Button;
