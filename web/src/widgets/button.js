import style from './button.scss';
import React, {Component} from 'react';

class Button extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    let classes = this.props.class.split(' ').map((name) => {
      return style[[name]]
    });
    return (
      <button type={this.props.type} className={classes} disabled={this.props.disabled}>{this.props.text}</button>
    )
  }
}

export default Button;
