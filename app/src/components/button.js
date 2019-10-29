import style from './button.module.scss';
import React, {Component} from 'react';
import {Link} from 'react-router-dom';
import Loading from './loading.js';

class Button extends Component {
  render() {
    let classes = this.props.classes.split(' ').map((name) => {
      return style[[name]]
    }).join(' ');

    if (this.props.original && this.props.type === 'link') {
      return (
        <a href={this.props.action} className={classes}>{this.props.text}</a>
      );
    }

    if (this.props.type === 'link') {
      return (
        <Link to={this.props.action} className={classes}>{this.props.text}</Link>
      );
    }

    if (this.props.type === 'button') {
      return (
        <button type={this.props.type} className={classes} disabled={this.props.disabled} onClick={this.props.click}>
            {this.props.disabled && <Loading class='small white' />}
              &nbsp;{this.props.text}
          </button>
      )
    }

    return (
      <button type={this.props.type} className={classes} disabled={this.props.disabled}>
          {this.props.disabled && <Loading class='small white' />}
            &nbsp;{this.props.text}
        </button>
    )
  }
}

export default Button;
