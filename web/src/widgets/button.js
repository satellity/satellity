import style from './button.scss';
import React, {Component} from 'react';
import Loading from './loading.js';

class Button extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    let classes = this.props.class.split(' ').map((name) => {
      return style[[name]]
    });
    return (
      <button type={this.props.type} className={classes} disabled={this.props.disabled}>
        {this.props.disabled && <Loading class='small white' />}
        &nbsp;{this.props.text}
      </button>
    )
  }
}

export default Button;
